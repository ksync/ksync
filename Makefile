
BINARY_VERSION  ?= corrupted-version
GIT_COMMIT      ?= $(shell git rev-parse --short HEAD)

DATE            := $(shell (which gdate > /dev/null && echo "gdate") || echo "date")
BUILD_DATE      := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
GO_VERSION      := $(shell go version | awk '{ print $$3 }')

IMAGE_BASE      := ksync/ksync
MUTABLE_VERSION := canary
DOCKER_VERSION  := git-${GIT_COMMIT}
export IMAGE    := ${IMAGE_BASE}:${DOCKER_VERSION}
MUTABLE_IMAGE   := ${IMAGE_BASE}:${MUTABLE_VERSION}

# ifdef CIRCLE_TAG
# IMAGE_TAG       := ${IMAGE_BASE}:${CIRCLE_TAG}
# endif

# CMD       ?= bin/ksync --log-level=debug watch

GO        ?= go
GOBIN     ?= $(shell go env GOPATH)/bin

TAGS      :=
LDFLAGS   := -w \
	-X github.com/ksync/ksync/pkg/ksync.GitCommit=${GIT_COMMIT} \
	-X github.com/ksync/ksync/pkg/ksync.BuildDate=${BUILD_DATE} \
	-X github.com/ksync/ksync/pkg/ksync.VersionString=${BINARY_VERSION} \
	-X github.com/ksync/ksync/pkg/ksync.GoVersion=${GO_VERSION} \
	-X github.com/ksync/ksync/pkg/ksync.GitTag=${CIRCLE_TAG} \
	-X github.com/ksync/ksync/pkg/radar.GitCommit=${GIT_COMMIT} \
	-X github.com/ksync/ksync/pkg/radar.BuildDate=${BUILD_DATE} \
	-X github.com/ksync/ksync/pkg/radar.VersionString=${BINARY_VERSION} \
	-X github.com/ksync/ksync/pkg/radar.GoVersion=${GO_VERSION} \
	-X github.com/ksync/ksync/pkg/radar.GitTag=${CIRCLE_TAG} \
	${LDFLAGS}

GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

SHELL=/bin/bash -o pipefail

GOOS=linux
GOARCH=amd64

.PHONY: build-cmd
build-cmd:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) \
		-tags '$(TAGS)' \
		-ldflags '$(LDFLAGS)' \
		github.com/ksync/ksync/cmd/ksync

HAS_GOX := $(shell command -v gox)

.PHONY: build-ci
build-ci:
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif
	${GOBIN}/gox --ldflags "${LDFLAGS}" \
		--parallel=10 \
		--output="bin/{{ .Dir }}_{{ .OS }}_{{ .Arch }}" \
		-os="!netbsd !freebsd !openbsd" -arch="amd64" \
		./cmd/...

.PHONY: docker-binary
docker-binary: BINDIR = $(CURDIR)/bin
docker-binary: GOFLAGS += -installsuffix cgo
docker-binary: docker-binary-radar

docker-binary-%:
	time GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		$(GO) build $(GOFLAGS) \
			-tags '$(TAGS)' \
			-ldflags '$(LDFLAGS)' \
			-o $(BINDIR)/$*_$(GOOS)_$(GOARCH) \
			github.com/ksync/ksync/cmd/$*

.PHONY: docker-build
docker-build:
	docker build --rm -t ${IMAGE} -f docker/Dockerfile ./
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

.PHONY: docker-import
docker-import:
	k3d import-images ${IMAGE}

.PHONY: docker-push
docker-push:
	docker push ${IMAGE}

# .PHONY: docker-tag-release
# docker-tag-release:
# 	docker pull ${IMAGE}
# ifdef CIRCLE_TAG
# 	docker tag ${IMAGE} ${IMAGE_TAG}
# 	docker push ${IMAGE_TAG}
# endif

.PHONY: ci-test
ci-test:
	go test -v -count=1 ./...

.PHONY: test
test:
	kubectl apply -f testdata/k8s/config/testing.yaml
	go test -v ./...

.PHONY: cluster-setup
cluster-setup:
ifndef HAS_K3D
	wget -q -O - https://raw.githubusercontent.com/rancher/k3d/master/install.sh | bash
endif
	k3d create --wait 30

## ------- UNUSED

# .PHONY: all
# all: build docker-binary docker-build

# .PHONY: push
# push: docker-push

# .PHONY: build
# build: build-proto build-cmd

# HAS_PROTOC := $(shell command -v protoc-gen-go)

# .PHONY: build-proto
# build-proto:
# ifndef HAS_PROTOC
# 	go get -u github.com/golang/protobuf/protoc-gen-go
# endif
# 	${GOBIN}/protoc proto/*.proto --go_out=plugins=grpc:pkg

# .PHONY: watch
# watch:
# 	ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) build && $(CMD)"

# .PHONY: update-radar
# update-radar: docker-binary-radar docker-build docker-push update-radar-image

# .PHONY: update-radar-image
# update-radar-image:
# 	bin/ksync* --log-level=debug --image=${IMAGE} init --upgrade --local=false

# HAS_LINT := $(shell command -v gometalinter)

# .PHONY: lint
# lint: install-linter most-lint

# .PHONY: install-linter
# install-linter:
# ifndef HAS_LINT
# 	go get -u github.com/alecthomas/gometalinter
# 	gometalinter --install
# endif

# .PHONY: most-lint
# most-lint:
# 	gometalinter ./... \
# 		--vendor \
# 		--skip "testdata" \
# 		--exclude "[a-zA-Z]*_test.go" \
# 		--exclude "[a-zA-Z]*.pb.go" \
# 		--tests \
# 		--sort=severity \
# 		--aggregate \
# 		--deadline=500s

# HAS_STERN := $(shell command -v stern)

# .PHONY: radar-logs
# radar-logs:
# ifndef HAS_STERN
# 	@printf "Install stern: https://github.com/wercker/stern/releases"; exit 1
# endif
# 	stern --namespace=kube-system --selector=app=radar
