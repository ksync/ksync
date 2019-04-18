package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type reloadCmd struct {
	cli.BaseCmd
}

func (r *reloadCmd) new() *cobra.Command {
	long := `Reload one or more remote specs.

	Initiates a manual reload of the remote side of one or more specs.`
	example := `ksync reload --all`

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
		} else {
			log.Fatal("cannot specify names when using `--all`")
		}
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
	// This is connecting locally and it is very unlikely watch is overloaded,
	// set the timeout *super* short to make it easier on the users when they
	// forgot to start watch.
	withTimeout, _ := context.WithTimeout(context.TODO(), 100*time.Millisecond)

	grpcConnection, err := grpc.DialContext(
		withTimeout,
		fmt.Sprintf("127.0.0.1:%d", viper.GetInt("port")),
		[]grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
		}...)
	if err != nil {
		// The assumption is that the only real error here is because watch isn't
		// running
		log.Debug(err)
		log.Fatal("Having problems querying status. Are you running watch?")
	}
	defer grpcConnection.Close() // nolint: errcheck

	ksyncClient := pb.NewKsyncClient(grpcConnection)

	resp, err := ksyncClient.GetSpecList(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	specs, err := ksync.DeserializeSpecList(resp)
	if err != nil {
		log.Fatal(err)
	}

	if !specs.Has(specName) {
		log.Fatalf("%s does not exist. Did you mean something else?", specName)
	}

	log.Infof("attempting reload of %s", specName)

	spec := specs.Items[specName]

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

}

func (r *reloadCmd) reloadAll() {
	specs := ksync.NewSpecList()
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	for name := range specs.Items {
		r.reload(name)
	}
}
