#!/usr/bin/env bash

# set -x
# set -eo pipefail

# Stamps to be evaluated at build time. They will be incorporated into the build
# via ldflags

BINARY_VERSION=${BINARY_VERSION:-"corrupted-version"}
GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
# This requires gnu-date!
BUILD_DATE=${BUILD_DATE:-$(date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')}
# This seems nasty
GO_VERSION=${GO_VERSION:-$(go version | awk '{ print $3 }')}

# Setup ldflags for runs
# THIS MUST BE RUN FROM THE ROOT REPO DIR!
export LDFLAGS="\
    -w \
    -X ksync.GitCommit=${GIT_COMMIT} \
    -X ksync.BuildDate=${BUILD_DATE} \
    -X ksync.VersionString=${BINARY_VERSION} \
    -X ksync.GoVersion=${GO_VERSION} \
    -X radar.GitCommit=${GIT_COMMIT} \
    -X radar.BuildDate=${BUILD_DATE} \
    -X radar.VersionString=${BINARY_VERSION} \
    -X radar.GoVersion=${GO_VERSION} \
    ${LDFLAGS:-} \
"

echo "Set stamps"
