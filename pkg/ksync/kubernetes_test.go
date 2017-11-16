package ksync

import (
	"testing"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitKubeClient(t *testing.T) {
	err := InitKubeClient("")

	// TODO: There has to be a better set of tests here
	require.NoError(t, err)
}
