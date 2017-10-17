package input

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

// Locator contains everything required to discover a remote container.
type Locator struct {
	PodName       string
	Selector      string
	ContainerName string
}

// GetLocator constructs a Locator from user configuration.
func GetLocator(cmdViper *viper.Viper) Locator {
	return Locator{
		cmdViper.GetString("pod"),
		cmdViper.GetString("selector"),
		cmdViper.GetString("container"),
	}
}

// Validator ensures that a Locator is valid for use.
func (l *Locator) Validator() {
	// TODO: something like cmdutil.UsageErrorf
	// TODO: move into its own function (add to command as a validator?)
	if l.Selector == "" && l.PodName == "" {
		log.Fatal("Must specify at least a selector or a pod name.")
	}
}

// Containers returns a list of all remote containers that match the Locator.
func (l *Locator) Containers() ([]*ksync.Container, error) {
	containerList, err := ksync.GetContainers(
		l.PodName, l.Selector, l.ContainerName)
	if err != nil {
		return nil, fmt.Errorf(
			"could not get containers for (pod:%s) (selector:%s) (container:%s): %v",
			l.PodName,
			l.Selector,
			l.ContainerName,
			err)
	}

	// TODO: maybe there's a better way to do this?
	if len(containerList) == 0 {
		return nil, fmt.Errorf(
			"no containers found for pod (%s) or selector (%s) with container (%s)",
			l.PodName,
			l.Selector,
			l.ContainerName)
	}

	return containerList, nil
}

// LocatorFlags sets the config options required to construct a locator.
func LocatorFlags(cmd *cobra.Command, cmdViper *viper.Viper) {
	cmd.Flags().StringP(
		"container",
		"c",
		"",
		"Container name. If omitted, the first container in the pod will be chosen.")

	cmdViper.BindPFlag("container", cmd.Flags().Lookup("container"))
	cmdViper.BindEnv("container", "KSYNC_CONTAINER")

	// TODO: is this best as an arg instead of positional?
	cmd.Flags().StringP(
		"pod",
		"p",
		"",
		"Pod name.")

	cmdViper.BindPFlag("pod", cmd.Flags().Lookup("pod"))
	cmdViper.BindEnv("pod", "KSYNC_POD")

	cmd.Flags().StringP(
		"selector",
		"l",
		"",
		"Selector (label query) to filter on, supports '=', '==', and '!='.")

	cmdViper.BindPFlag("selector", cmd.Flags().Lookup("selector"))
	cmdViper.BindEnv("selector", "KSYNC_SELECTOR")
}
