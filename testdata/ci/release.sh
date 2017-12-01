#!/usr/bin/env bash

# set -x
set -eo pipefail

# Colors for reporting
RED='\033[0;31m'
GREEN='\033[0;32m'
PURPLE='\033[0;35m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Check if this is a tag we're building. If it is then push an unpublished release
echo -e "${BLUE}Checking if this should be a release${NC}"
if ${CIRCLE_TAG} == ""; then
  echo -e "${YELLOW}No tag detected. Not pushing a release.${NC}"
  exit 0
else
  echo -e "${GREEN} Tag detected. Creating release for ${CIRCLE_TAG}.${NC}"
fi

# Check if `ghr` (https://github.com/tcnksm/ghr) is installed.
echo -e "${BLUE}Checking if GHR is installed${NC}"
if ! command -v ghr; then
  echo -e "${RED}GHR Is not installed. It must be installed to run.${NC}"
  exit 1
elif ${GITHUB_TOKEN} == ""; then
  echo -e "${RED}No GitHub token is set! A token must be passed to upload releases.${NC}"
  exit 1
else
  echo -e "${GREEN}GHR Is not installed. It must be installed to run.${NC}"
fi

# Create a release for the given tag and push it
echo -e "${BLUE}Tag: ${CIRCLE_TAG}${NC}"
echo -e "${BLUE}Commit: ${CIRCLE_SHA1}${NC}"
echo -e "${BLUE}Changes: ${CIRCLE_COMPARE_URL}${NC}"

ghr \
  -t ${GITHUB_TOKEN} \
  -p 5 \
  -draft \
  ${CIRCLE_TAG} bin/
