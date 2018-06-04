package main

import (
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type reloadCmd struct {
	cli.BaseCmd
}

func (r *reloadCmd) new() *cobra.Command {
	long := `Reload one or more remote specs.

	Initiates a manual reload of the remote side of one or more specs.`
	example := ``

	r.Init("ksync", &cobra.Command{
		Use:     "reload",
		Short:   "Reload a remote spec.",
		Long:    long,
		Example: example,
		Run:     r.run,
	})

	flags := r.Cmd.Flags()

	flags.BoolP(
		"all",
		"a",
		false,
		"reload all specs")
	if err := r.BindFlag("all"); err != nil {
		log.Fatal(err)
	}

	return r.Cmd
}

func (r *reloadCmd) run(cmd *cobra.Command, args []string) {
	if r.Viper.GetBool("all") {
		if len(r.Cmd.Flags().Args()) == 0 {
			r.reloadAll()
		}
		log.Fatal("cannot specify names when using `--all`")
	} else if len(r.Cmd.Flags().Args()) != 0 {
		sort.Strings(args)
		for specName := range args {
			r.reload(args[specName])
		}
	} else {
		log.Fatal("reload requires at least one spec or `--all`")
	}
}

func (r *reloadCmd) reload(specName string) {
	specs := &ksync.SpecList{}
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	if !specs.Has(specName) {
		log.Fatalf("%s does not exist. Did you mean something else?", specName)
	}

	log.Debugf("attempting reload of %s", specName)

	spec, err := specs.Get(specName)
	log.Debugf("%+v", spec)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("%s", )
	service, err := spec.Services.Get(spec.Details.Name)
	if err != nil {
		log.Fatal(err)
	}

	containerID := service.RemoteContainer.ID

	conn, err := cluster.NewConnection(service.RemoteContainer.NodeName).Radar()
	client := pb.NewRadarClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Restart(context.Background(), &pb.ContainerPath{
		ContainerId: containerID,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("restart initiated for %s", specName)
}

func (r *reloadCmd) reloadAll() {
	specs := &ksync.SpecList{}
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	for name, _ := range specs.Items {
		r.reload(name)
	}
}
