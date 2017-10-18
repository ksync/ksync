package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BaseCmd struct {
	Root  string
	Cmd   *cobra.Command
	Viper *viper.Viper
}

func (b *BaseCmd) Init(root string, cmd *cobra.Command) {
	b.Root = root
	b.Cmd = cmd
	b.Viper = viper.New()
}

func (b *BaseCmd) BindFlag(name string) error {
	if err := b.Viper.BindPFlag(name, b.Cmd.Flags().Lookup(name)); err != nil {
		return err
	}

	if err := b.Viper.BindEnv(
		name, strings.ToUpper(fmt.Sprintf("%s_%s", b.Root, name))); err != nil {
		return err
	}

	return nil
}
