package service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/docker"
)

// Start takes container configuration and starts it locally.
func Start(
	containerConfig *container.Config,
	hostConfig *container.HostConfig,
	networkConfig *network.NetworkingConfig,
	containerName string) error {

	imageExists, err := docker.HasImage(containerConfig.Image)
	if err != nil {
		log.Fatal(err)
	}

	if !imageExists {
		// Note: this won't work from inside the watch/run containers as there is no
		// docker client and/or connection.
		if pullErr := Pull(containerConfig.Image); pullErr != nil {
			log.Fatal(pullErr)
		}
	}

	cntr, err := Create(containerConfig, hostConfig, networkConfig, containerName)

	if err != nil {
		return err
	}

	if err := docker.Client.ContainerStart(
		context.Background(),
		cntr.ID,
		types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id": cntr.ID,
	}).Debug("container started")

	return nil
}
