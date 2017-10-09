package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EnvName constructs the name of the environment from the current command
// optionally provided name
func envName(cmd string, name string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(cmd), strings.ToUpper(name))
}

// DefaultFlags provides a definition for the command flags
// TODO: add doc strings
func DefaultFlags(cmd *cobra.Command, name string) {
	cmd.PersistentFlags().String(
		"config",
		"",
		fmt.Sprintf("config file (default is $HOME/.%s.yaml", name))

	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
	viper.BindEnv("config", envName(name, "config"))

	// TODO: can this be limited to a selection?
	cmd.PersistentFlags().String(
		"log-level",
		"warn",
		"log level to use.")

	viper.BindPFlag("log-level", cmd.PersistentFlags().Lookup("log-level"))
	viper.BindEnv("log-level", envName(name, "log_level"))
}
