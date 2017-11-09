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
	long := `Watch configured syncs and start them when required.

	Note: this is run automatically for you by 'ksync init'. It expects to be run
	from inside a container.`
	example := ``

	w.Init("ksync", &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured syncs and start them when required.",
		Long:    long,
		Example: example,
		Run:     w.run,
		// TODO: remove this when the command can be run by users.
		Hidden: true,
	})

	return w.Cmd
}

// TODO: hook up to k8s and watch for changes
// TODO: handle Normalize errors.
func (w *watchCmd) run(cmd *cobra.Command, args []string) {
	// 1. Watch config file for updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		list, err := ksync.GetServices()
		if err != nil {
			log.Fatal(err)
		}
		if normerr := list.Normalize(); normerr != nil {
			log.Fatal(normerr)
		}
	})

	list, err := ksync.GetServices()
	if err != nil {
		log.Fatal(err)
	}
	if normerr := list.Normalize(); normerr != nil {
		log.Fatal(normerr)
	}
	// 2. Watch k8s API for updates

	waitForSignal()
}

func waitForSignal() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM) // nolint: megacheck
	<-exitSignal
}
