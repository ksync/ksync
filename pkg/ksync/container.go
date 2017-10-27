package ksync

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// Container is a specific container running on the remote cluster.
type Container struct {
	ID       string
	Name     string
	NodeName string
	PodName  string
}

func (c *Container) String() string {
	return YamlString(c)
}

// Fields returns a set of structured fields for logging.
func (c *Container) Fields() log.Fields {
	return StructFields(c)
}

// Radar connects to the server component (radar) and returns a client.
func (c *Container) Radar() (pb.RadarClient, error) {
	conn, err := NewRadarInstance().RadarConnection(c.NodeName)
	if err != nil {
		return nil, ErrorOut("Could not connect to radar", err, c)
	}
	// TODO: what's a better way to handle c?
	// defer conn.Close()

	log.WithFields(c.Fields()).Debug("radar connected")

	return pb.NewRadarClient(conn), nil
}

// RestartMirror restarts the remote mirror container responsible for this
// container.
func (c *Container) RestartMirror() error {
	client, err := c.Radar()
	if err != nil {
		return err
	}

	if _, err := client.RestartMirror(
		context.Background(), &empty.Empty{}); err != nil {
		return err
	}

	return nil
}

func getContainer(pod *apiv1.Pod, containerName string) (*Container, error) {
	// TODO: runtime error because there are no container statuses while
	// k8s master is restarting.
	if containerName == "" {
		if len(pod.Status.ContainerStatuses) == 0 {
			return nil, fmt.Errorf("no status for container")
		}

		return &Container{
			pod.Status.ContainerStatuses[0].ContainerID[9:],
			pod.Status.ContainerStatuses[0].Name,
			pod.Spec.NodeName,
			pod.Name}, nil
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.Name != containerName {
			continue
		}

		return &Container{
			status.ContainerID[9:],
			status.Name,
			pod.Spec.NodeName,
			pod.Name,
		}, nil
	}

	// TODO: should this work like `GetContainers` does and just return nil?
	return nil, fmt.Errorf(
		"could not find container (%s) in pod (%s)",
		containerName,
		pod.Name)
}

// GetByName takes a pod and container name and looks for a running Container.
func GetByName(podName string, containerName string) (*Container, error) {
	pod, err := kubeClient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
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
	// TODO: namespace is not global anywhere else.
	pods, err := kubeClient.CoreV1().Pods(namespace).List(opts)
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

// GetContainers uses a Locator (podName, selector, containerName) and provides
// a list of the currently running containers. It is possible that this list is
// empty, for when the locator does not match anything currently running.
// TODO: this takes a little bit to execute, is there any kind of progress or output
// that would be useful to the user?
// TODO: make this into a channel
func GetContainers(
	podName string,
	selector string,
	containerName string) ([]*Container, error) {

	containerList := []*Container{}

	if podName != "" {
		container, err := GetByName(podName, containerName)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}

		// It is possible that a container wasn't found. We don't want to error out,
		// instead just return an empty list.
		if container != nil {
			containerList = append(containerList, container)
		}
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
