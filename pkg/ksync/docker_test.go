package ksync

import (
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

func TestInitDockerClient(t *testing.T) {
	err := InitDockerClient()

	require.NoError(t, err)
}
