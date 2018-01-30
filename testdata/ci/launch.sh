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

# Install the binaries
# TODO: This would seem to be a nonsequitor as we have to compile and launch
# the binaries _before_ running tests. Need to figure out a better way to do this.
cd ${CIRCLE_WORKING_DIRECTORY}/cmd/ksync
go install -v
cd ${CIRCLE_WORKING_DIRECTORY}/cmd/radar
go install -v

# Deploy radar to the cluster
ksync init --local=false

# Get absolute path for kubectl in case it isn't in our shell
TEST_KUBECTL="/home/circleci/google-cloud-sdk/bin/kubectl"

# Set a namespace to use
TEST_NAMESPACE="default"
TEST_RADAR_NAMESPACE="kube-system"

timeout -k 2m 2m ${TEST_KUBECTL} run --rm -it wait-for-ksync \
  --restart Never \
  --image=groundnuty/k8s-wait-for \
  --requests='cpu=10m' \
  -- pod -lapp=ksync --all-namespaces

# Fetch the names of the pods launched
PODS=($(${TEST_KUBECTL} get pods --namespace ${TEST_NAMESPACE} --selector app=test -o json | jq --raw-output '.items[].metadata.name'))
echo -e "${BLUE}${PODS[@]}${NC}"

# Fetch the nodes pods are scheduled on
NODES=($(${TEST_KUBECTL} get pods --namespace ${TEST_NAMESPACE} --selector app=test -o json | jq --raw-output '.items[].spec.nodeName'))
echo -e "${BLUE}${NODES[@]}${NC}"

# Fetch the pod's ID
CONTAINERIDS=($(${TEST_KUBECTL} get pods --namespace ${TEST_NAMESPACE} --selector app=test -o json | jq --raw-output '.items[].status.containerStatuses[].containerID' | awk '{print $NF}' FS=/))
echo -e "${BLUE}${CONTAINERIDS[@]}${NC}"

# Fetch the names of the pods launched
RADAR_PODS=($(${TEST_KUBECTL} get pods --namespace ${TEST_RADAR_NAMESPACE} --selector app=radar -o json | jq --raw-output '.items[].metadata.name'))
echo -e "${BLUE}${RADAR_PODS[@]}${NC}"

# Fetch the nodes pods are scheduled on
RADAR_NODES=($(${TEST_KUBECTL} get pods --namespace ${TEST_RADAR_NAMESPACE} --selector app=radar -o json | jq --raw-output '.items[].spec.nodeName'))
echo -e "${BLUE}${RADAR_NODES[@]}${NC}"

# Fetch the pod's ID
RADAR_CONTAINERIDS=($(${TEST_KUBECTL} get pods --namespace ${TEST_RADAR_NAMESPACE} --selector app=radar -o json | jq --raw-output '.items[].status.containerStatuses[].containerID' | awk '{print $NF}' FS=/))
echo -e "${BLUE}${RADAR_CONTAINERIDS[@]}${NC}"

# Export all necessary variables to a temporary holding place for later sourcing
echo "export TEST_NAMESPACE=${TEST_NAMESPACE}" >> ${BASH_ENV}
echo "export TEST_POD=${PODS[0]}" >> ${BASH_ENV}
echo "export TEST_CONTAINERID=${CONTAINERIDS[0]}" >> ${BASH_ENV}
echo "export TEST_NODE=${NODES[0]}" >> ${BASH_ENV}
echo "export TEST_RADAR_NAMESPACE=${TEST_RADAR_NAMESPACE}" >> ${BASH_ENV}
echo "export TEST_RADAR_POD=${RADAR_PODS[0]}" >> ${BASH_ENV}
echo "export TEST_RADAR_CONTAINERID=${RADAR_CONTAINERIDS[0]}" >> ${BASH_ENV}
echo "export TEST_RADAR_NODE=${RADAR_NODES[0]}" >> ${BASH_ENV}
