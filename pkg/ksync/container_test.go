package ksync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

func TestRadar(t *testing.T) {
	con := &Container{
		NodeName: "nope",
	}
	_, err := con.Radar()

  // TODO: Is there anything else we can test? There has to be a better way of doing this.
	require.EqualError(t, err, "Could not connect to radar")
}

func TestGetByName(t *testing.T) {
	podName := "somestupidname"
	containerName := "someequallystupidname"

  // Test erroring on empty containerName
	_, err := GetByName(podName, "")
	assert.EqualError(t, err, "no status for container")

  // Test default error
	_, err = GetByName(podName, containerName)
	require.Error(t, err)
}
