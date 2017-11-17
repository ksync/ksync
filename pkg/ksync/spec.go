package ksync

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/vapor-ware/ksync/pkg/debug"
)

var (
	specEquivalenceFields = []string{"Name"}
)

// Spec is all the configuration required to setup a sync from a local directory
// to a remote directory in a specific remote container.
type Spec struct {
	// Local config
	Name string

	// RemoteContainer Locator
	ContainerName string
	Pod           string
	Selector      string
	Namespace     string

	// File config
	LocalPath  string
	RemotePath string

	// Reload related options
	Reload bool

	services     *ServiceList
	stopWatching chan bool
}

func (s *Spec) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Spec) Fields() log.Fields {
	return debug.StructFields(s)
}

// TODO: implement status now that mirror is being run from inside watch.
// Status returns the current status of a spec.
// TODO: this requires a lot more thought and effort, status is complex.
// func (s *Spec) Status() (string, error) {
// 	status := "inactive"
// 	// TODO: this is super naive and should be handled differently

// 	if len(s.services.Items) == 0 {
// 		return status, nil
// 	}

// 	var statuses []string
// 	for _, service := range s.services.Items {
// 		status, err := service.Status()
// 		if err != nil {
// 			return "", err
// 		}

// 		statuses = append(statuses, status.Status)
// 	}

// 	return strings.Join(statuses, ", "), nil
// }

// IsValid returns an error if the spec is not valid.
func (s *Spec) IsValid() error {

	// Cannot sync files, must do directories.
	fstat, err := os.Stat(s.LocalPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if fstat != nil && !fstat.IsDir() {
		return fmt.Errorf("local path cannot be a single file, please use a directory")
	}

	return nil
}

// Watch monitors the remote status of this spec.
func (s *Spec) Watch() error {
	if s.stopWatching != nil {
		log.WithFields(s.Fields()).Debug("already watching")
		return nil
	}

	s.services = &ServiceList{}

	opts := metav1.ListOptions{}
	opts.LabelSelector = s.Selector
	watcher, err := kubeClient.CoreV1().Pods(s.Namespace).Watch(opts)
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
		return s.addService(pod)
	}

	return nil
}

func (s *Spec) addService(pod *v1.Pod) error {
	if err := s.services.Add(pod, s); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (s *Spec) cleanService(pod *v1.Pod) error {
	service := s.services.Pop(pod.Name)
	if service == nil {
		log.WithFields(s.Fields()).Debug("service not found")
		return nil
	}

	if err := service.Stop(); err != nil {
		return err
	}

	return nil
}

// Equivalence returns a set of fields that can be used to compare specs for
// equivalence via. reflect.DeepEqual.
func (s *Spec) Equivalence() map[string]interface{} {
	vals := structs.Map(s)
	for _, k := range specEquivalenceFields {
		delete(vals, k)
	}
	return vals
}

// Cleanup will remove anything running in the background, meant to be used when
// a spec is deleted.
func (s *Spec) Cleanup() error {
	if s.stopWatching != nil {
		close(s.stopWatching)
	}
	return s.services.Stop()
}
