package ksync

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	labels = map[string]string{
		"name": "ksync-radar",
		"app":  "radar",
	}

	radarNamespace = "kube-system"
	radarName      = "ksync-radar"
	radarPort      = 40321

	grpcOpts = []grpc.DialOption{
		grpc.WithTimeout(5 * time.Second),
		grpc.WithBlock(),
		// TODO: add client side tracing
	}

	// TODO: make namespace, name?, service account configurable
	radarDaemonSet = &v1beta1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			// TODO: configurable
			Namespace: radarNamespace,
			Name:      radarName,
			Labels:    labels,
		},
		Spec: v1beta1.DaemonSetSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						// TODO: this should only be set on --upgrade --force
						"forceUpdate": fmt.Sprint(time.Now().Unix()),
						// TODO: set inotify sysctl high en
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: radarName,
							// TODO: configurable
							Image:           "gcr.io/elated-embassy-152022/ksync/radar:canary",
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{ContainerPort: 40321, Name: "grpc"},
							},
							// TODO: resources
							VolumeMounts: []v1.VolumeMount{
								v1.VolumeMount{
									Name:      "dockerfs",
									MountPath: "/var/lib/docker",
								},
								v1.VolumeMount{
									Name:      "dockersock",
									MountPath: "/var/run/docker.sock",
								},
								v1.VolumeMount{
									Name:      "kubelet",
									MountPath: "/var/lib/kubelet",
								},
							},
						},
						{
							Name: "mirror",
							// TODO: configurable
							Image:           "gcr.io/elated-embassy-152022/ksync/mirror:canary",
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{ContainerPort: 49172, Name: "grpc"},
							},
							// TODO: resources
							VolumeMounts: []v1.VolumeMount{
								v1.VolumeMount{
									Name:      "dockerfs",
									MountPath: "/var/lib/docker",
								},
								v1.VolumeMount{
									Name:      "dockersock",
									MountPath: "/var/run/docker.sock",
								},
								v1.VolumeMount{
									Name:      "kubelet",
									MountPath: "/var/lib/kubelet",
								},
							},
						},
					},
					NodeSelector: map[string]string{
						"beta.kubernetes.io/os": "linux",
					},
					// TODO: add HostPathType
					Volumes: []v1.Volume{
						v1.Volume{
							Name: "dockerfs",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/var/lib/docker",
								},
							},
						},
						v1.Volume{
							Name: "dockersock",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/var/run/docker.sock",
								},
							},
						},
						v1.Volume{
							Name: "kubelet",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: "/var/lib/kubelet",
								},
							},
						},
					},
				},
			},
			UpdateStrategy: v1beta1.DaemonSetUpdateStrategy{
				Type: "RollingUpdate",
			},
		},
	}
)

// InitRadar initializes a new instance of radar (server) and deploys it
// as a DaemonSet into the cluster
// TODO: spin up on demand
// TODO: wait for ready
func InitRadar(upgrade bool) error {
	fn := KubeClient.DaemonSets(radarDaemonSet.Namespace).Create

	if upgrade {
		fn = KubeClient.DaemonSets(radarDaemonSet.Namespace).Update
	}

	_, err := fn(radarDaemonSet)

	// TODO: need better error
	if err != nil {
		return err
	}

	log.Debug("started DaemonSet")

	return nil
}

// InitRadarOpts initializes the grpc options for an instance of radar
func InitRadarOpts() {
	// TODO: add TLS
	// TODO: add grpc_retry?
	grpcOpts = append(grpcOpts, grpc.WithInsecure())
}

// radarPodName returns the pod name where the launched instance of
// radar is running
func radarPodName(nodeName string) (string, error) {
	// TODO: error handling for nodes that don't exist.
	pods, err := KubeClient.CoreV1().Pods(radarNamespace).List(
		metav1.ListOptions{
			LabelSelector: "app=radar",
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

// NewRadarConnection creates a new tunnel connection to a instance of radar
func NewRadarConnection(nodeName string) (*grpc.ClientConn, error) {
	tun, err := NewTunnel(nodeName, radarPort)
	if err != nil {
		return nil, err
	}
	if err := tun.Start(); err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("127.0.0.1:%d", tun.LocalPort), grpcOpts...)
}
