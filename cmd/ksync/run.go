package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
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

	localPath := args[0]
	remotePath := args[1]

	if !filepath.IsAbs(localPath) {
		log.Fatal("Local path must be absolute.")
	}

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		log.Fatal("Local path must exist.")
	}

	if !filepath.IsAbs(remotePath) {
		log.Fatal("Local path must be absolute.")
	}

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

	client, err := container.Radar()
	if err != nil {
		log.Fatal(err)
	}

	path, err := client.GetAbsPath(
		context.Background(), &pb.ContainerPath{container.ID, remotePath})
	if err != nil {
		log.Fatalf("Could not get root path: %v", err)
	}

	tun, err := ksync.NewTunnel(container.NodeName, 49172)
	if err != nil {
		log.Fatal(err)
	}

	if err := tun.Start(); err != nil {
		log.Fatalf("Error starting tunnel: %v", err)
	}

	cmdArgs := []string{
		"-Xmx2G",
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-cp", "/home/thomas/work/bin/mirror-all.jar",
		"mirror.Mirror",
		"client",
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", tun.LocalPort),
		"-l", localPath,
		"-r", path.Full,
	}

	cmd := exec.Command("java", cmdArgs...)

	// cmd := exec.Command("/bin/bash", "-c", "while true; do sleep 1; echo 'adsf'; done")

	if err := ksync.LogCmdStream(cmd); err != nil {
		log.Fatal(err)
	}

	// Setup the k8s runtime to fail on unreturnable error (instead of looping).
	// This helps cleanup zombie java processes.
	runtime.ErrorHandlers = append(runtime.ErrorHandlers, func(err error) {
		cmd.Process.Kill()
		// TODO: this makes me feel dirty, there must be a better way.
		if strings.Contains(err.Error(), "Connection refused") {
			log.Fatal(
				"Lost connection to remote radar pod. Please try again (it should restart).")
		}

		log.Fatalf("unreturnable error: %v", err)
	})

	if err := cmd.Start(); err != nil {
		log.Fatalf("%v", err)
	}

	log.WithFields(log.Fields{
		"cmd":  cmd.Path,
		"args": cmd.Args,
	}).Debug("starting mirror")

	cmd.Wait()
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
