package ksync

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// ServiceList is a list of services.
type ServiceList struct {
	Items []*Service
}

// Add takes a pod/spec, creates a new service, adds it to the list and starts it.
func (s *ServiceList) Add(pod *v1.Pod, spec *Spec) error {
	cntr, err := NewRemoteContainer(pod, spec.ContainerName)
	if err != nil {
		return err
	}

	service := NewService(cntr, spec)
	if s.Has(service) {
		return &errors.StatusError{
			ErrStatus: metav1.Status{
				Status:  metav1.StatusFailure,
				Reason:  metav1.StatusReasonAlreadyExists,
				Message: fmt.Sprintf("service for %s already exists", pod.Name),
			}}
	}

	s.Items = append(s.Items, service)

	log.WithFields(service.Fields()).Debug("added service")

	return service.Start()
}

// Has checks for equivalence between the existing services.
func (s *ServiceList) Has(target *Service) bool {
	for _, service := range s.Items {
		if reflect.DeepEqual(service.RemoteContainer, target.RemoteContainer) {
			return true
		}
	}

	return false
}

// Pop fetches a service by pod name and removes it from the list.
func (s *ServiceList) Pop(podName string) *Service {
	for i, service := range s.Items {
		if service.RemoteContainer.PodName == podName {
			s.Items[i] = s.Items[len(s.Items)-1]
			s.Items = s.Items[:len(s.Items)-1]
			return service
		}
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
