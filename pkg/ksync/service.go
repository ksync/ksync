package ksync

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/debug"
)

// Service reflects a sync that can be run in the background.
type Service struct {
	RemoteContainer *RemoteContainer `structs:"-"`
	Spec            *Spec            `structs:"-"`

	mirror *Mirror
}

// ServiceStatus contains the current status of a given service
type ServiceStatus string

// NewService constructs a Service to manage and run local syncs from.
func NewService(cntr *RemoteContainer, spec *Spec) *Service {
	return &Service{
		RemoteContainer: cntr,
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

// Start runs this service in the background.
func (s *Service) Start() error {
	if s.mirror != nil {
		return fmt.Errorf("already running")
	}

	if err := s.RemoteContainer.RestartMirror(); err != nil {
		return err
	}

	s.mirror = &Mirror{
		SpecName:        s.Spec.Name,
		RemoteContainer: s.RemoteContainer,
		Reload:          s.Spec.Reload,
		LocalPath:       s.Spec.LocalPath,
		RemotePath:      s.Spec.RemotePath,
	}

	if err := s.mirror.Run(); err != nil { // nolint: megacheck
		return err
	}

	return nil
}

// Stop halts a service that has been running in the background.
func (s *Service) Stop() error {
	log.WithFields(s.Fields()).Debug("stopping service")
	return s.mirror.Stop()
}

// Status checks to see if a service is currently running and looks at its
// status.
func (s *Service) Status() (ServiceStatus, error) {
	return "", nil
}
