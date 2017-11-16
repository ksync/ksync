package ksync

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/docker"
)

// ServiceList is a list of services.
type ServiceList struct {
	Items []*Service
}

// ServiceListOptions is query options for ServiceList
type ServiceListOptions struct {
	Name string
}

// AllServices creates a ServiceList containing all the running services.
func AllServices() (*ServiceList, error) {
	list := &ServiceList{}

	err := list.Update(ServiceListOptions{})
	if err != nil { // nolint: megacheck
		return nil, err
	}
	return list, nil
}

// Update looks at the locally running containers and updates the list based on
// that state.
func (s *ServiceList) Update(opts ServiceListOptions) error {
	args := filters.NewArgs()
	args.Add("label", "heritage=ksync")
	args.Add("label", "service=true")
	if opts.Name != "" {
		args.Add("label", fmt.Sprintf("name=%s", opts.Name))
	}

	cntrs, err := docker.Client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: args,
		},
	)

	if err != nil {
		return errors.Wrap(err, "could not get container list from docker.")
	}

	for _, cntr := range cntrs {
		service := &Service{
			Name: cntr.Labels["name"],
			RemoteContainer: &RemoteContainer{
				PodName:  cntr.Labels["pod"],
				Name:     cntr.Labels["container"],
				NodeName: cntr.Labels["node"],
			},
		}
		s.Items = append(s.Items, service)
		log.WithFields(service.Fields()).Debug("found service")
	}

	return nil
}

// Stop takes all the services in a list and stops them.
func (s *ServiceList) Stop() error {
	for _, service := range s.Items {
		if err := service.Stop(); err != nil {
			return err
		}
	}

	return nil
}

// StopByName stops a specific service by pod name.
func (s *ServiceList) StopByName(name string) error {
	for _, service := range s.Items {
		if service.RemoteContainer.PodName == name {
			return service.Stop()
		}
	}

	return fmt.Errorf("no service to stop")
}

// Clean looks for running services that are no longer needed.
func (s *ServiceList) Clean() error {
	for _, service := range s.Items {
		remove, err := service.ShouldStop()
		if err != nil {
			return err
		}

		if remove {
			return service.Stop()
		}
	}

	return nil
}
