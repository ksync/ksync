package main

import (
	"os"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ksync/ksync/pkg/cli"
	"github.com/ksync/ksync/pkg/ksync"
	"github.com/ksync/ksync/pkg/ksync/cluster"
)

type versionCmd struct {
	cli.BaseCmd
}

func (v *versionCmd) new() *cobra.Command {
	long := `View the versions of both the local binary and remote service.`
	example := `ksync version`

	v.Init("ksync", &cobra.Command{
		Use:     "version",
		Short:   "View the versions of both the local binary and remote service.",
		Long:    long,
		Example: example,
		Run:     v.run,
	})

	return v.Cmd
}

var versionTemplate = `{{define "local"}}ksync:
	Version:    {{.Version}}
	Go Version: {{.GoVersion}}
	Git Commit: {{.GitCommit}}
	Git Tag:    {{if ne .GitTag ""}}{{.GitTag}}{{end}}
	Built:      {{date .BuildDate}}
	OS/Arch:    {{.OS}}/{{.Arch}}
{{end}}

{{define "service"}}service:
	Version:    {{.Version}}
	Go Version: {{.GoVersion}}
	Git Commit: {{.GitCommit}}
	Git Tag:    {{if ne .GitTag ""}}{{.GitTag}}{{end}}
	Built:      {{ date .BuildDate}}
{{end}}`

func parseDate(version string) string {
	t, err := time.Parse(time.RFC3339, version)
	if err != nil {
		return ""
	}

	return t.Format(time.UnixDate)
}

func (v *versionCmd) run(cmd *cobra.Command, args []string) {
	tmpl, tmplErr := template.New("local").Funcs(template.FuncMap{
		"date": parseDate,
	}).Parse(versionTemplate)
	if tmplErr != nil {
		log.Fatal(tmplErr)
	}

	if err := tmpl.ExecuteTemplate(
		os.Stdout, "local", ksync.Version()); err != nil {
		log.Fatal(err)
	}

	service := cluster.NewService()
	state, err := service.IsInstalled()
	if err != nil {
		log.Fatal(err)
	}

	if !state {
		log.Print("Remote service is not running. Run init to start it.")
	}

	radarVersion, err := service.Version()
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.ExecuteTemplate(
		os.Stdout, "service", radarVersion); err != nil {
		log.Fatal(err)
	}
}
