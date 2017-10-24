package cli

import (
	"testing"
	// "reflect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "github.com/spf13/pflag"
)

// TODO: This is what's left of the unit testing. Need to write some integration tests here, but leaving this for inspiration when we eventually write proper unit tests.

func TestInit(t *testing.T) {
	base := &BaseCmd{}
	cmd := &cobra.Command{}

	// Test for panics during setup
	require.NotPanics(t, func() { base.Init("rooty", cmd) })

	// Test to see if values are correctly set
	assert.Contains(t, base.Root, "rooty")
}

func TestCmdBindFlag(t *testing.T) {
	base := &BaseCmd{
		Root:  "rooty",
		Cmd:   &cobra.Command{},
		Viper: viper.New(),
	}
	base.Cmd.Flags().String(
		"someflag",
		"",
		"Some flag that I don't like")
	base.Cmd.ParseFlags([]string{"someflag"})

	// Test for run errors
	err := base.BindFlag("someflag")
	assert.NoError(t, err)
}
