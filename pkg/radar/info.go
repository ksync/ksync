package radar

import (
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/ksync/ksync/pkg/debug"
	pb "github.com/ksync/ksync/pkg/proto"
)

// These values will be stamped at build time
var (
	// GitCommit is the commit hash of the commit used to build
	GitCommit string
	// VersionString is the canonical version string
	VersionString string
	// BuildDate contains the build timestamp
	BuildDate string
	// GitTag optionally contains the git tag used in build
	GitTag string
	// GoVersion contains the Go version used in build
	GoVersion string
)

func (r *radarServer) GetVersionInfo(
	ctx context.Context, _ *empty.Empty) (*pb.VersionInfo, error) {

	version := &pb.VersionInfo{
		Version:   VersionString,
		GoVersion: GoVersion,
		GitCommit: GitCommit,
		GitTag:    GitTag,
		BuildDate: BuildDate,
	}

	log.WithFields(debug.StructFields(version)).Debug("getting version info")

	return version, nil
}
