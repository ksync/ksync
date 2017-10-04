package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// TODO: update the usage instructions
	runHelp = `
    Start syncing between a local and remote directory.
    `

	runCmd = &cobra.Command{
		Use:   "run [flags] [local path] [remote path]",
		Short: "Start syncing between a local and remote directory.",
		Long:  runHelp,
		Args:  cobra.ExactArgs(2),
		Run:   runStart,
		// TODO: BashCompletionFunction
	}

	runViper = viper.New()
)

// TODO: check for existence of java (and the right version)
// TODO: download the jar locally (into a ksync home directory?)
// TODO: move checks/downloads into init?
func runStart(_ *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	if runViper.GetString("pod") == "" {
		log.Fatal("Must specify --pod.")
	}

	paths := input.GetPaths(args)
	paths.Validator()

	container, err := ksync.GetByName(
		runViper.GetString("pod"),
		runViper.GetString("container"))
	if err != nil {
		log.Fatalf(
			"Could not get pod(%s) container(%s): %v",
			runViper.GetString("pod"),
			runViper.GetString("container"),
			err)
	}

	mirror := &ksync.Mirror{
		Container:  container,
		LocalPath:  paths.Local,
		RemotePath: paths.Remote,
	}
	if err := mirror.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP(
		"container",
		"c",
		"",
		"Container name. If omitted, the first container in the pod will be chosen.")

	runViper.BindPFlag("container", runCmd.Flags().Lookup("container"))
	runViper.BindEnv("container", "KSYNC_CONTAINER")

	// TODO: is this best as an arg instead of positional?
	runCmd.Flags().StringP(
		"pod",
		"p",
		"",
		"Pod name.")

	runViper.BindPFlag("pod", runCmd.Flags().Lookup("pod"))
	runViper.BindEnv("pod", "KSYNC_POD")
}
