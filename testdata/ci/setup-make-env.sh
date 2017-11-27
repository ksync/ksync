#!/usr/bin/env bash
# TODO: This is all duplication and should be removed for something sleeker

set -x
set -eo pipefail

# Colors for reporting
RED='\033[0;31m'
GREEN='\033[0;32m'
PURPLE='\033[0;35m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Check if we are running in CircleCI
if [[ -z $CIRCLECI ]]; then
  echo -e "${RED}We are not running on CircleCI! We must be running there.${NC}"
  exit 1
fi

# Add google cli utilites to our path if necessary
source /home/circleci/google-cloud-sdk/path.bash.inc || echo -e "${BLUE}Google install path not dectected, not modifying path${NC}"
