package ksync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	apiclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/docker"
	"github.com/vapor-ware/ksync/pkg/service"
)

// Service reflects a sync that can be run in the background.
type Service struct {
	Name            string
	RemoteContainer *RemoteContainer `structs:"-"`
	Spec            *Spec            `structs:"-"`

	image string // TODO: make this configurable
}

// NewService constructs a Service to manage and run local syncs from.
func NewService(name string, cntr *RemoteContainer, spec *Spec) *Service {
	return &Service{
		Name:            name,
		RemoteContainer: cntr,
		image:           "gcr.io/elated-embassy-152022/ksync/ksync:canary",
		Spec:            spec,
	}
}

func (s *Service) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Service) Fields() log.Fields {
	return debug.StructFields(s)
}

func (s *Service) containerName() string {
	return fmt.Sprintf("%s-%s", s.Name, s.RemoteContainer.PodName)
}

func (s *Service) create() (*container.ContainerCreateCreatedBody, error) {
	cntr, err := docker.Client.ContainerCreate(
		context.Background(),
		&container.Config{
			// TODO: make most of these options configurable.
			// TODO: missing context
			Cmd: []string{
				"/ksync",
				"--log-level=debug",
				"run",
				fmt.Sprintf("--pod=%s", s.RemoteContainer.PodName),
				fmt.Sprintf("--container=%s", s.RemoteContainer.Name),
				s.Spec.LocalPath,
				s.Spec.RemotePath,
			},
			Image: s.image,
			Labels: map[string]string{
				"name":       s.Name,
				"pod":        s.RemoteContainer.PodName,
				"container":  s.RemoteContainer.Name,
				"node":       s.RemoteContainer.NodeName,
				"localPath":  s.Spec.LocalPath,
				"remotePath": s.Spec.RemotePath,
				"heritage":   "ksync",
				"service":    "true",
			},
			User: s.Spec.User,
		},
		&container.HostConfig{
			// TODO: need to make this configurable
			Binds: []string{
				fmt.Sprintf("%s:/.kube/config", s.Spec.KubeCfgPath),
				fmt.Sprintf("%s:%s", s.Spec.LocalPath, s.Spec.LocalPath),
			},
			RestartPolicy: container.RestartPolicy{Name: "on-failure"},
		},
		&network.NetworkingConfig{},
		s.containerName())

	return &cntr, err
}

// Start runs a service in the background.
// TODO: pull image for users.
func (s *Service) Start() error {
	status, err := s.Status()
	if err != nil {
		return err
	}

	if status.Running {
		return serviceRunningError{
			service: s,
		}
	}

	// Watch is run from inside a container, because of volume mounts, `/host`
	// must be prepended here.
	internalPath := filepath.Join("/host", s.Spec.LocalPath)
	if _, existErr := os.Stat(internalPath); os.IsNotExist(existErr) {
		if mkErr := os.MkdirAll(internalPath, 0755); mkErr != nil {
			return errors.Wrap(mkErr, "local path does not exist, cannot create")
		}

		// TODO: make this more than just a debug statement, it is important to the
		// user.
		log.WithFields(log.Fields{
			"path":       s.Spec.LocalPath,
			"permission": 0755,
		}).Debug("created missing local directory")
	}

	// TODO: check whether the configured container user can write to localPath

	cntr, err := s.create()
	if err != nil {
		if apiclient.IsErrImageNotFound(err) {
			return fmt.Errorf("run `docker pull %s`", s.image)
		}

		return err
	}

	log.WithFields(MergeFields(s.Fields(), log.Fields{
		"id": cntr.ID,
	})).Debug("container created")

	return service.Start(cntr)
}

// Stop halts a service that has been running in the background.
func (s *Service) Stop() error {
	return service.Stop(s.containerName())
}

// Status checks to see if a service is currently running and looks at its
// status.
func (s *Service) Status() (*service.Status, error) {
	return service.GetStatus(s.containerName())
}
