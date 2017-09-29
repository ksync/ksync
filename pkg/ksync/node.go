package ksync

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// TODO: make this configurable?
	label = "radar-name"
)

// TODO: run in parallel
// TODO: this isn't really used yet, it will be employed for dynamic startup
// of pods instead of using the DaemonSet.
func PrepareNodes() error {
	nodes, err := KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		if _, ok := node.Labels[label]; ok {
			continue
		}

		node.Labels[label] = node.Name
		if _, err := KubeClient.CoreV1().Nodes().Update(&node); err != nil {
			log.Print(err)
		}
	}

	return nil
}
