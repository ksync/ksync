package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var findercmd = &FinderCmd{
	BaseCmd: BaseCmd{
		Root:  "testing",
		Cmd:   &cobra.Command{},
		Viper: viper.New(),
	},
}

func TestDefaultFlags(t *testing.T) {
	t.Log(*findercmd)
	err := findercmd.DefaultFlags()

	assert.NoError(t, err)
}

func TestValidator(t *testing.T) {
	err := findercmd.Validator()
	findercmd.BaseCmd.Viper.Set("selector", "app=testapp")

	assert.NoError(t, err)
}

func TestContainers(t *testing.T) {
	containerList, err := findercmd.Containers()

	assert.NoError(t, err)
	assert.NotNil(t, containerList)
}
