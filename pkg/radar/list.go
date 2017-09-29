package radar

import (
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

type ContainerFileList struct {
	ContainerPath *pb.ContainerPath
	Files         *pb.Files
}

func (this *ContainerFileList) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.WithFields(log.Fields{
			"path":  path,
			"error": err,
		}).Debug("could not lookup path")
	}

	return nil
}

// TODO: is there a better way to pass errors back than just pushing the string?
func (s *radarServer) ListContainerFiles(
	ctx context.Context,
	containerPath *pb.ContainerPath) (*pb.Files, error) {

	fileList := &ContainerFileList{containerPath, &pb.Files{}}

	grpc_ctxtags.Extract(ctx).Set(
		"container", fileList.ContainerPath.ContainerId).Set(
		"path", fileList.ContainerPath.PathName)

	// TODO: need some kind of error handling here
	rootPath, err := GetRootPath(fileList.ContainerPath)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"path": rootPath,
	}).Debug("root path found")

	joinPath := filepath.Join(rootPath, containerPath.PathName)

	filepath.Walk(joinPath, fileList.walk)

	return &pb.Files{}, nil
}
