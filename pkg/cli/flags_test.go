package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
// "github.com/stretchr/testify/require"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"

)

func TestBindFlag(t *testing.T) {
	viper := viper.New()
  // Create a fake flag to replace
	flag := pflag.Flag{
		Name: "testflag",
	}
	err := BindFlag(viper, &flag, "ksync")

	assert.NoError(t, err)
}

func TestDefaultFlags(t *testing.T) {
	cmd := &cobra.Command{}

	err := DefaultFlags(cmd, "ksync")

	assert.NoError(t, err)
	assert.NotNil(t, cmd.PersistentFlags())
}
