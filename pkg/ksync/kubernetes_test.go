package ksync

import (
	"testing"
	"os"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitKubeClient(t *testing.T) {
	err := InitKubeClient("", os.Getenv("TEST_NAMESPACE"))

	// TODO: There has to be a better set of tests here
	require.NoError(t, err)
}
