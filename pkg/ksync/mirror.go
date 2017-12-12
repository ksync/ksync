package ksync

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/vapor-ware/ksync/pkg/debug"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

var (
	maxConnectionRetries = 1
	tooSoonReset         = 3 * time.Second
	// The assumption is that it'll not take more than 5 seconds to send a file,
	// while this isn't actually true for large files .. it is good enough for
	// now as we're unable to get real progress from mirror.
	resetStatusTime = 5 * time.Second
)

// Mirror is the definition of a sync from the local host to a remote container.
type Mirror struct {
	SpecName        string
	RemoteContainer *RemoteContainer
	Reload          bool
	LocalPath       string
	RemotePath      string
	Status          ServiceStatus

	cmd               *exec.Cmd
	connectionRetries int
	retryLock         sync.Mutex

	resetTimer       *time.Timer
	restartContainer chan bool
	clean            chan bool
}

// NewMirror is a constructor for Mirror
func NewMirror(service *Service) *Mirror {
	return &Mirror{
		SpecName:        service.SpecDetails.Name,
		RemoteContainer: service.RemoteContainer,
		Reload:          service.SpecDetails.Reload,
		LocalPath:       service.SpecDetails.LocalPath,
		RemotePath:      service.SpecDetails.RemotePath,
	}
}

func (m *Mirror) String() string {
	return debug.YamlString(m)
}

// Fields returns a set of structured fields for logging.
func (m *Mirror) Fields() log.Fields {
	return debug.StructFields(m)
}

func (m *Mirror) hotReload() error {
	m.restartContainer = make(chan bool)
	conn, err := NewRadarInstance().RadarConnection(
		m.RemoteContainer.NodeName)
	if err != nil {
		return err
	}

	client := pb.NewRadarClient(conn)

	// TODO: this is pretty naive, there are definite edge cases here where the
	// reload will happen but not actually get some files.
	go func() {
		defer conn.Close() // nolint: errcheck

		tooSoon := false
		for {
			select {
			case <-m.restartContainer:
				if tooSoon {
					continue
				}
				tooSoon = true

				log.WithFields(m.RemoteContainer.Fields()).Debug("issuing reload")

				if _, err := client.Restart(
					context.Background(), &pb.ContainerPath{
						ContainerId: m.RemoteContainer.ID,
					}); err != nil {
					log.WithFields(m.RemoteContainer.Fields()).Error(err)
					continue
				}

				log.WithFields(m.RemoteContainer.Fields()).Debug("reloaded")
			case <-time.After(tooSoonReset):
				tooSoon = false
			case <-m.clean:
				return
			}
		}
	}()

	return nil
}

func (m *Mirror) resetStatus(status ServiceStatus) {
	log.Print(status)
	if m.resetTimer != nil && m.resetTimer.Reset(resetStatusTime) {
		return
	}

	m.resetTimer = time.AfterFunc(resetStatusTime, func() { m.Status = status })
}

func (m *Mirror) errHandler(logger func(...interface{})) error {
	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)

	go func() {
		for scanner.Scan() {
			logger(scanner.Text())
		}
	}()

	return nil
}

func (m *Mirror) lineHandler(logger func(...interface{})) error {
	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	handler, err := NewLineStatus(map[ServiceStatus]string{
		ServiceConnecting: "INFO  Increasing file limit",
		ServiceConnected:  "INFO  Connected",
		ServiceWatching:   "INFO  Tree populated",
		ServiceSending:    "INFO  Sending",
		ServiceReceiving:  "INFO  Remote update",
		ServiceError:      "ERROR",
	})

	if err != nil {
		return err
	}

	log.Error("fail boat")
	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			logger(line)

			switch status := handler.Next(line); status {
			case "":
			default:
				m.Status = status
				fallthrough
			case ServiceSending:
				if m.Reload {
					m.restartContainer <- true
				}
				m.resetStatus(ServiceWatching)
			case ServiceReceiving:
				m.resetStatus(ServiceWatching)
			}
		}
	}()

	return nil
}

