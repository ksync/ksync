package ksync

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RadarInstance is the remote server component of ksync.
type RadarInstance struct {
	namespace  string
	name       string
	labels     map[string]string
	mirrorPort int32
	radarPort  int32
}

func (r *RadarInstance) String() string {
	return YamlString(r)
}

// Fields returns a set of structured fields for logging.
func (r *RadarInstance) Fields() log.Fields {
	return StructFields(r)
}

// NewRadarInstance constructs a RadarInstance to track the remote status.
// TODO: make namespace, name?, service account configurable
func NewRadarInstance() *RadarInstance {
	return &RadarInstance{
		namespace:  "kube-system",
		name:       "ksync-radar",
		mirrorPort: 49172,
		radarPort:  40321,
		labels: map[string]string{
			"name": "ksync-radar",
			"app":  "radar",
		},
	}
}

// Run starts (or upgrades) radar on the remote cluster.
// TODO: spin up on demand
// TODO: wait for ready
func (r *RadarInstance) Run(upgrade bool) error {
	fn := kubeClient.DaemonSets(r.namespace).Create

	if upgrade {
		fn = kubeClient.DaemonSets(r.namespace).Update
	}

	_, err := fn(r.daemonSet())

	// TODO: need better error
	if err != nil {
		return err
	}

	log.WithFields(MergeFields(r.Fields(), log.Fields{
		"upgrade": upgrade,
	})).Debug("started DaemonSet")

	return nil
}

// TODO: add TLS
// TODO: add grpc_retry?
func (r *RadarInstance) opts() []grpc.DialOption {
	return append([]grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
		grpc.WithBlock(),
		// TODO: add client side tracing
	}, grpc.WithInsecure())
}

func (r *RadarInstance) podName(nodeName string) (string, error) {
	// TODO: error handling for nodes that don't exist.
	pods, err := kubeClient.CoreV1().Pods(r.namespace).List(
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", r.labels["app"]),
			FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
		})

	if err != nil {
		return "", nil
	}

	// TODO: provide a better error here, explain to users how to fix it.
	if len(pods.Items) != 1 {
		return "", fmt.Errorf(
			"unexpected result looking up radar pod (count:%s) (node:%s)",
			len(pods.Items),
			nodeName)
	}

	return pods.Items[0].Name, nil
}

func (r *RadarInstance) connection(nodeName string, port int32) (int32, error) {
	podName, err := r.podName(nodeName)
	if err != nil {
		return 0, ErrorOut("cannot get pod name", err, r)
	}

	tun, err := NewTunnel(r.namespace, podName, r.radarPort)
	if err != nil {
		return 0, ErrorOut("unable to create tunnel", err, r)
	}

	if err := tun.Start(); err != nil {
		return 0, ErrorOut("unable to start tunnel", err, r)
	}

	return tun.LocalPort, nil
}

// RadarConnection creates a new gRPC connection to a radar instance running on
// the specified node.
func (r *RadarInstance) RadarConnection(nodeName string) (*grpc.ClientConn, error) {
	localPort, err := r.connection(nodeName, r.radarPort)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(fmt.Sprintf("127.0.0.1:%d", localPort), r.opts()...)
}

// MirrorConnection creates a tunnel to the remote mirror instance running on
// the specified node.
func (r *RadarInstance) MirrorConnection(nodeName string) (int32, error) {
	return r.connection(nodeName, r.mirrorPort)
}
