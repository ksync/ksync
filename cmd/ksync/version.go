package main

import (
	"os"
	"runtime"
	"text/template"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
	pb "github.com/vapor-ware/ksync/pkg/proto"
	"github.com/vapor-ware/ksync/pkg/radar"
)

type versionCmd struct {
	cli.BaseCmd
}

func (v *versionCmd) new() *cobra.Command {
	long := `Print version information.`
	example := ``

	v.Init("ksync", &cobra.Command{
		Use:     "version",
		Short:   "Print version information.",
		Long:    long,
		Example: example,
		Run:     v.run,
	})

	return v.Cmd
}

var ksyncVersionTemplate = `{{define "ksync"}}ksync:
	Version:    {{.Client.Version}}
	Go Version: {{.Client.GoVersion}}
	Git Commit: {{.Client.GitCommit}}
	Git Tag:    {{if ne .Client.GitTag ""}}{{.Client.GitTag}}{{end}}
	Built:      {{.Client.BuildDate}}
	OS/Arch:    {{.Client.OS}}/{{.Client.Arch}}{{println}}{{end}}`

var radarVersionTemplate = `{{define "radar"}}radar:
	Version:    {{.Server.Version}}
	Go Version: {{.Server.GoVersion}}
	Git Commit: {{.Server.GitCommit}}
	Git Tag:    {{if ne .Server.GitTag ""}}{{.Server.GitTag}}{{end}}
	Built:      {{.Server.BuildDate}}
	Healthy:    {{.Server.Healthy}}{{println}}{{end}}`

type versionInfo struct {
	Client *ksync.Version
	Server *radar.Version
}

func (v *versionCmd) run(cmd *cobra.Command, args []string) { // nolint: gocyclo
	template, err := template.New("ksync").Parse(ksyncVersionTemplate)
	if err != nil {
		log.Fatal(err)
	}

	template, err = template.New("radar").Parse(radarVersionTemplate)
	if err != nil {
		log.Fatal(err)
	}

	version := versionInfo{
		Client: &ksync.Version{
			Version:   ksync.VersionString,
			GoVersion: ksync.GoVersion,
			GitCommit: ksync.GitCommit,
			GitTag:    ksync.GitTag,
			BuildDate: fixVersionTime(ksync.BuildDate),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
	}

	err = template.ExecuteTemplate(os.Stdout, "ksync", version)
	if err != nil {
		log.Fatal(err)
	}

	radarVersion, err := radarVersion()
	if err != nil {
		log.Fatal(err)
	}

	if radarVersion == nil {
		return
	}

	version.Server = &radar.Version{
		Version:   radarVersion.Version,
		GoVersion: radarVersion.GoVersion,
		GitCommit: radarVersion.GitCommit,
		GitTag:    radarVersion.GitTag,
		BuildDate: fixVersionTime(radarVersion.BuildDate),
		Healthy:   true,
	}

	if err := template.ExecuteTemplate(os.Stdout, "radar", version); err != nil {
		log.Fatal(err)
	}
}

func fixVersionTime(version string) string {
	timeRadar, err := time.Parse(time.RFC3339, version)
	if err != nil {
		return ""
	}

	return timeRadar.Format(time.UnixDate)
}

// TODO: temporary
func radarVersion() (*pb.VersionInfo, error) {
	service := cluster.NewService()
	nodes, err := service.NodeNames()
	if err != nil {
		return nil, nil
	}

	if len(nodes) == 0 {
		return nil, nil
	}

	if _, healthErr := service.IsHealthy(nodes[0]); err != nil {
		return nil, healthErr
	}

	conn, err := cluster.NewConnection(nodes[0]).Radar()
	if err != nil {
		return nil, err
	}

	return pb.NewRadarClient(conn).GetVersionInfo(
		context.Background(), &empty.Empty{})
}
