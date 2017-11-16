package main

import (
	"os"
	"runtime"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vapor-ware/ksync/pkg/cli"
	"github.com/vapor-ware/ksync/pkg/ksync"
)

var (
	// GitCommit is the commit hash of the commit used to build
	GitCommit string
	// VersionString is the canonical version string
	VersionString string
	// BuildDate contains the build timestamp
	BuildDate string
	// GitTag optionally contains the git tag used in build
	GitTag string
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
	Client ksyncVersion
	Server radarVersion
}

type ksyncVersion struct {
	Version   string
	GoVersion string
	GitCommit string
	GitTag    string
	BuildDate string
	OS        string
	Arch      string
}

type radarVersion struct {
	Version   string
	GoVersion string
	GitCommit string
	GitTag    string
	BuildDate string
	Healthy   bool
}

func (v *versionCmd) run(cmd *cobra.Command, args []string) {
	templateKsync, err := template.New("ksync").Parse(ksyncVersionTemplate)
	if err != nil {
		log.Fatal(err)
	}
	templateRadar, err := template.New("radar").Parse(radarVersionTemplate)
	if err != nil {
		log.Fatal(err)
	}

	version := versionInfo{
		Client: ksyncVersion{
			Version:   VersionString,
			GoVersion: runtime.Version(),
			GitCommit: GitCommit,
			GitTag:    GitTag,
			BuildDate: BuildDate,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
		// TODO: get this from radar
		Server: radarVersion{
			Version:   VersionString,
			GoVersion: runtime.Version(),
			GitCommit: GitCommit,
			GitTag:    GitTag,
			BuildDate: BuildDate,
			Healthy:   radarCheck(),
		},
	}

	// Convert time to a human readable format
	timeKsync, timeErr := time.Parse(time.RFC3339Nano, version.Client.BuildDate)
	if timeErr == nil {
		version.Client.BuildDate = timeKsync.Format(time.UnixDate)
	} else {
		log.Fatal(timeErr)
	}

	timeRadar, timeErr := time.Parse(time.RFC3339, version.Server.BuildDate)
	if timeErr == nil {
		version.Server.BuildDate = timeRadar.Format(time.UnixDate)
	} else {
		log.Fatal(timeErr)
	}

	// If radar is reachable, print that part of the template
	// TODO: Change this to use template.ExecuteTemplate
	if radarCheck() {
		err := templateKsync.Execute(os.Stdout, version)
		if err != nil {
			log.Fatal(err)
		}
		err = templateRadar.Execute(os.Stdout, version)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := templateKsync.Execute(os.Stdout, version)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// TODO: temporary
func radarCheck() bool {
	radar := ksync.NewRadarInstance()
	containers, err := ksync.GetRemoteContainers("", "app=test", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(containers)
	h, err := radar.IsHealthy(containers[0].NodeName)
	log.WithError(err)
	return h
}
