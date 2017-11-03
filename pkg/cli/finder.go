package cli

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

// FinderCmd parses the options required to discover a remote container.
type FinderCmd struct {
	BaseCmd
}

// DefaultFlags configures the default flags to find a container.
func (cmd *FinderCmd) DefaultFlags() error {
	flags := cmd.Cmd.Flags()
	flags.StringP(
		"container",
		"c",
		"",
		"Container name. Defaults to first container.")
	if err := cmd.BindFlag("container"); err != nil {
		return err
	}

	// TODO: is this best as an arg instead of positional?
	flags.StringP(
		"pod",
		"p",
		"",
		"Pod name.")
	if err := cmd.BindFlag("pod"); err != nil {
		return err
	}

	flags.StringP(
		"selector",
		"l",
		"",
		"Selector (label query) to filter on, supports '=', '==', and '!='.")
	return cmd.BindFlag("selector")
}

// Validator ensures that the command has valid arguments.
func (cmd *FinderCmd) Validator() error {
	// TODO: something like cmdutil.UsageErrorf
	// TODO: move into its own function (add to command as a validator?)
	if cmd.Viper.GetString("selector") == "" && cmd.Viper.GetString("pod") == "" {
		return fmt.Errorf("must specify at least a selector or a pod name")
	}

	return nil
}

// RemoteContainers returns a list of all remote containers that match the finder args.
func (cmd *FinderCmd) RemoteContainers() ([]*ksync.RemoteContainer, error) {
	containerList, err := ksync.GetRemoteContainers(
		cmd.Viper.GetString("pod"),
		cmd.Viper.GetString("selector"),
		cmd.Viper.GetString("container"))
	if err != nil {
		return nil, errors.Wrap(err, "could not get containers")
	}

	// TODO: maybe there's a better way to do this?
	if len(containerList) == 0 {
		return nil, fmt.Errorf("no containers found")
	}

	return containerList, nil
}
