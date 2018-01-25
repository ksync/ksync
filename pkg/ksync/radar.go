package ksync

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/vapor-ware/ksync/pkg/debug"
)

var (
	maxReadyRetries = uint64(10)
	// RadarImageName is the docker image to use for running radar.
	RadarImageName = "vaporio/ksync"
)

// SetImage sets the package-wide image to use for launching tasks
// (both local and remote).
func SetImage(name string) {
	RadarImageName = name
}

// RadarInstance is the remote server component of ksync.
type RadarInstance struct {
	namespace         string
	name              string
	labels            map[string]string
	radarPort         int32
	syncthingAPI      int32
	syncthingListener int32
}

func (r *RadarInstance) String() string {
	return debug.YamlString(r)
}

// Fields returns a set of structured fields for logging.
func (r *RadarInstance) Fields() log.Fields {
	return debug.StructFields(r)
}

// NewRadarInstance constructs a RadarInstance to track the remote status.
// TODO: make namespace, name?, service account configurable
func NewRadarInstance() *RadarInstance {
	return &RadarInstance{
		namespace:         "kube-system",
		name:              "ksync",
		syncthingAPI:      8384,
		syncthingListener: 22000,
		radarPort:         40321,
		labels: map[string]string{
			"name": "ksync-radar",
			"app":  "radar",
		},
	}
}

// IsInstalled makes sure radar has been submitted to the remote cluster.
// TODO: add version checking here.
func (r *RadarInstance) IsInstalled() (bool, error) {
	if _, err := kubeClient.DaemonSets(r.namespace).Get(
		r.name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// IsHealthy verifies the target node is running radar and it is not scheduled
// for deletion.
func (r *RadarInstance) IsHealthy(nodeName string) (bool, error) {
	log.WithFields(log.Fields{
		"nodeName": nodeName,
	}).Debug("checking to see if radar is ready")

	installed, err := r.IsInstalled()
	if err != nil {
		return false, err
	}

	if !installed {
		return false, fmt.Errorf("radar is not installed")
	}

	podName, err := r.podName(nodeName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, debug.ErrorOut("cannot get pod name", err, r)
		}

		return false, nil
	}

	log.WithFields(MergeFields(r.Fields(), log.Fields{
		"nodeName": nodeName,
		"podName":  podName,
	})).Debug("found pod name")

	pod, err := kubeClient.CoreV1().Pods(r.namespace).Get(
		podName, metav1.GetOptions{})
	if err != nil {
		return false, debug.ErrorOut("cannot get pod details", err, r)
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

// NodeNames returns a list of all the nodes radar is currently running on.
func (r *RadarInstance) NodeNames() ([]string, error) {
	result := []string{}

	opts := metav1.ListOptions{}
	opts.LabelSelector = "app=radar"
	pods, err := kubeClient.CoreV1().Pods(r.namespace).List(opts)
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

// Run starts (or upgrades) radar on the remote cluster.
// TODO: spin up on demand
// TODO: wait for ready
func (r *RadarInstance) Run(upgrade bool) error {
	daemonSets := kubeClient.DaemonSets(r.namespace)

	if _, err := daemonSets.Create(r.daemonSet()); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := daemonSets.Update(r.daemonSet()); err != nil {
			return err
		}
	}

	log.WithFields(MergeFields(r.Fields(), log.Fields{
		"upgrade": upgrade,
	})).Debug("started DaemonSet")

	return nil
}

// TODO: add TLS
// TODO: add grpc_retry?
func (r *RadarInstance) opts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
}

func (r *RadarInstance) podName(nodeName string) (string, error) {
	// TODO: error handling for nodes that don't exist.
	pods, err := kubeClient.CoreV1().Pods(r.namespace).List(
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", r.labels["app"]),
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
				Message: fmt.Sprintf("%s not found on %s", r.name, nodeName),
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

func (r *RadarInstance) waitForHealthy(nodeName string) error {
	test := func() error {
		ready, err := r.IsHealthy(nodeName)
		if err != nil {
			return backoff.Permanent(err)
		}

		if !ready {
			return fmt.Errorf("radar on %s not ready", nodeName)
		}

		return nil
	}

	return backoff.Retry(
		test,
		backoff.WithMaxTries(backoff.NewExponentialBackOff(), maxReadyRetries))
}

func (r *RadarInstance) connection(nodeName string, port int32) (int32, error) {
	if err := r.waitForHealthy(nodeName); err != nil {
		return 0, err
	}

	podName, err := r.podName(nodeName)
	if err != nil {
		return 0, debug.ErrorOut("cannot get pod name", err, r)
	}

	tun, err := NewTunnel(r.namespace, podName, port)
	if err != nil {
		return 0, debug.ErrorOut("unable to create tunnel", err, r)
	}

	if err := tun.Start(); err != nil {
		return 0, debug.ErrorOut("unable to start tunnel", err, r)
	}

	return tun.LocalPort, nil
}

// RadarConnection creates a new gRPC connection to a radar instance running on
// the specified node.
func (r *RadarInstance) RadarConnection(nodeName string) (*grpc.ClientConn, error) {
	localPort, err := r.connection(nodeName, r.radarPort)
	if err != nil {
		return nil, debug.ErrorLocation(err)
	}

	return grpc.Dial(fmt.Sprintf("127.0.0.1:%d", localPort), r.opts()...)
}

// SyncthingConnection creates a tunnel to the remote syncthing instance running on
// the specified node.
func (r *RadarInstance) SyncthingConnection(nodeName string) (int32, int32, error) {
	apiPort, err := r.connection(nodeName, r.syncthingAPI)
	if err != nil {
		return 0, 0, err
	}

	listenerPort, err := r.connection(nodeName, r.syncthingListener)
	if err != nil {
		return 0, 0, err
	}

	return apiPort, listenerPort, nil
}
