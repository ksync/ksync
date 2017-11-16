package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"github.com/vapor-ware/ksync/pkg/docker"
)

// Create builds a docker container after verifying that the image exists.
func Create(
	containerConfig *container.Config,
	hostConfig *container.HostConfig,
	networkConfig *network.NetworkingConfig,
	containerName string) (*container.ContainerCreateCreatedBody, error) {

	imageExists, err := docker.HasImage(containerConfig.Image)
	if err != nil {
		return nil, err
	}

	if !imageExists {
		return nil, fmt.Errorf(
			"%s does not exist, docker pull", containerConfig.Image)
	}

	cntr, err := docker.Client.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		networkConfig,
		containerName)

	return &cntr, err
}
