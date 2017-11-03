package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/docker"
)

// Stop halts a container
func Stop(name string) error {
	status, err := GetStatus(name)
	if err != nil {
		return err
	}

	if !status.Running {
		return fmt.Errorf("must start before you can stop: %s", name)
	}

	if err := docker.Client.ContainerRemove(
		context.Background(),
		status.ID,
		types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrap(err, "could not remove")
	}

	log.WithFields(status.Fields()).Debug("container removed")

	return nil
}
