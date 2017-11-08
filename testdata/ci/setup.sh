#!/usr/bin/env bash

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

# Check if require utilities are installed and accessible
echo -e "${BLUE}Checking if kubectl and gcloud are installed${NC}"
if ! command -v kubectl; then
  echo -e "${RED}Kubectl Is not installed. It must be installed to run.${NC}"
  exit 1
elif ! command -v gcloud; then
  echo -e "${RED}Gcloud Is not installed. It must be installed to run.${NC}"
  exit 1
elif ! command -v jq; then
  echo -e "${RED}JQ is not installed. Attempting to install it for you.${NC}"
  curl -L --progress-bar https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 -o /usr/local/bin/jq
  chmod +x /usr/local/bin/jq
fi
echo -e "${GREEN}Kubectl, gcloud, and jq installed${NC}"

# Make sure we have the info for the correct cluster
echo -e "${BLUE}Checking for rights to the currently configured cluster${NC}"
gcloud container clusters get-credentials ${CLUSTER_NAME} --zone ${CLUSTER_ZONE}
echo -e "${GREEN}Got credentials from ${PURPLE}${CLUSTER_NAME}${NC}"

# Launch our test deployment
echo -e "${BLUE}Getting necessary image (${PURPLE}This will be removed when pulling is added)${NC}"
gcloud docker -- pull gcr.io/elated-embassy-152022/ksync/ksync:canary
echo -e "${BLUE}Launching test deployment${NC}"
kubectl apply -f ${CIRCLE_WORKING_DIRECTORY}/testdata/k8s/config/test-app.yaml --validate true -o json | jq
