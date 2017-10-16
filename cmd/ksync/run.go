package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/input"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type runCmd struct {
	viper *viper.Viper
}

func (this *runCmd) new() *cobra.Command {
	long := `
    Start syncing between a local and remote directory.`
	example := ``

	cmd := &cobra.Command{
		Use:     "run [flags] [local path] [remote path]",
		Short:   "Start syncing between a local and remote directory.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(2),
		Run:     this.run,
		// TODO: BashCompletionFunction
	}
	this.viper = viper.New()

	flags := cmd.Flags()
	flags.StringP(
		"container",
		"c",
		"",
		"Container name. If omitted, the first container in the pod will be chosen.")

	this.viper.BindPFlag("container", flags.Lookup("container"))
	this.viper.BindEnv("container", "KSYNC_CONTAINER")

	// TODO: is this best as an arg instead of positional?
	flags.StringP(
		"pod",
		"p",
		"",
		"Pod name.")

	this.viper.BindPFlag("pod", flags.Lookup("pod"))
	this.viper.BindEnv("pod", "KSYNC_POD")

	return cmd
}

// TODO: check for existence of java (and the right version)
// TODO: message (and fail) when this is not run from the expected environment -
//       the docker container.
func (this *runCmd) run(cmd *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	if this.viper.GetString("pod") == "" {
		log.Fatal("Must specify --pod.")
	}

	syncPath := input.GetSyncPath(args)
	syncPath.Validator()

	container, err := ksync.GetByName(
		this.viper.GetString("pod"),
		this.viper.GetString("container"))
	if err != nil {
		log.Fatalf(
			"Could not get pod(%s) container(%s): %v",
			this.viper.GetString("pod"),
			this.viper.GetString("container"),
			err)
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
