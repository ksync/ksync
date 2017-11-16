package ksync

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
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
	Name    string
	User    string
	CfgPath string

	// Kubernetes config
	Namespace   string
	Context     string
	KubeCfgPath string

	// RemoteContainer Locator
	// TODO: use a locator instead?
	Container string
	Pod       string
	Selector  string

	// File config
	LocalPath  string
	RemotePath string

	// Reload related options
	Reload bool

	stopWatching chan bool
}

func (s *Spec) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Spec) Fields() log.Fields {
	return debug.StructFields(s)
}

// Status returns the current status of a spec.
// TODO: this requires a lot more thought and effort, status is complex.
func (s *Spec) Status() (string, error) {
	status := "inactive"
	// TODO: this is super naive and should be handled differently
	list, err := s.ServiceList()
	if err != nil {
		return "", err
	}

	if len(list.Items) == 0 {
		return status, nil
	}

	var statuses []string
	for _, service := range list.Items {
		status, err := service.Status()
		if err != nil {
			return "", err
		}

		statuses = append(statuses, status.Status)
	}

	return strings.Join(statuses, ", "), nil
}

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

	// User is of format int:int
	if result, _ := regexp.MatchString("[0-9]+:[0-9]+", s.User); !result {
		return fmt.Errorf("user must be of format uid:gid, got: %s", s.User)
	}

	return nil
}

// Watch monitors the remote status of this spec.
func (s *Spec) Watch() error {
	if s.stopWatching != nil {
		log.WithFields(s.Fields()).Debug("already watching")
		return nil
	}

	opts := metav1.ListOptions{}
	opts.LabelSelector = s.Selector
	watcher, err := kubeClient.CoreV1().Pods(namespace).Watch(opts)
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
		return s.stopPod(pod)
	}

	if pod.Status.Phase == v1.PodRunning {
		return s.Start()
	}

	return nil
}

func (s *Spec) stopPod(pod *v1.Pod) error {
	log.WithFields(s.Fields()).Debug("stopping service")

	list, err := s.ServiceList()
	if err != nil {
		return err
	}

	return list.StopByName(pod.Name)
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
		s.stopWatching <- true
	}
	return s.Stop()
}

// ServiceList gets all the running services scoped to this specific spec.
func (s *Spec) ServiceList() (*ServiceList, error) {
	list := &ServiceList{}
	if err := list.Update(ServiceListOptions{Name: s.Name}); err != nil {
		return nil, err
	}

	return list, nil
}

// Start runs a service for every matching remote container.
func (s *Spec) Start() error {
	containers, err := GetRemoteContainers(s.Pod, s.Selector, s.Container)
	if err != nil {
		return debug.ErrorOut("unable to get container list", err, nil)
	}

	if len(containers) == 0 {
		log.WithFields(s.Fields()).Debug("no matching running containers.")
	}

	for _, cntr := range containers {
		if err := NewService(s.Name, cntr, s).Start(); err != nil &&
			!IsServiceRunning(err) {
			return err
		}

		log.WithFields(cntr.Fields()).Debug("container running")
	}

	return nil
}

// Stop will stop every service for this spec that is running locally.
func (s *Spec) Stop() error {
	list, err := s.ServiceList()
	if err != nil {
		return err
	}
	return list.Stop()
}
