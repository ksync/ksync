package ksync

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/syncthing"
)

// Syncthing represents the local syncthing process.
type Syncthing struct {
	cmd *exec.Cmd
}

// NewSyncthing constructs a new Syncthing.
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

// Propogate stdout/stderr into the ksync logs for debugging.
func (s *Syncthing) outputHandler() error {
	logger := log.WithFields(log.Fields{
		"name": "syncthing",
	})

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	outScanner := bufio.NewScanner(stdout)

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return err
	}

	errScanner := bufio.NewScanner(stderr)

	go func() {
		for outScanner.Scan() {
			logger.Debug(outScanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			logger.Warn(errScanner.Text())
		}
	}()

	return nil
}

func (s *Syncthing) binPath() string {
	return filepath.Join(cli.ConfigPath(), "bin", "syncthing")
}

// HasBinary checks whether the syncthing binary exists in the correct location
// or not.
func (s *Syncthing) HasBinary() bool {
	if _, err := os.Stat(s.binPath()); err != nil {
		return false
	}

	return true
}

// Fetch the latest syncthing binary to Syncthing.binPath().
func (s *Syncthing) Fetch() error {
	return syncthing.Fetch(s.binPath())
}

// To make sure no odd devices or folders are being synced after edge cases
// (such as the process being kill'd), the config and db are blown away before
// each run.
func (s *Syncthing) resetState() error {
	base := filepath.Join(cli.ConfigPath(), "syncthing")
	if err := os.RemoveAll(base); err != nil {
		return err
	}

	return syncthing.ResetConfig(filepath.Join(base, "config.xml"))
}

// Run starts up a local syncthing process to serve files from.
func (s *Syncthing) Run() error {
	if !s.HasBinary() {
		return fmt.Errorf("missing pre-requisites, run init to fix")
	}

	if err := s.resetState(); err != nil {
		return err
	}

	path := filepath.Join(cli.ConfigPath(), "bin", "syncthing")

	address := fmt.Sprintf("localhost:%d", viper.GetInt("syncthing-port"))

	cmdArgs := []string{
		"-gui-address", address,
		"-gui-apikey", viper.GetString("apikey"),
		"-home", filepath.Join(cli.ConfigPath(), "syncthing"),
		"-no-browser",
	}

	s.cmd = exec.Command(path, cmdArgs...) //nolint: gas

	if err := s.outputHandler(); err != nil {
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
