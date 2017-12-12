package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	// "github.com/vapor-ware/ksync/pkg/ksync"
)

type getCmd struct {
	cli.BaseCmd
}

func (g *getCmd) new() *cobra.Command {
	long := `Get all configured syncs and their status.`
	example := ``

	g.Init("ksync", &cobra.Command{
		Use:     "get",
		Short:   "Get all configured syncs and their status.",
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
	table.SetHeader([]string{"Name", "Local", "Remote", "Status", "Pod"})

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

		local, err := filepath.Rel(cwd, spec.Details.LocalPath)
		if err != nil {
			log.Fatal(err)
		}

		table.Append([]string{
			name,
			local,
			spec.Details.RemotePath,
			status,
		})

		for _, service := range spec.Services.Items {
			table.Append([]string{
				"",
				"",
				"",
				service.Status,
				service.RemoteContainer.PodName,
			})
		}
	}

	table.Render()
}

// TODO: TLS?
func (g *getCmd) run(cmd *cobra.Command, args []string) {
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%d", viper.GetInt("port")),
		[]grpc.DialOption{
			grpc.WithTimeout(5 * time.Second),
			grpc.WithBlock(),
			grpc.WithInsecure(),
		}...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() // nolint: errcheck

	client := pb.NewKsyncClient(conn)

	resp, err := client.GetSpecList(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	g.out(resp)
}
