package radar

import (
	// "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/vapor-ware/ksync/pkg/proto"
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

type radarVersion struct {
	Version   string
	GoVersion string
	GitCommit string
	GitTag    string
	BuildDate string
	Healthy   bool
}


func (r *radarServer) GetVersionInfo(ctx context.Context) (*pb.VersionInfo, error) {
	radar := radarVersion{
		Version:   VersionString,
		GoVersion: GoVersion,
		GitCommit: GitCommit,
		GitTag:    GitTag,
		BuildDate: BuildDate,
	}

	log.WithFields(log.Fields{
		"Version":   VersionString,
		"GoVersion": GoVersion,
		"GitCommit": GitCommit,
		"GitTag":    GitTag,
		"BuildDate": BuildDate,
	}).Debug("getting version info")

	return &pb.VersionInfo{
		Version:   VersionString,
		GoVersion: GoVersion,
		GitCommit: GitCommit,
		GitTag:    GitTag,
		BuildDate: BuildDate,
		}, nil
}
