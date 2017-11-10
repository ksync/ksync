package service

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Pull takes an image name and pulls it locally.
func Pull(name string) error {
	log.WithFields(log.Fields{
		"image": name,
	}).Debug("pulling image")

	cmd := exec.Command("docker", "pull", name)

	// TODO: make this configurable
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
