package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

type watchCmd struct {
	cli.BaseCmd
}

func (w *watchCmd) new() *cobra.Command {
	long := `Watch configured syncs and start them when required.`
	example := ``

	w.Init("ksync", &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured syncs and start them when required.",
		Long:    long,
		Example: example,
		Run:     w.run,
	})

	return w.Cmd
}

// TODO: hook up to k8s and watch for changes
// TODO: handle Normalize errors.
func (w *watchCmd) run(cmd *cobra.Command, args []string) {
	// 1. Watch config file for updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := ksync.GetServices().Normalize(); err != nil {
			log.Fatal(err)
		}
	})

	if err := ksync.GetServices().Normalize(); err != nil {
		log.Fatal(err)
	}
	// 2. Watch k8s API for updates

	waitForSignal()
}

func waitForSignal() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM) // nolint: megacheck
	<-exitSignal
}
