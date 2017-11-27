package ksync

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

var (
	maxConnectionRetries = 1
	tooSoonReset         = 3 * time.Second
)

// Mirror is the definition of a sync from the local host to a remote container.
type Mirror struct {
	SpecName        string
	RemoteContainer *RemoteContainer
	Reload          bool
	// TODO: should this be a SyncPath? Seems like it ...
	LocalPath  string
	RemotePath string

	cmd               *exec.Cmd
	connectionRetries int
	retryLock         sync.Mutex

	restartContainer chan bool
	clean            chan bool
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

func (m *Mirror) scanner(pipe io.Reader, logger func(...interface{})) error {
	scanner := bufio.NewScanner(pipe)
	pattern, err := regexp.Compile("INFO  Sending")
	if err != nil {
		return err
	}

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if m.Reload && pattern.MatchString(line) {
				m.restartContainer <- true
			}
			logger(line)
		}
	}()

	return nil
}

func (m *Mirror) initLogs() error {
	logger := log.WithFields(log.Fields{
		"name": m.SpecName,
	})

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}
	if scanErr := m.scanner(stderr, logger.Warn); scanErr != nil {
		return scanErr
	}

	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	return m.scanner(stdout, logger.Debug)
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
	signal.Notify(teardown, os.Interrupt, os.Kill)
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

	m.cmd = exec.Command("java", cmdArgs...)// #nosec
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
