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

func init() {
	findercmd.BaseCmd.Viper.Set("selector", "app=auth")
}

func TestFinderDefaultFlags(t *testing.T) {
	err := findercmd.DefaultFlags()

	assert.NoError(t, err)
}

func TestValidator(t *testing.T) {
	err := findercmd.Validator()

	assert.NoError(t, err)
}
