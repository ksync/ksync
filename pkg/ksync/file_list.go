package ksync

import (
	"fmt"
	"os"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/golang/protobuf/ptypes"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// FileList is the set of all Files that exist in the remote Container.
type FileList struct {
	Container *Container
	Path      string
	Files     *pb.Files
}

// Get populates the FileList with the set of all Files in the remote Container.
func (f *FileList) Get() error {
	client, err := f.Container.Radar()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	f.Files, err = client.ListContainerFiles(
		context.Background(), &pb.ContainerPath{
			ContainerId: f.Container.ID,
			PathName:    f.Path,
		})
	if err != nil {
		return fmt.Errorf("Could not list files: %v", err)
	}

	return nil
}

// Output prints a table of the Files in f FileList.
func (f *FileList) Output() error {

	fmt.Println(tm.Color(fmt.Sprintf("==> %s:%s:%s <==",
		f.Container.PodName, f.Container.Name, f.Path), tm.CYAN))

	// TODO: should output be configurable?
	// TODO: should this be a common table format?
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	// TODO: can I map this instead?
	// TODO: add color (directories, links, ...)
	for _, file := range f.Files.Items {
		modTime, _ := ptypes.Timestamp(file.ModTime)

		// TODO: show link path eg. foo -> ../bar
		table.Append([]string{
			file.Mode,
			// TODO: make size human readable (via config?)
			fmt.Sprintf("%d", file.Size),
			modTime.Format("Jan 2 15:4"),
			// TODO: path output needs to be improved
			tm.Color(strings.TrimPrefix(file.Path, f.Path), f.pathColor(file)),
		})
	}
	table.Render()

	return nil
}

func (f *FileList) pathColor(file *pb.File) int {
	if file.IsDir {
		// TODO: this isn't the best blue ... is there a better way to handle this?
		return tm.BLUE
	}

	// TODO: color links cyan

	return tm.WHITE
}
