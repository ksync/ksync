#!/usr/bin/env bash

set -x
# set -eo pipefail

# Stamps to be evaluated at build time. They will be incorporated into the build
# via ldflags

BINARY_VERSION=${BINARY_VERSION:-"corrupted version"}
GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD 2> /dev/null || true)}
# This requires gnu-date!
BUILD_DATE=${BUILD_DATE:-$(date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')}

# Setup ldflags for runs
export LDFLAGS="\
    -w \
    -X github.com/vapor-ware/ksync/pkg/ksync/ksync.GitCommit=${GIT_COMMIT} \
    -X github.com/vapor-ware/ksync/pkg/ksync/ksync.BuildDate=${BUILD_DATE} \
    -X github.com/vapor-ware/ksync/pkg/ksync/ksync.VersionString=${BINARY_VERSION} \
    ${LDFLAGS:-} \
"

echo "Set stamps"
