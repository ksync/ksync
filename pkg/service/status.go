package service

import (
	"context"

	apiclient "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/docker"
)

// Status is the status of a specific service.
type Status struct {
	ID        string
	Status    string
	Running   bool
	StartedAt string
}

func (s *Status) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Status) Fields() log.Fields {
	return debug.StructFields(s)
}

// GetStatus returns a Status based on container name.
func GetStatus(name string) (*Status, error) {
	cntr, err := docker.Client.ContainerInspect(context.Background(), name)
	if err != nil {
		if !apiclient.IsErrNotFound(err) {
			return nil, err
		}
		return &Status{
			ID:        "",
			Status:    "not created",
			Running:   false,
			StartedAt: "",
		}, nil
	}

	return &Status{
		ID:        cntr.ID,
		Status:    cntr.State.Status,
		Running:   cntr.State.Running,
		StartedAt: cntr.State.StartedAt}, nil
}
