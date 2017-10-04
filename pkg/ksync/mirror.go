package ksync

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type Mirror struct {
	Container  *Container
	LocalPath  string
	RemotePath string
	cmd        *exec.Cmd
}

func (this *Mirror) scanner(pipe io.ReadCloser, logger func(...interface{})) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			logger(scanner.Text())
		}
	}()
}

func (this *Mirror) initLogs() error {
	logger := log.WithFields(log.Fields{
		"path": this.cmd.Path,
		"args": this.cmd.Args,
	})

	stderr, err := this.cmd.StderrPipe()
	if err != nil {
		return err
	}
	this.scanner(stderr, logger.Warn)

	stdout, err := this.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	this.scanner(stdout, logger.Debug)

	return nil
}

func (this *Mirror) path() (string, error) {
	client, err := this.Container.Radar()
	if err != nil {
		return "", err
	}

	path, err := client.GetAbsPath(
		context.Background(), &pb.ContainerPath{this.Container.ID, this.RemotePath})
	if err != nil {
		return "", err
	}

	return path.Full, nil
}

func (this *Mirror) tunnel() (int, error) {
	tun, err := NewTunnel(this.Container.NodeName, 49172)
	if err != nil {
		return 0, err
	}

	if err := tun.Start(); err != nil {
		return 0, err
	}

	return tun.LocalPort, nil
}

func (this *Mirror) initErrorHandler() {
	// Setup the k8s runtime to fail on unreturnable error (instead of looping).
	// This helps cleanup zombie java processes.
	runtime.ErrorHandlers = append(runtime.ErrorHandlers, func(err error) {
		this.cmd.Process.Kill()
		// TODO: this makes me feel dirty, there must be a better way.
		if strings.Contains(err.Error(), "Connection refused") {
			log.Fatal(
				"Lost connection to remote radar pod. Try again (it should restart).")
		}

		log.Fatalf("unreturnable error: %v", err)
	})
}

// TODO: this takes maybe 5 seconds or so to start, show a progress bar.
// TODO: the output for this needs some thought. There should be:
//   - debug output (raw sync), this is a little tough to read right now
//   -
func (this *Mirror) Run() error {
	path, err := this.path()
	if err != nil {
		return err
	}

	port, err := this.tunnel()
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"-Xmx2G",
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-cp", "/home/thomas/work/bin/mirror-all.jar",
		"mirror.Mirror",
		"client",
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", port),
		"-l", this.LocalPath,
		"-r", path,
	}

	this.cmd = exec.Command("java", cmdArgs...)
	this.initErrorHandler()

	if err := this.initLogs(); err != nil {
		return err
	}

	if err := this.cmd.Start(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cmd":  this.cmd.Path,
		"args": this.cmd.Args,
	}).Debug("starting mirror")

	return this.cmd.Wait()
}
