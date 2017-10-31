package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

var findercmd = &FinderCmd{
	BaseCmd: BaseCmd{
		Root:  "testing",
		Cmd:   &cobra.Command{},
		Viper: viper.New(),
	},
}

// We have to import the ksync pkg and initialize the k8s client otherwise the deathstar goes boom
func init() {
	if err := ksync.InitKubeClient("", "default"); err != nil {
		log.Fatal(err)
	}
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

func TestContainers(t *testing.T) {
	containerList, err := findercmd.Containers()

	assert.NoError(t, err)
	assert.NotNil(t, containerList)
}
