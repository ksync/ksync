package ksync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/service"
)

var (
	imageName string
)

// SetImage sets the package-wide image to use for launching tasks
// (both local and remote).
func SetImage(name string) {
	imageName = name
}

// Service reflects a sync that can be run in the background.
type Service struct {
	Name            string
	RemoteContainer *RemoteContainer `structs:"-"`
	Spec            *Spec            `structs:"-"`
}

// NewService constructs a Service to manage and run local syncs from.
func NewService(name string, cntr *RemoteContainer, spec *Spec) *Service {
	return &Service{
		Name:            name,
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

func (s *Service) containerName() string {
	return fmt.Sprintf("%s-%s", s.Name, s.RemoteContainer.PodName)
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
	return service.Start(
		&container.Config{
			// TODO: make most of these options configurable.
			// TODO: missing context
			Cmd: []string{
				"/ksync",
				"--log-level=debug",
				fmt.Sprintf("--context=%s", s.Spec.Context),
				"run",
				fmt.Sprintf("--pod=%s", s.RemoteContainer.PodName),
				fmt.Sprintf("--container=%s", s.RemoteContainer.Name),
				s.Spec.LocalPath,
				s.Spec.RemotePath,
			},
			Image: imageName,
			Labels: map[string]string{
				"name":       s.Name,
				"specName":   s.Spec.Name,
				"pod":        s.RemoteContainer.PodName,
				"container":  s.RemoteContainer.Name,
				"node":       s.RemoteContainer.NodeName,
				"localPath":  s.Spec.LocalPath,
				"remotePath": s.Spec.RemotePath,
				"heritage":   "ksync",
				"service":    "true",
			},
			User: s.Spec.User,
			Env:  []string{"KUBECONFIG=/.kube/config"},
		},
		&container.HostConfig{
			// TODO: need to make this configurable
			Binds: []string{
				fmt.Sprintf("%s:/.kube/config", s.Spec.KubeCfgPath),
				fmt.Sprintf("%s:%s", s.Spec.LocalPath, s.Spec.LocalPath),
				fmt.Sprintf("%s:/.ksync", s.Spec.CfgPath),
			},
			RestartPolicy: container.RestartPolicy{Name: "on-failure"},
		},
		&network.NetworkingConfig{},
		s.containerName())
}

// Stop halts a service that has been running in the background.
func (s *Service) Stop() error {
	log.WithFields(s.Fields()).Debug("stopping service")
	return service.Stop(s.containerName())
}

// ShouldStop checks to see if this service should still run or not.
func (s *Service) ShouldStop() (bool, error) {
	// remote container still running
	if _, err := GetByName(
		s.RemoteContainer.PodName, s.RemoteContainer.Name); err != nil {
		if apiErrors.IsNotFound(err) {
			return true, nil
		}

		return false, err
	}

	list := &SpecList{}
	if err := list.Update(); err != nil {
		return false, err
	}

	if !list.Has(s.Name) {
		return true, nil
	}

	return false, nil
}

// Status checks to see if a service is currently running and looks at its
// status.
func (s *Service) Status() (*service.Status, error) {
	return service.GetStatus(s.containerName())
}
