package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type runCmd struct {
	cli.BaseCmd
}

// TODO: update to the container method.
func (r *runCmd) new() *cobra.Command {
	long := `
    Start syncing between a local and remote directory.

    Note: this is meant to be run from within the ksync container.`
	example := ``

	r.Init("ksync", &cobra.Command{
		Use:     "run [flags] [local path] [remote path]",
		Short:   "Start syncing between a local and remote directory.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(2),
		Run:     r.run,
		// TODO: BashCompletionFunction
	})

	r.Cmd.Flags().StringP(
		"container",
		"c",
		"",
		"Container name. Defaults to the first container in pod.")
	if err := r.BindFlag("container"); err != nil {
		log.Fatal(err)
	}

	// TODO: is this best as an arg instead of positional?
	r.Cmd.Flags().StringP(
		"pod",
		"p",
		"",
		"Pod name.")
	if err := r.BindFlag("pod"); err != nil {
		log.Fatal(err)
	}

	return r.Cmd
}

// TODO: check for existence of java (and the right version)
// TODO: message (and fail) when this is not run from the expected environment -
//       the docker container.
func (r *runCmd) run(cmd *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	if r.Viper.GetString("pod") == "" {
		log.Fatal("Must specify --pod.")
	}

	syncPath := input.GetSyncPath(args)
	if err := syncPath.Validator(); err != nil {
		log.Fatal(err)
	}

	container, err := ksync.GetByName(
		r.Viper.GetString("pod"),
		r.Viper.GetString("container"))
	if err != nil {
		log.Fatalf(
			"Could not get pod(%s) container(%s): %v",
			r.Viper.GetString("pod"),
			r.Viper.GetString("container"),
			err)
	}

	// TODO: can/should we be a little bit more intelligent here?
	if err := container.RestartMirror(); err != nil {
		log.Fatal(err)
	}

	mirror := &ksync.Mirror{
		Container:  container,
		LocalPath:  syncPath.Local,
		RemotePath: syncPath.Remote,
	}
	if err := mirror.Run(); err != nil {
		log.Fatal(err)
	}
}
