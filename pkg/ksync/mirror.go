package ksync

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// Mirror is the definition of a sync from the local host to a remote container.
type Mirror struct {
	Container *Container
	// TODO: should this be a SyncPath? Seems like it ...
	LocalPath  string
	RemotePath string
	cmd        *exec.Cmd
}

func (m *Mirror) scanner(pipe io.Reader, logger func(...interface{})) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			logger(scanner.Text())
		}
	}()
}

func (m *Mirror) initLogs() error {
	logger := log.WithFields(log.Fields{
		"path": m.cmd.Path,
		"args": m.cmd.Args,
	})

	stderr, err := m.cmd.StderrPipe()
	if err != nil {
		return err
	}
	m.scanner(stderr, logger.Warn)

	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	m.scanner(stdout, logger.Debug)

	return nil
}

func (m *Mirror) path() (string, error) {
	client, err := m.Container.Radar()
	if err != nil {
		return "", err
	}

	path, err := client.GetBasePath(
		context.Background(), &pb.ContainerPath{
			ContainerId: m.Container.ID,
		})
	if err != nil {
		return "", err
	}

	return filepath.Join(path.Full, m.RemotePath), nil
}

func (m *Mirror) initErrorHandler() {
	// Setup the k8s runtime to fail on unreturnable error (instead of looping).
	// This helps cleanup zombie java processes.
	runtime.ErrorHandlers = append(runtime.ErrorHandlers, func(fromHandler error) {
		if err := m.cmd.Process.Kill(); err != nil {
			log.Fatalf("couldn't kill %v", err)
		}
		// TODO: this makes me feel dirty, there must be a better way.
		if strings.Contains(fromHandler.Error(), "Connection refused") {
			log.Fatal(
				"Lost connection to remote radar pod. Try again (it should restart).")
		}

		log.Fatalf("unreturnable error: %v", fromHandler)
	})
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
	path, err := m.path()
	if err != nil {
		return err
	}

	port, err := NewRadarInstance().MirrorConnection(m.Container.NodeName)
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"-Xmx2G",
		"-XX:+HeapDumpOnOutOfMemoryError",
		// TODO: make this generic
		"-cp", "/mirror/mirror-all.jar",
		"mirror.Mirror",
		"client",
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", port),
		"-l", m.LocalPath,
		"-r", path,
	}

	m.cmd = exec.Command("java", cmdArgs...)
	m.initErrorHandler()

	if err := m.initLogs(); err != nil {
		return err
	}

	if err := m.cmd.Start(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cmd":  m.cmd.Path,
		"args": m.cmd.Args,
	}).Debug("starting mirror")

	return m.cmd.Wait()
}
