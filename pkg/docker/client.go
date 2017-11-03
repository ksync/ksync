package docker

import (
	apiclient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var (
	// Client is a singleton docker client that is already configured.
	Client *apiclient.Client
)

// InitClient sets up the singleton for use by the ksync package.
func InitClient() error {
	client, err := apiclient.NewEnvClient()
	if err != nil {
		return err
	}

	log.Debug("docker client created")

	Client = client

	return nil
}
