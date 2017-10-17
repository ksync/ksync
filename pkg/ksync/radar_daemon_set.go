package ksync

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func (r *RadarInstance) daemonSet() *v1beta1.DaemonSet {
	return &v1beta1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			// TODO: configurable
			Namespace: r.namespace,
			Name:      r.name,
			Labels:    r.labels,
		},
		Spec: v1beta1.DaemonSetSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: r.labels,
					Annotations: map[string]string{
						// TODO: this should only be set on --upgrade --force
						"forceUpdate": fmt.Sprint(time.Now().Unix()),
						// TODO: set inotify sysctl high en
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: r.name,
							// TODO: configurable
							Image:           "gcr.io/elated-embassy-152022/ksync/radar:canary",
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{ContainerPort: r.radarPort, Name: "grpc"},
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
								{ContainerPort: r.mirrorPort, Name: "grpc"},
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
}
