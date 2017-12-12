package ksync

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/vapor-ware/ksync/pkg/debug"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// SpecStatus is the status of a spec
type SpecStatus string

// See docs/spec-lifecycle.png
const (
	SpecWaiting SpecStatus = "waiting"
	SpecRunning SpecStatus = "running"
)

// Spec is all the configuration required to setup a sync from a local directory
// to a remote directory in a specific remote container.
type Spec struct {
	Details  *SpecDetails
	Services *ServiceList `structs:"-"`

	Status SpecStatus

	stopWatching chan bool
}

func (s *Spec) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Spec) Fields() log.Fields {
	return s.Details.Fields()
}

// Message is used to serialize over gRPC
func (s *Spec) Message() (*pb.Spec, error) {
	details, err := s.Details.Message()
	if err != nil {
		return nil, err
	}

	services, err := s.Services.Message()
	if err != nil {
		return nil, err
	}

	return &pb.Spec{
		Details:  details,
		Services: services,
		Status:   string(s.Status),
	}, nil
}

// NewSpec is a constructor for Specs
func NewSpec(details *SpecDetails) *Spec {
	return &Spec{
		Details:  details,
		Services: NewServiceList(),
		Status:   SpecWaiting,
	}
}

// Watch monitors the remote status of this spec.
func (s *Spec) Watch() error {
	if s.stopWatching != nil {
		log.WithFields(s.Fields()).Debug("already watching")
		return nil
	}

	opts := metav1.ListOptions{}
	opts.LabelSelector = s.Details.Selector
	watcher, err := kubeClient.CoreV1().Pods(s.Details.Namespace).Watch(opts)
	if err != nil {
		return err
	}

	log.WithFields(s.Fields()).Debug("watching for updates")

	s.stopWatching = make(chan bool)
	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-s.stopWatching:
				log.WithFields(s.Fields()).Debug("stopping watch")
				return

			case event := <-watcher.ResultChan():
				if event.Object == nil {
					continue
				}

				if err := s.handleEvent(event); err != nil {
					log.WithFields(s.Fields()).Error(err)
				}
			}
		}
	}()

	return nil
}

func (s *Spec) handleEvent(event watch.Event) error {
	pod := event.Object.(*v1.Pod)
	if event.Type != watch.Modified && event.Type != watch.Added {
		return nil
	}

	log.WithFields(log.Fields{
		"type":    event.Type,
		"name":    pod.Name,
		"status":  pod.Status.Phase,
		"deleted": pod.DeletionTimestamp != nil,
	}).Debug("new event")

	if pod.DeletionTimestamp != nil {
		return s.cleanService(pod)
	}

	if pod.Status.Phase == v1.PodRunning {
		s.Status = SpecRunning
		return s.addService(pod)
	}

	return nil
}

func (s *Spec) addService(pod *v1.Pod) error {
	if err := s.Services.Add(pod, s.Details); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (s *Spec) cleanService(pod *v1.Pod) error {
	service := s.Services.Pop(pod.Name)
	if service == nil {
		log.WithFields(s.Fields()).Debug("service not found")
		return nil
	}

	if err := service.Stop(); err != nil {
		return err
	}

	if len(s.Services.Items) == 0 {
		s.Status = SpecWaiting
	}

	return nil
}

// Cleanup will remove anything running in the background, meant to be used when
// a spec is deleted.
func (s *Spec) Cleanup() error {
	if s.stopWatching != nil {
		close(s.stopWatching)
	}
	return s.Services.Stop()
}
