package cluster

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	policyv1beta "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type creationFunc func(bool) error

func (s *Service) creationFuncs(withPSP bool) []creationFunc {
	funcs := []creationFunc{s.createDaemonSet, s.createServiceAccount}
	if withPSP {
		funcs = append(funcs, s.createPSP, s.createClusterRole, s.createClusterRoleBinding)
	}
	return funcs
}

func (s *Service) createDaemonSet(upgrade bool) error {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.Namespace,
			Name:      s.name,
			Labels:    s.labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: s.labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: s.labels,
					Annotations: map[string]string{
						// TODO: this should only be set on --upgrade --force
						"forceUpdate": fmt.Sprint(time.Now().Unix()),
						// TODO: set inotify sysctl high en
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: s.name,
							// TODO: configurable
							Image:           ImageName,
							ImagePullPolicy: "Always",
							Command:         []string{"/radar", "--log-level=debug", "serve"},
							Env: []v1.EnvVar{
								{
									Name: "RADAR_POD_NAME",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							Ports: []v1.ContainerPort{
								{ContainerPort: s.RadarPort, Name: "grpc"},
							},
							// TODO: resources
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "dockersock",
									MountPath: viper.GetString("docker-socket"),
								},
							},
						},
						{
							Name:            "syncthing",
							Image:           ImageName,
							ImagePullPolicy: "Always",
							Command: []string{
								"/syncthing/syncthing",
								"-home", "/var/syncthing/config",
								"-gui-apikey", viper.GetString("apikey"),
								"-verbose",
							},
							Ports: []v1.ContainerPort{
								{ContainerPort: s.SyncthingAPI, Name: "rest"},
								{ContainerPort: s.SyncthingListener, Name: "sync"},
							},
							// TODO: resources
							VolumeMounts: []v1.VolumeMount{
								v1.VolumeMount{
									Name:      "dockerfs",
									MountPath: viper.GetString("docker-root"),
								},
								v1.VolumeMount{
									Name:      "dockersock",
									MountPath: viper.GetString("docker-socket"),
								},
								v1.VolumeMount{
									Name:      "kubelet",
									MountPath: "/var/lib/kubelet",
								},
							},
							LivenessProbe: &v1.Probe{
								Handler: v1.Handler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.FromInt(int(s.SyncthingAPI)),
									},
								},
								InitialDelaySeconds: 10,
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.FromInt(int(s.SyncthingAPI)),
									},
								},
								InitialDelaySeconds: 10,
							},
						},
					},
					NodeSelector: map[string]string{
						"beta.kubernetes.io/os": "linux",
					},
					ServiceAccountName: s.name,
					// TODO: add HostPathType
					Volumes: []v1.Volume{
						v1.Volume{
							Name: "dockerfs",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: viper.GetString("docker-root"),
								},
							},
						},
						v1.Volume{
							Name: "dockersock",
							VolumeSource: v1.VolumeSource{
								HostPath: &v1.HostPathVolumeSource{
									Path: viper.GetString("docker-socket"),
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
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: "RollingUpdate",
			},
		},
	}

	collection := Client.AppsV1().DaemonSets(s.Namespace)

	if _, err := collection.Create(daemonSet); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := collection.Update(daemonSet); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) createServiceAccount(upgrade bool) error {
	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.Namespace,
			Name:      s.name,
			Labels:    s.labels,
		},
	}

	collection := Client.CoreV1().ServiceAccounts(s.Namespace)

	if _, err := collection.Create(serviceAccount); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := collection.Update(serviceAccount); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) createPSP(upgrade bool) error {
	psp := &policyv1beta.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.name,
			Labels: s.labels,
		},
		Spec: policyv1beta.PodSecurityPolicySpec{
			FSGroup: policyv1beta.FSGroupStrategyOptions{
				Rule: policyv1beta.FSGroupStrategyRunAsAny,
			},
			RunAsUser: policyv1beta.RunAsUserStrategyOptions{
				Rule: policyv1beta.RunAsUserStrategyRunAsAny,
			},
			SELinux: policyv1beta.SELinuxStrategyOptions{
				Rule: policyv1beta.SELinuxStrategyRunAsAny,
			},
			SupplementalGroups: policyv1beta.SupplementalGroupsStrategyOptions{
				Rule: policyv1beta.SupplementalGroupsStrategyRunAsAny,
			},
			Volumes: []v1beta1.FSType{
				v1beta1.HostPath,
				v1beta1.Secret,
			},
		},
	}

	collection := Client.PolicyV1beta1().PodSecurityPolicies()

	if _, err := collection.Create(psp); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := collection.Update(psp); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) createClusterRole(upgrade bool) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.name,
			Labels: s.labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:         []string{"use"},
				APIGroups:     []string{"policy"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{s.name},
			},
		},
	}

	collection := Client.RbacV1().ClusterRoles()

	if _, err := collection.Create(clusterRole); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := collection.Update(clusterRole); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) createClusterRoleBinding(upgrade bool) error {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.name,
			Labels: s.labels,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     s.name,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      s.name,
			Namespace: s.Namespace,
		}},
	}

	collection := Client.RbacV1().ClusterRoleBindings()

	if _, err := collection.Create(clusterRoleBinding); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	if upgrade {
		if _, err := collection.Update(clusterRoleBinding); err != nil {
			return err
		}
	}
	return nil
}
