package ksync

import (
	"fmt"
	"os"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/golang/protobuf/ptypes"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type FileList struct {
	Container *Container
	Path      string
	Files     *pb.Files
}

func (this *FileList) Get() error {
	conn, err := NewRadarConnection(this.Container.NodeName)
	if err != nil {
		return fmt.Errorf("Could not connect to radar: %v", err)
	}
	defer conn.Close()

	log.WithFields(log.Fields{
		"node": this.Container.NodeName,
		"pod":  this.Container.PodName,
		"name": this.Container.Name,
		"id":   this.Container.ID,
	}).Debug("radar connected")

	this.Files, err = pb.NewRadarClient(conn).ListContainerFiles(
		context.Background(), &pb.ContainerPath{this.Container.ID, this.Path})
	if err != nil {
		return fmt.Errorf("Could not list files: %v", err)
	}

	return nil
}

func (this *FileList) Output() error {

	fmt.Println(tm.Color(fmt.Sprintf("==> %s:%s:%s <==",
		this.Container.PodName, this.Container.Name, this.Path), tm.CYAN))

	// TODO: should output be configurable?
	// TODO: should this be a common table format?
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	// TODO: can I map this instead?
	// TODO: add color (directories, links, ...)
	for _, file := range this.Files.Items {
		modTime, _ := ptypes.Timestamp(file.ModTime)

		table.Append([]string{
			file.Mode,
			// TODO: make size human readable (via config?)
			fmt.Sprintf("%d", file.Size),
			modTime.Format("Jan 2 15:4"),
			// TODO: path output needs to be improved
			tm.Color(strings.TrimPrefix(file.Path, this.Path), this.pathColor(file)),
		})
	}
	table.Render()

	return nil
}

func (this *FileList) pathColor(file *pb.File) int {
	if file.IsDir {
		// TODO: this isn't the best blue ... is there a better way to handle this?
		return tm.BLUE
	}

	return tm.WHITE
}
