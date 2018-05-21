package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/vapor-ware/ksync/pkg/cli"
	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type getCmd struct {
	cli.BaseCmd
}

func (g *getCmd) new() *cobra.Command {
	long := `Get all specs.

	Based off what specs have been created, returns the current status of each spec.`
	example := ``

	g.Init("ksync", &cobra.Command{
		Use:     "get",
		Short:   "Get all specs.",
		Long:    long,
		Example: example,
		Run:     g.run,
	})

	return g.Cmd
}

func (g *getCmd) out(specs *pb.SpecList) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetHeader([]string{"Name", "Local", "Remote", "Status", "Pod", "Container"})

	var keys []string
	for name := range specs.Items {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		spec := specs.Items[name]

		status := spec.Status
		if len(spec.Services.Items) > 0 {
			status = ""
		}

		relPath, err := filepath.Rel(cwd, spec.Details.LocalPath)
		if err != nil {
			log.Fatal(err)
		}

		local := relPath
		if strings.Count(local, "/") > strings.Count(spec.Details.LocalPath, "/") {
			local = spec.Details.LocalPath
		}

		// Print "read-only" status only if it is set
		if spec.Details.LocalReadOnly {
			local = fmt.Sprintf("%s:%s", local, "ro")
		}

		var remote string
		if spec.Details.RemoteReadOnly {
			remote = fmt.Sprintf("%s:%s", spec.Details.RemotePath, "ro")
		} else {
			remote = spec.Details.RemotePath
		}

		table.Append([]string{
			name,
			local,
			remote,
			status,
		})

		for _, service := range spec.Services.Items {
			table.Append([]string{
				"",
				"",
				"",
				service.Status,
				service.RemoteContainer.PodName,
				spec.Details.ContainerName,
			})
		}
	}

	table.Render()
}

func (g *getCmd) run(cmd *cobra.Command, args []string) {
	// This is connecting locally and it is very unlikely watch is overloaded,
	// set the timeout *super* short to make it easier on the users when they
	// forgot to start watch.
	withTimeout, _ := context.WithTimeout(context.TODO(), 100 * time.Millisecond)

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

	g.out(resp)
}
