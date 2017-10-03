package ksync

import (
	"bufio"
	"io"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func pipeScanner(pipe io.ReadCloser, logger func(...interface{})) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			logger(scanner.Text())
		}
	}()
}

// TODO: this is kinda ugly, should it be a type?
// TODO: logging is tough to read right now because of all the extra stuff (tags).
func LogCmdStream(cmd *exec.Cmd) error {
	logger := log.WithFields(log.Fields{
		"path": cmd.Path,
		"args": cmd.Args,
	})

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	pipeScanner(stderr, logger.Warn)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	pipeScanner(stdout, logger.Debug)

	return nil
}
