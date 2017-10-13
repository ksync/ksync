package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vapor-ware/ksync/pkg/ksync"
)

type WatchCmd struct {
	viper *viper.Viper
}

func (this *WatchCmd) New() *cobra.Command {
	long := `Watch configured syncs and start them when required.`
	example := ``

	cmd := &cobra.Command{
		Use:     "watch",
		Short:   "Watch configured syncs and start them when required.",
		Long:    long,
		Example: example,
		Run:     this.run,
	}

	this.viper = viper.New()

	return cmd
}

func (this *WatchCmd) run(cmd *cobra.Command, args []string) {
	// 1. Watch config file for updates
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		this.manageSpecs()
	})

	this.manageSpecs()

	// 2. Watch k8s API for updates
	// 3. Add/remove runs

	waitForSignal()
}

// TODO: how to test what *shouldn't* be running?
func (this *WatchCmd) manageSpecs() {
	specMap, _ := ksync.AllSpecs()
	for name, spec := range specMap.Items {
		// Should run?
		containerList, err := ksync.GetContainers(
			spec.Pod, spec.Selector, spec.Container)
		if err != nil {
			log.Fatal(err)
		}

		if len(containerList) == 0 {
			log.WithFields(spec.Fields()).Debug("no matching running containers.")
			continue
		}

		for _, cntr := range containerList {
			service, err := ksync.NewService(name, cntr, spec)
			if err != nil {
				log.Fatal(err)
			}

			status, err := service.Status()
			if err != nil {
				log.Fatal(err)
			}

			log.WithFields(status.Fields()).Debug("service status")

			// Container is running, leave it alone.
			if status.Running {
				continue
			}

			if err := service.Start(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func waitForSignal() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
