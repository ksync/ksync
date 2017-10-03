package main

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

var (
	// TODO: update the usage instructions
	startHelp = `
    Start syncing between a local and remote directory.
    `

	startCmd = &cobra.Command{
		Use:   "start [flags] [local path] [remote path]",
		Short: "Start syncing between a local and remote directory.",
		Long:  startHelp,
		Args:  cobra.ExactArgs(2),
		Run:   runStart,
		// TODO: BashCompletionFunction
	}

	startViper = viper.New()
)

func runStart(_ *cobra.Command, args []string) {
	loc := ksync.GetLocator(startViper)
	// Usage validation ------------------------------------
	loc.Validator()

	localPath := args[0]
	remotePath := args[1]

	containerList, err := loc.Containers()
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, cntr := range containerList {
		client, err := cntr.Radar()
		if err != nil {
			log.Fatalf("%v", err)
		}

		path, err := client.GetAbsPath(
			context.Background(), &pb.ContainerPath{cntr.ID, remotePath})
		if err != nil {
			log.Fatalf("Could not get root path: %v", err)
		}

		tun, err := ksync.NewTunnel(cntr.NodeName, 49172)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if err := tun.Start(); err != nil {
			log.Fatalf("Error starting tunnel: %v", err)
		}

		args := []string{
			"client",
			"-h", "localhost",
			"-p", fmt.Sprintf("%d", tun.LocalPort),
			"-l", localPath,
			"-r", path.Full,
		}

		cmd := exec.Command("mirror", args...)
		if err := cmd.Start(); err != nil {
			log.Fatalf("%v", err)
		}

		log.WithFields(log.Fields{
			"cmd":  "mirror",
			"args": args,
		}).Debug("starting mirror")

		cmd.Wait()

		log.Print(path)
	}
}

func init() {
	RootCmd.AddCommand(startCmd)

	ksync.LocatorFlags(startCmd, startViper)
}
