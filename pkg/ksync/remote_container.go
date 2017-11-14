package ksync

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/vapor-ware/ksync/pkg/debug"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// RemoteContainer is a specific container running on the remote cluster.
type RemoteContainer struct {
	ID       string
	Name     string
	NodeName string
	PodName  string
}

func (c *RemoteContainer) String() string {
	return debug.YamlString(c)
}

// Fields returns a set of structured fields for logging.
func (c *RemoteContainer) Fields() log.Fields {
	return debug.StructFields(c)
}

// Radar connects to the server component (radar) and returns a client.
func (c *RemoteContainer) Radar() (pb.RadarClient, error) {
	conn, err := NewRadarInstance().RadarConnection(c.NodeName)
	if err != nil {
		return nil, err
	}
	// TODO: what's a better way to handle c?
	// defer conn.Close()

	log.WithFields(c.Fields()).Debug("radar connected")

	return pb.NewRadarClient(conn), nil
}

// RestartMirror restarts the remote mirror container responsible for this
// container.
func (c *RemoteContainer) RestartMirror() error {
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

func getRemoteContainer(
	pod *apiv1.Pod, containerName string) (*RemoteContainer, error) {
	// TODO: I don't want to go and setup a whole bunch of error handling code
	// right now. The NotFound error works perfectly here for now.
	if pod.DeletionTimestamp != nil {
		return nil, &errors.StatusError{
			ErrStatus: metav1.Status{
				Status:  metav1.StatusFailure,
				Reason:  metav1.StatusReasonNotFound,
				Message: fmt.Sprintf("%s scheduled for deletion", pod.Name),
			}}
	}

	// TODO: runtime error because there are no container statuses while
	// k8s master is restarting.
	// TODO: added, but non-running containers don't have any status. This should
	// be converted to something closer to `IsNotRunning()`
	// and `IsMissingContainer()`
	if containerName == "" && len(pod.Status.ContainerStatuses) > 0 {
		return &RemoteContainer{
			pod.Status.ContainerStatuses[0].ContainerID[9:],
			pod.Status.ContainerStatuses[0].Name,
			pod.Spec.NodeName,
			pod.Name}, nil
	}

	for _, status := range pod.Status.ContainerStatuses {
		if status.Name != containerName {
			continue
		}

		return &RemoteContainer{
			status.ContainerID[9:],
			status.Name,
			pod.Spec.NodeName,
			pod.Name,
		}, nil
	}

	return nil, &errors.StatusError{
		ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Reason:  metav1.StatusReasonNotFound,
			Message: fmt.Sprintf("%s not running %s", pod.Name, containerName),
		}}
}

// GetByName takes a pod and container name and looks for a running RemoteContainer.
func GetByName(podName string, containerName string) (*RemoteContainer, error) {
	pod, err := kubeClient.CoreV1().Pods(namespace).Get(
		podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"name": pod.Name,
	}).Debug("found pod")

	return getRemoteContainer(pod, containerName)
}

func getBySelector(selector string, containerName string) ([]*RemoteContainer, error) {
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

	containerList := []*RemoteContainer{}
	for _, pod := range pods.Items {
		cntr, err := getRemoteContainer(&pod, containerName)
		if errors.IsNotFound(err) {
			continue
		}

		if err != nil {
			return nil, err
		}

		containerList = append(containerList, cntr)
	}

	return containerList, nil
}

// GetRemoteContainers uses a Locator (podName, selector, containerName) and provides
// a list of the currently running containers. It is possible that this list is
// empty, for when the locator does not match anything currently running.
// TODO: this takes a little bit to execute, is there any kind of progress or output
// that would be useful to the user?
// TODO: make this into a channel
func GetRemoteContainers(
	podName string,
	selector string,
	containerName string) ([]*RemoteContainer, error) {

	containerList := []*RemoteContainer{}

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
