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

func (r *runCmd) new() *cobra.Command {
	long := `
    Start syncing between a local and remote directory.`
	example := ``

	cmd := &cobra.Command{
		Use:     "run [flags] [local path] [remote path]",
		Short:   "Start syncing between a local and remote directory.",
		Long:    long,
		Example: example,
		Args:    cobra.ExactArgs(2),
		Run:     r.run,
		// TODO: BashCompletionFunction
	}
	r.viper = viper.New()

	flags := cmd.Flags()
	flags.StringP(
		"container",
		"c",
		"",
		"Container name. If omitted, the first container in the pod will be chosen.")

	r.viper.BindPFlag("container", flags.Lookup("container"))
	r.viper.BindEnv("container", "KSYNC_CONTAINER")

	// TODO: is this best as an arg instead of positional?
	flags.StringP(
		"pod",
		"p",
		"",
		"Pod name.")

	r.viper.BindPFlag("pod", flags.Lookup("pod"))
	r.viper.BindEnv("pod", "KSYNC_POD")

	return cmd
}

// TODO: check for existence of java (and the right version)
// TODO: message (and fail) when this is not run from the expected environment -
//       the docker container.
func (r *runCmd) run(cmd *cobra.Command, args []string) {
	// Usage validation ------------------------------------
	if r.viper.GetString("pod") == "" {
		log.Fatal("Must specify --pod.")
	}

	syncPath := input.GetSyncPath(args)
	syncPath.Validator()

	container, err := ksync.GetByName(
		r.viper.GetString("pod"),
		r.viper.GetString("container"))
	if err != nil {
		log.Fatalf(
			"Could not get pod(%s) container(%s): %v",
			r.viper.GetString("pod"),
			r.viper.GetString("container"),
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
