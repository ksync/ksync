package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

var (
	// TODO: update the usage instructions
	listHelp = `
    List the files from a remote container.
    `

	// TODO: add path in here. Should it be something like the ssh form? (see oc rsync)
	listCmd = &cobra.Command{
		Use:     "list [flags] [path]",
		Short:   "List files from a remote container.",
		Long:    listHelp,
		Aliases: []string{"ls"},
		Run:     runList,
		// TODO: BashCompletionFunction
	}
)

func getContainers(
	selector string, container string) ([]*pb.ContainerPath, error) {
	return nil, nil
}

func runList(_ *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	// TODO: something like cmdutil.UsageErrorf
	// TODO: move into its own function (add to command as a validator?)
	if viper.GetString("selector") == "" && viper.GetString("pod") == "" {
		log.Fatal("Must specify at least a selector or a pod name.")
	}

	if len(args) == 0 {
		log.Fatal("Must specify a container path.")
	}

	// REMOVE ME
	ksync.PrepareNodes()

	path := args[0]

	containerList, err := ksync.GetContainers(
		viper.GetString("pod"),
		viper.GetString("selector"),
		viper.GetString("container"))
	if err != nil {
		log.Fatalf(
			"could not get containers for (pod:%s) (selector:%s) (container:%s): %v",
			viper.GetString("pod"),
			viper.GetString("selector"),
			viper.GetString("container"),
			err)
	}

	// TODO: maybe there's a better way to do this?
	if len(containerList) == 0 {
		log.Fatalf(
			"no containers found for pod (%s) or selector (%s) with container (%s)",
			viper.GetString("pod"),
			viper.GetString("selector"),
			viper.GetString("container"))
	}

	// TODO: make this into a channel?
	for _, cntr := range containerList {
		conn, err := ksync.NewRadarConnection(cntr.NodeName)
		if err != nil {
			log.Fatalf("Could not connect to radar: %v", err)
		}
		defer conn.Close()

		log.WithFields(log.Fields{
			"node": cntr.NodeName,
		}).Debug("radar connected")

		files, err := pb.NewRadarClient(conn).ListContainerFiles(
			context.Background(), &pb.ContainerPath{cntr.ID, path})
		if err != nil {
			log.Fatalf("Failed getting files: %v", err)
		}

		// TODO: improve output
		log.Debugf("%s", files)
	}
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP(
		"container",
		"c",
		"",
		"Container name. If omitted, the first container in the pod will be chosen.")

	viper.BindPFlag("container", listCmd.Flags().Lookup("container"))
	viper.BindEnv("container", "KSYNC_CONTAINER")

	// TODO: is this best as an arg instead of positional?
	listCmd.Flags().StringP(
		"pod",
		"p",
		"",
		"Pod name.")

	viper.BindPFlag("pod", listCmd.Flags().Lookup("pod"))
	viper.BindEnv("pod", "KSYNC_POD")

	listCmd.Flags().StringP(
		"selector",
		"l",
		"",
		"Selector (label query) to filter on, supports '=', '==', and '!='.")

	viper.BindPFlag("selector", listCmd.Flags().Lookup("selector"))
	viper.BindEnv("selector", "KSYNC_SELECTOR")
}
