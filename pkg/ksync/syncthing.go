package ksync

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	daemon "github.com/sevlyar/go-daemon"

	"github.com/vapor-ware/ksync/pkg/debug"
)

type Syncthing struct {
	cmd *exec.Cmd
}

func NewSyncthing() *Syncthing {
	return &Syncthing{}
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

	cmdArgs := []string{
		"-gui-address", fmt.Sprintf("localhost:%d", viper.GetInt("syncthing-port")),
		"-gui-apikey", viper.GetString("apikey"),
		"-home", filepath.Dir(viper.ConfigFileUsed()),
		"-no-browser",
	}

	s.cmd = exec.Command(path, cmdArgs...)

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


	return nil
}

// Stop halts the background process and cleans up.
func (s *Syncthing) Stop() error {
	defer s.cmd.Process.Wait() //nolint: errcheck
	return s.cmd.Process.Kill()
}

func (s *Syncthing) Daemonize() error {
	context := &daemon.Context{
		PidFileName: filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "daemon.pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "daemon.log"),
		LogFilePerm: 0640,
		WorkDir:     filepath.Dir(viper.ConfigFileUsed()),
		// Umask:       027,
		Args: []string{"","watch"},
	}

	daemon, err := context.Reborn()
	if err != nil {
		return err
	}
	if daemon != nil {
		return nil
	}

	log.Debug("daemonizing")

	defer context.Release()

	return nil
}
