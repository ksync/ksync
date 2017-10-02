package radar

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type ContainerFileList struct {
	ContainerPath *pb.ContainerPath
	Files         *pb.Files

	rootPath string
}

func (this *ContainerFileList) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.WithFields(log.Fields{
			"path":  path,
			"error": err,
		}).Debug("could not lookup path")

		return err
	}

	// No idea why this would ever fail ...
	modTime, _ := ptypes.TimestampProto(info.ModTime())

	this.Files.Items = append(this.Files.Items, &pb.File{
		strings.TrimPrefix(path, this.rootPath),
		info.Size(),
		info.Mode().String(),
		modTime,
		info.IsDir(),
	})

	return nil
}

// TODO: is there a better way to pass errors back than just pushing the string?
func (this *radarServer) ListContainerFiles(
	ctx context.Context,
	containerPath *pb.ContainerPath) (*pb.Files, error) {

	// TODO: need some kind of error handling here
	rootPath, err := GetRootPath(containerPath)
	if err != nil {
		return nil, err
	}

	fileList := &ContainerFileList{containerPath, &pb.Files{}, rootPath}

	grpc_ctxtags.Extract(ctx).Set(
		"container", fileList.ContainerPath.ContainerId).Set(
		"path", fileList.ContainerPath.PathName).Set(
		"rootPath", rootPath)

	log.WithFields(log.Fields{
		"path": rootPath,
	}).Debug("root path found")

	joinPath := filepath.Join(rootPath, containerPath.PathName)

	filepath.Walk(joinPath, fileList.walk)

	return fileList.Files, nil
}
