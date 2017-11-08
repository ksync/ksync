package ksync

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	InitKubeClient("", os.Getenv("TEST_RADAR_NAMESPACE"))
}

func TestRadar(t *testing.T) {
	con := &RemoteContainer{
		ID:   os.Getenv("TEST_RADAR_CONTAINERID"),
		Name: "",
		// TODO: This has to be dynamic
		NodeName: os.Getenv("TEST_RADAR_NODE"),
		PodName:  os.Getenv("TEST_RADAR_POD"),
	}
	_, err := con.Radar()
	// TODO: Remove logging
	t.Log(err)

	// TODO: Is there anything else we can test? There has to be a better way of doing this.
	require.NoError(t, err)
}

func TestGetByName(t *testing.T) {
	// TODO: This has to be dynamic
	podName := os.Getenv("TEST_POD")
	containerName := "someequallystupidname"

	// Test erroring on empty containerName
	_, err := GetByName(podName, "")
	assert.NoError(t, err)

	// Test default error
	_, err = GetByName(podName, containerName)
	require.Error(t, err)
}
