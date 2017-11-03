package ksync

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/vapor-ware/ksync/pkg/debug"
	"github.com/vapor-ware/ksync/pkg/docker"
)

// ServiceList is a list of services.
type ServiceList struct {
	Items []*Service
}

// GetServices creates a ServiceList containing all the running services.
func GetServices() *ServiceList {
	list := &ServiceList{}

	list.populate() // nolint: errcheck

	return list
}

func (s *ServiceList) populate() error {
	args := filters.NewArgs()
	args.Add("label", "heritage=ksync")
	args.Add("label", "service=true")

	cntrs, err := docker.Client.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Filters: args,
		},
	)

	// TODO: is this even possible?
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

// Normalize starts services for any specs that don't have ones and stops the
// services that are no longer required.
func (s *ServiceList) Normalize() error {
	specs, _ := AllSpecs()

	if err := s.compact(specs); err != nil {
		return err
	}

	return s.update(specs)
}

// Filter takes a name and returns a new instance of ServiceList populated with
// elements that have that name.
func (s *ServiceList) Filter(name string) *ServiceList {
	list := &ServiceList{}
	for _, service := range s.Items {
		if service.Name == name {
			list.Items = append(list.Items, service)
		}
	}

	return list
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

func (s *ServiceList) compact(specs *SpecMap) error {
	for _, service := range s.Items {
		if _, ok := specs.Items[service.Name]; ok {
			continue
		}

		if err := service.Stop(); err != nil {
			return errors.Wrap(
				err, "unable to stop service that is no longer needed.")
		}
	}

	return nil
}

func (s *ServiceList) update(specs *SpecMap) error {
	for name, spec := range specs.Items {
		containerList, err := GetRemoteContainers(
			spec.Pod, spec.Selector, spec.Container)
		if err != nil {
			return debug.ErrorOut("unable to get container list", err, nil)
		}

		if len(containerList) == 0 {
			log.WithFields(spec.Fields()).Debug("no matching running containers.")

			if err := s.Filter(name).Stop(); err != nil {
				return err
			}
			continue
		}

		// TODO: should this be on its own?
		for _, cntr := range containerList {
			if err := NewService(name, cntr, spec).Start(); err != nil {

				if IsServiceRunning(err) {
					log.WithFields(MergeFields(
						cntr.Fields(), spec.Fields())).Debug("already running")
					continue
				}

				return err
			}
		}
	}

	return nil
}
