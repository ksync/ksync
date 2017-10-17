package ksync

import (
	apiclient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var (
	dockerClient *apiclient.Client
)

func InitDockerClient() error {
	client, err := apiclient.NewEnvClient()
	if err != nil {
		return err
	}

	log.Debug("docker client created")

	dockerClient = client

	return nil
}
