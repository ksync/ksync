package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

func init() {
	InitKubeClient("", "kube-system")
}

func TestRadar(t *testing.T) {
	con := &Container{
		ID: "",
		Name: "",
    // TODO: This has to be dynamic
		NodeName: "gke-tim-dev-default-pool-8e1c25a0-j3kp",
		PodName: "",
	}
	_, err := con.Radar()
  // TODO: Remove logging
	t.Log(err)

  // TODO: Is there anything else we can test? There has to be a better way of doing this.
	require.NoError(t, err)
}

func TestGetByName(t *testing.T) {
	// TODO: This has to be dynamic
	podName := "ksync-radar-01q6r"
	containerName := "someequallystupidname"

  // Test erroring on empty containerName
	_, err := GetByName(podName, "")
	assert.NoError(t, err)

  // Test default error
	_, err = GetByName(podName, containerName)
	require.Error(t, err)
}
