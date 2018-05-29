package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type deleteCmd struct {
	cli.BaseCmd
}

func (d *deleteCmd) new() *cobra.Command {
	long := `Delete an existing spec. This will stop syncing files between your
	local directory and the remote containers.

	The files you've synced are not touched and the remote container is left as is.`
	example := ``

	d.Init("ksync", &cobra.Command{
		Use:     "delete [flags] [name]...",
		Short:   "Delete an existing spec",
		Long:    long,
		Example: example,
		Aliases: []string{"d"},
		Args:    cobra.ArbitraryArgs(),
		Run:     d.run,
	})

	flags := d.Cmd.Flags()

	flags.Bool(
		"all",
		false,
		"delete all specs")
	if err := d.BindFlag("all"); err != nil {
		log.Fatal(err)
	}

	return d.Cmd
}

func (d *deleteCmd) run(cmd *cobra.Command, args []string) {
	if viper.GetBool("all") {
		if len(cmd.Flags().Args()) == 0 {
			d.deleteAll()
		}
		log.Fatal("cannot specify names when using `--all`")
	} else {
		for name := range args {
			d.delete(args[name])
		}
	}
}

func (d *deleteCmd) delete(name string) {
	specs := &ksync.SpecList{}
	if err := specs.Update(); err != nil {
		log.Fatal(err)
	}

	if !specs.Has(name) {
		log.Fatalf("%s does not exist. Did you mean something else?", name)
	}

	if err := specs.Delete(name); err != nil {
		log.Fatalf("Could not delete %s: %v", name, err)
	}

	if err := specs.Save(); err != nil {
		log.Fatal(err)
	}
}

func (d *deleteCmd) deleteAll() {
	// This is connecting locally and it is very unlikely watch is overloaded,
	// set the timeout *super* short to make it easier on the users when they
	// forgot to start watch.
	withTimeout, _ := context.WithTimeout(context.TODO(), 100*time.Millisecond)

	conn, err := grpc.DialContext(
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
	defer conn.Close() // nolint: errcheck

	client := pb.NewKsyncClient(conn)

	resp, err := client.GetSpecList(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for name := range resp.Items {
		d.delete(name)
	}
}
