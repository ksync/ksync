#!/usr/bin/env bash

# set -x
# set -eo pipefail

# Stamps to be evaluated at build time. They will be incorporated into the build
# via ldflags

# BINARY_VERSION is currently inherited from an organization wide context in
# CI. It is set to "Release" for all tagged releases for now. It can be separated
# into "dev", "edge", etc. in the future via more fine grained substitution
BINARY_VERSION=${BINARY_VERSION:-"corrupted-version"}
GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
# Deal with macOS
if (( $(uname) == "Darwin" )); then
  BUILD_DATE_UNIX=${BUILD_DATE_UNIX:-$(date -u +%Y-%m-%dT%T 2> /dev/null)}
  BUILD_DATE_OFFSET=${BUILD_DATE_OFFSET:-"+00:00"}
  BUILD_DATE_FAKE_NANOSECONDS=${BUILD_DATE_FAKE_NANOSECONDS:-".000000000"}
  BUILD_DATE=${BUILD_DATE:-${BUILD_DATE_UNIX}${BUILD_DATE_FAKE_NANOSECONDS}${BUILD_DATE_OFFSET}}
fi
# This requires gnu-date!
BUILD_DATE=${BUILD_DATE:-$(date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')}
# This seems nasty
GO_VERSION=${GO_VERSION:-$(go version | awk '{ print $3 }')}

# Setup ldflags for runs
export LDFLAGS="\
    -w \
    -X github.com/vapor-ware/ksync/pkg/ksync.GitCommit=${GIT_COMMIT} \
    -X github.com/vapor-ware/ksync/pkg/ksync.BuildDate=${BUILD_DATE} \
    -X github.com/vapor-ware/ksync/pkg/ksync.VersionString=${BINARY_VERSION} \
    -X github.com/vapor-ware/ksync/pkg/ksync.GoVersion=${GO_VERSION} \
    -X github.com/vapor-ware/ksync/pkg/ksync.GitTag=${CIRCLE_TAG} \
    -X github.com/vapor-ware/ksync/pkg/radar.GitCommit=${GIT_COMMIT} \
    -X github.com/vapor-ware/ksync/pkg/radar.BuildDate=${BUILD_DATE} \
    -X github.com/vapor-ware/ksync/pkg/radar.VersionString=${BINARY_VERSION} \
    -X github.com/vapor-ware/ksync/pkg/radar.GoVersion=${GO_VERSION} \
    -X github.com/vapor-ware/ksync/pkg/radar.GitTag=${CIRCLE_TAG} \
    ${LDFLAGS:-} \
"

echo "Set stamps"
