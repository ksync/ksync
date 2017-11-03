package service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/docker"
)

// Start runs a container
func Start(cntr *container.ContainerCreateCreatedBody) error {
	if err := docker.Client.ContainerStart(
		context.Background(),
		cntr.ID,
		types.ContainerStartOptions{}); err != nil {
		return err
	}

	// TODO: make sure that the container comes up successfully.

	log.WithFields(log.Fields{
		"id": cntr.ID,
	}).Debug("container started")

	return nil
}
