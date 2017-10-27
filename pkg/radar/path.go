package radar

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
)

// GetAbsPath takes a container path and returns the absolute path on the
// current node for that directory.
func (r *radarServer) GetBasePath(
	ctx context.Context,
	containerPath *pb.ContainerPath) (*pb.BasePath, error) {

	rootPath, err := getRootPath(containerPath)
	if err != nil {
		return nil, err
	}

	grpc_ctxtags.Extract(ctx).Set(
		"container", containerPath.ContainerId).Set(
		"rootPath", rootPath)

	log.WithFields(log.Fields{
		"path": rootPath,
	}).Debug("root path found")

	return &pb.BasePath{
		Full: rootPath}, nil
}
