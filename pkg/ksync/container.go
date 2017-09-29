package ksync

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

// TODO: is NodeName always there on the spec post-start?
type Container struct {
	ID       string
	NodeName string
}

func (this *Container) String() string {
	return fmt.Sprintf("ID: %s, NodeName: %s", this.ID, this.NodeName)
}

func getContainer(pod *apiv1.Pod, containerName string) (*Container, error) {
	// TODO: runtime error because there are no container statuses while
	// k8s master is restarting.
	log.Print(pod.Status.ContainerStatuses)
	if containerName == "" {
		return &Container{
			pod.Status.ContainerStatuses[0].ContainerID[9:],
			pod.Spec.NodeName}, nil
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.Name != containerName {
			continue
		}

		return &Container{
			status.ContainerID[9:],
			pod.Spec.NodeName,
		}, nil
	}

	return nil, fmt.Errorf(
		"could not find container (%s) in pod (%s)",
		containerName,
		pod.Name)
}

func getByName(podName string, containerName string) (*Container, error) {
	pod, err := KubeClient.CoreV1().Pods(Namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"name": pod.Name,
	}).Debug("found pod")

	return getContainer(pod, containerName)
}

func getBySelector(selector string, containerName string) ([]*Container, error) {
	opts := metav1.ListOptions{}
	opts.LabelSelector = selector
	pods, err := KubeClient.CoreV1().Pods(Namespace).List(opts)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"length":   len(pods.Items),
		"selector": selector,
	}).Debug("found pods by selector")

	containerList := []*Container{}
	for _, pod := range pods.Items {
		cntr, err := getContainer(&pod, containerName)
		if err != nil {
			return nil, err
		}
		containerList = append(containerList, cntr)
	}

	return containerList, nil
}

// TODO: this takes a little bit to execute, is there any kind of progress or output
// that would be useful to the user?
// TODO: make this into a channel
func GetContainers(
	podName string,
	selector string,
	containerName string) ([]*Container, error) {

	containerList := []*Container{}

	if podName != "" {
		container, err := getByName(podName, containerName)
		if err != nil {
			return nil, err
		}

		containerList = append(containerList, container)
	}

	if selector != "" {
		selectorList, err := getBySelector(selector, containerName)
		if err != nil {
			return nil, err
		}
		containerList = append(containerList, selectorList...)
	}

	log.WithFields(log.Fields{
		"podName":       podName,
		"selector":      selector,
		"containerName": containerName,
		"length":        len(containerList),
	}).Debug("found containers")

	return containerList, nil
}
