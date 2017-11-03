package docker

import (
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitClient(t *testing.T) {
	err := InitClient()

	require.NoError(t, err)
}
