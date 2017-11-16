package ksync

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

type ksyncVersion struct {
	Version   string
	GoVersion string
	GitCommit string
	GitTag    string
	BuildDate string
	OS        string
	Arch      string
}
