package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitLogging(t *testing.T) {
	// TODO: There must be something more we can test here.
	require.NotPanics(t, func() { InitLogging() })
}