func (m *Mirror) initLogs() error {
	logger := log.WithFields(log.Fields{
		"name": m.SpecName,
	})

	if err := m.errHandler(logger.Warn); err != nil {
		return err
	}

	return m.lineHandler(logger.Debug)
}

func (m *Mirror) path() (string, error) {
	client, err := m.RemoteContainer.Radar()
	if err != nil {
		return "", err
	}

	path, err := client.GetBasePath(
		context.Background(), &pb.ContainerPath{
			ContainerId: m.RemoteContainer.ID,
		})
	if err != nil {
		return "", err
	}

	return filepath.Join(path.Full, m.RemotePath), nil
}

// TODO: this will fire for *every* disconnect (no matter what it is for). Need
// to filter down.
func (m *Mirror) initErrorHandler() {
	// Setup the k8s runtime to fail on unreturnable error (instead of looping).
	// This helps cleanup zombie java processes.
	runtime.ErrorHandlers = append(runtime.ErrorHandlers, func(fromHandler error) {
		// Connection refused errors suggest that mirror is restarting remotely
		// and we should just be patient. If the error isn't connection refused,
		// it is likely we have a more serious problem (such as the entire pod
		// getting rescheduled).

		if m.connectionRetries < maxConnectionRetries {
			m.retryLock.Lock()
			defer m.retryLock.Unlock()

			m.connectionRetries++
			log.WithFields(log.Fields{
				"retries": m.connectionRetries,
			}).Debug("lost connection to remote mirror server")

			return
		}

		if err := m.Stop(); err != nil {
			log.Fatalf("couldn't stop %v", err)
		}

		log.Fatal(fromHandler)
	})
}

func (m *Mirror) handleTeardown() {
	teardown := make(chan os.Signal, 2)
	signal.Notify(teardown, os.Interrupt) // Removed `os.Kill` as SIGTERMs will always be caught
	go func() {
		for {
			select {
			case <-teardown:
				m.Stop() //nolint: errcheck
			case <-m.clean:
				return
			}
		}
	}()
}

// Run starts a sync from the local host to a remote container. This is a
// long running process and will wait indefinitely (or until the background
// process dies).
// TODO: this takes maybe 5 seconds or so to start, show a progress bar.
// TODO: the output for this needs some thought. There should be:
//   - debug output (raw sync), this is a little tough to read right now
//   - state updates (disconnected, active, idle)
// TODO: stop gracefully when the remote pod goes away.
func (m *Mirror) Run() error {
	m.clean = make(chan bool)
	m.connectionRetries = 0

	path, err := m.path()
	if err != nil {
		return err
	}

	m.Status = ServiceStarting

	port, err := NewRadarInstance().MirrorConnection(m.RemoteContainer.NodeName)
	if err != nil {
		return err
	}

	jarPath := filepath.Join(
		filepath.Dir(viper.ConfigFileUsed()), "mirror-all.jar")
	cmdArgs := []string{
		"-Xmx2G",
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-cp", jarPath,
		"mirror.Mirror",
		"client",
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", port),
		"-l", m.LocalPath,
		"-r", path,
	}

	m.cmd = exec.Command("java", cmdArgs...) // #nosec
	m.initErrorHandler()

	if m.Reload {
		if err := m.hotReload(); err != nil {
			return err
		}
	}

	if err := m.initLogs(); err != nil {
		return err
	}

	m.handleTeardown()

	if err := m.cmd.Start(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cmd":  m.cmd.Path,
		"args": m.cmd.Args,
	}).Debug("starting mirror")

	return nil
}

// Stop halts the background process and cleans up.
func (m *Mirror) Stop() error {
	defer m.cmd.Process.Wait() //nolint: errcheck
	if m.clean != nil {
		close(m.clean)
	}
	m.clean = nil
	return m.cmd.Process.Kill()
}
