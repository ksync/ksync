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

func InitRadarOpts() {
	// TODO: add TLS
	// TODO: add grpc_retry?
	grpcOpts = append(grpcOpts, grpc.WithInsecure())
}

func NewRadarConnection(nodeName string) (*grpc.ClientConn, error) {
	tun := NewTunnel("ksync-radar-1jf5p")
	if err := tun.Start(); err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("127.0.0.1:%d", tun.LocalPort), grpcOpts...)
}
