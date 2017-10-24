package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	configname := "ksync"
	err := InitConfig(configname)

	require.NoError(t, err)
}
