package ksync

import (
	"os"
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	InitKubeClient("") // nolint: errcheck
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
