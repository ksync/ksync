package radar

import (
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// GetAbsPath takes a container path and returns the absolute path on the
// current node for that directory.
func (r *radarServer) GetAbsPath(
	ctx context.Context,
	containerPath *pb.ContainerPath) (*pb.AbsolutePath, error) {

	rootPath, err := GetRootPath(containerPath)
	if err != nil {
		return nil, err
	}

	grpc_ctxtags.Extract(ctx).Set(
		"container", containerPath.ContainerId).Set(
		"path", containerPath.PathName).Set(
		"rootPath", rootPath)

	log.WithFields(log.Fields{
		"path": rootPath,
	}).Debug("root path found")

	joinPath := filepath.Join(rootPath, containerPath.PathName)

	if _, err := os.Lstat(joinPath); err != nil {
		return nil, err
	}

	return &pb.AbsolutePath{
		Full: joinPath}, nil
}
