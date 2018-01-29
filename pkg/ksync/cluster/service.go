package cluster

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/vapor-ware/ksync/pkg/debug"
)

var (
	// ImageName is the docker image to use for running the cluster service.
	ImageName = "vaporio/ksync"
)

// SetImage sets the package-wide image to used for running the cluster service.
func SetImage(name string) {
	ImageName = name
}

// Service is the remote server component of ksync.
type Service struct {
	Namespace string
	name      string
	labels    map[string]string

	RadarPort         int32
	SyncthingAPI      int32
	SyncthingListener int32
}

func (s *Service) String() string {
	return debug.YamlString(s)
}

// Fields returns a set of structured fields for logging.
func (s *Service) Fields() log.Fields {
	return debug.StructFields(s)
}

// NewRadarInstance constructs a RadarInstance to track the remote status.
// TODO: make namespace, name?, service account configurable
func NewService() *Service {
	return &Service{
		Namespace: "kube-system",
		name:      "ksync",
		labels: map[string]string{
			"name": "ksync",
			"app":  "ksync",
		},

		RadarPort:         40321,
		SyncthingAPI:      8384,
		SyncthingListener: 22000,
	}
}

// IsInstalled makes sure the cluster service has been installed.
// TODO: add version checking here.
func (s *Service) IsInstalled() (bool, error) {
	if _, err := Client.DaemonSets(s.Namespace).Get(
		s.name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (s *Service) PodName(nodeName string) (string, error) {
	// TODO: error handling for nodes that don't exist.
	pods, err := Client.CoreV1().Pods(s.Namespace).List(
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", s.labels["app"]),
			FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
		})

	if err != nil {
		return "", err
	}

	// TODO: I don't want to go and setup a whole bunch of error handling code
	// right now. The NotFound error works perfectly here for now.
	if len(pods.Items) == 0 {
		return "", &errors.StatusError{
			ErrStatus: metav1.Status{
				Status:  metav1.StatusFailure,
				Reason:  metav1.StatusReasonNotFound,
				Message: fmt.Sprintf("%s not found on %s", s.name, nodeName),
			}}
	}

	if len(pods.Items) > 1 {
		return "", fmt.Errorf(
			"unexpected result looking up radar pod (count:%d) (node:%s)",
			len(pods.Items),
			nodeName)
	}

	return pods.Items[0].Name, nil
}

// IsHealthy verifies the target node is running the service container and it
// is not scheduled for deletion.
func (s *Service) IsHealthy(nodeName string) (bool, error) {
	log.WithFields(log.Fields{
		"nodeName": nodeName,
	}).Debug("checking to see if radar is ready")

	installed, err := s.IsInstalled()
	if err != nil {
		return false, err
	}

	if !installed {
		return false, fmt.Errorf("radar is not installed")
	}

	podName, err := s.PodName(nodeName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, debug.ErrorOut("cannot get pod name", err, s)
		}

		return false, nil
	}

	log.WithFields(debug.MergeFields(s.Fields(), log.Fields{
		"nodeName": nodeName,
		"podName":  podName,
	})).Debug("found pod name")

	pod, err := Client.CoreV1().Pods(s.Namespace).Get(
		podName, metav1.GetOptions{})
	if err != nil {
		return false, debug.ErrorOut("cannot get pod details", err, s)
	}

	log.WithFields(log.Fields{
		"podName":  pod.Name,
		"nodeName": nodeName,
		"status":   pod.Status.Phase,
	}).Debug("found pod")

	if pod.Status.Phase != v1.PodRunning || pod.DeletionTimestamp != nil {
		return false, nil
	}

	return true, nil
}

// NodeNames returns a list of all the nodes the cluster service is running on.
func (s *Service) NodeNames() ([]string, error) {
	result := []string{}

	opts := metav1.ListOptions{}
	opts.LabelSelector = "app=radar"
	pods, err := Client.CoreV1().Pods(s.Namespace).List(opts)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"count": len(pods.Items),
	}).Debug("radar nodes")

	for _, pod := range pods.Items {
		result = append(result, pod.Spec.NodeName)
	}

	return result, nil
}

// Run starts (or upgrades) the cluster service.
// TODO: spin up on demand
// TODO: wait for ready
func (s *Service) Run(upgrade bool) error {
	daemonSets := Client.DaemonSets(s.Namespace)

	if _, err := daemonSets.Create(s.daemonSet()); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := daemonSets.Update(s.daemonSet()); err != nil {
			return err
		}
	}

	log.WithFields(debug.MergeFields(s.Fields(), log.Fields{
		"upgrade": upgrade,
	})).Debug("started DaemonSet")

	return nil
}

func (s *Service) Remove() error {
	daemonSets := Client.DaemonSets(s.Namespace)

	if err := daemonSets.Delete(s.name, &metav1.DeleteOptions{}); err != nil {
		return err
	}

	log.WithFields(debug.MergeFields(s.Fields(), log.Fields{
		"Name": s.name,
	})).Debug("Removed DaemonSet")

	return nil
}
