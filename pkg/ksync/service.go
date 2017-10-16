package ksync

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	apiclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Name      string
	Container *Container `structs:"-"`
	client    *apiclient.Client
	image     string // TODO: make this configurable
	Spec      *Spec  `structs:"-"`
}

type ServiceStatus struct {
	ID        string
	Status    string
	Running   bool
	StartedAt string
}

func (s *ServiceStatus) String() string {
	return YamlString(s)
}

func (s *ServiceStatus) Fields() log.Fields {
	return StructFields(s)
}

// NewService constructs a Service to manage and run local syncs from.
func NewService(name string, cntr *Container, spec *Spec) (*Service, error) {
	cli, err := apiclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create new service")
	}

	log.Debug("docker client created")

	return &Service{
		Name:      name,
		Container: cntr,
		client:    cli,
		image:     "busybox",
		Spec:      spec}, nil
}

func (s *Service) String() string {
	return YamlString(s)
}

func (s *Service) Fields() log.Fields {
	return StructFields(s)
}

func (s *Service) containerName() string {
	return fmt.Sprintf("%s-%s", s.Name, s.Container.PodName)
}

// TODO: pull image for users.
func (s *Service) create() (*container.ContainerCreateCreatedBody, error) {
	cntr, err := s.client.ContainerCreate(
		context.Background(),
		&container.Config{
			Cmd:   []string{"/bin/sh", "-c", "while true; do sleep 100; done"},
			Image: "busybox",
			Labels: map[string]string{
				"name":       s.Name,
				"pod":        s.Container.PodName,
				"container":  s.Container.Name,
				"node":       s.Container.NodeName,
				"localPath":  s.Spec.LocalPath,
				"remotePath": s.Spec.RemotePath,
			},
			// Volumes
		},
		&container.HostConfig{
			RestartPolicy: container.RestartPolicy{Name: "on-failure"},
		},
		&network.NetworkingConfig{},
		s.containerName())

	if err != nil {
		if !apiclient.IsErrImageNotFound(err) {
			return nil, ErrorOut("could not create", err, s)
		}

		return nil, fmt.Errorf("run `docker pull %s`", s.image)
	}

	return &cntr, nil
}

func (s *Service) Start() error {
	cntr, err := s.create()
	if err != nil {
		return err
	}

	log.WithFields(MergeFields(s.Fields(), log.Fields{
		"id": cntr.ID,
	})).Debug("container created")

	if err := s.client.ContainerStart(
		context.Background(),
		cntr.ID,
		types.ContainerStartOptions{}); err != nil {
		return ErrorOut("could not start", err, s)
	}

	log.WithFields(MergeFields(s.Fields(), log.Fields{
		"id": cntr.ID,
	})).Debug("container started")

	return nil
}

func (s *Service) Stop() error {
	cntr, err := s.Status()
	if err != nil {
		return err
	}

	if !cntr.Running {
		return fmt.Errorf("must start before you can stop: %s", s.containerName())
	}

	if err := s.client.ContainerRemove(
		context.Background(),
		cntr.ID,
		types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrap(err, "could not remove")
	}

	log.WithFields(MergeFields(s.Fields(), log.Fields{
		"id": cntr.ID,
	})).Debug("container removed")

	return nil
}

// Status checks to see if a service is currently running and looks at its
// status.
func (s *Service) Status() (*ServiceStatus, error) {
	cntr, err := s.client.ContainerInspect(
		context.Background(), s.containerName())
	if err != nil {
		if !apiclient.IsErrNotFound(err) {
			return nil, ErrorOut("could not get status", err, s)
		}
		return &ServiceStatus{
			ID:        "",
			Status:    "not created",
			Running:   false,
			StartedAt: "",
		}, nil
	}

	return &ServiceStatus{
		ID:        cntr.ID,
		Status:    cntr.State.Status,
		Running:   cntr.State.Running,
		StartedAt: cntr.State.StartedAt}, nil
}
