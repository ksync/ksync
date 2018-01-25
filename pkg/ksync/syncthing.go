package ksync

import (
	"bufio"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	daemon "github.com/sevlyar/go-daemon"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/syncthing"
)

type Syncthing struct {
	Server *syncthing.Server

	cmd *exec.Cmd
}

func NewSyncthing() *Syncthing {
	return &Syncthing{
		Server: &syncthing.Server{},
	}
}

func (s *Syncthing) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Syncthing) Fields() log.Fields {
	return debug.StructFields(s)
}

func (s *Syncthing) errHandler(logger func(...interface{})) error {
	stderr, err := s.cmd.StderrPipe()
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

func (s *Syncthing) lineHandler(logger func(...interface{})) error {
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			logger(line)
		}
	}()

	return nil
}

func (s *Syncthing) initLogs() error {
	logger := log.WithFields(log.Fields{
		"name": "syncthing",
	})

	if err := s.errHandler(logger.Warn); err != nil {
		return err
	}

	return s.lineHandler(logger.Debug)
}

func (s *Syncthing) Run() error {
	path := filepath.Join(
		filepath.Dir(viper.ConfigFileUsed()), "syncthing")

	s.cmd = exec.Command(path)

	if err := s.initLogs(); err != nil {
		return err
	}

	if err := s.cmd.Start(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cmd":  s.cmd.Path,
		"args": s.cmd.Args,
	}).Debug("starting syncthing")

	s.Daemonize()

	return nil
}

// Stop halts the background process and cleans up.
func (s *Syncthing) Stop() error {
	defer s.cmd.Process.Wait() //nolint: errcheck
	return s.cmd.Process.Kill()
}

func (s *Syncthing) Daemonize() error {
	context := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args: []string{"version"},
	}

	daemon, err := context.Reborn()
	if err != nil {
		return err
	}
	if daemon != nil {
		return nil
	}

	log.WithFields(log.Fields{
		"cmd":  s.cmd.Path,
		"args": s.cmd.Args,
	}).Debug("daemonizing")

	defer context.Release()

	return nil
}