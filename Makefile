
IMAGE_BASE      := gcr.io/elated-embassy-152022/ksync/radar
MUTABLE_VERSION := canary
DOCKER_VERSION  := git-$(shell git rev-parse --short HEAD)
IMAGE           := ${IMAGE_BASE}:${DOCKER_VERSION}
MUTABLE_IMAGE   := ${IMAGE_BASE}:${MUTABLE_VERSION}

CMD       ?= bin/ksync --log-level=debug init --upgrade

GO        ?= go
TAGS      :=
LDFLAGS   :=
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

SHELL=/bin/bash

.PHONY: all
all: build

.PHONY: build
build: build-proto build-cmd

.PHONY: build-cmd
build-cmd:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) \
		-tags '$(TAGS)' \
		-ldflags '$(LDFLAGS)' \
		github.com/vapor-ware/ksync/cmd/ksync

.PHONY: build-proto
build-proto:
	protoc proto/*.proto --go_out=plugins=grpc:pkg

.PHONY: watch
watch:
	ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) build docker-build docker-push && $(CMD) && stern --namespace=kube-system --selector=app=radar"

HAS_DEP := $(shell command -v dep)

.PHONY: bootstrapmake
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure

.PHONY: docker-binary
docker-binary: BINDIR = $(CURDIR)/rootfs/bin
docker-binary: GOFLAGS += -a -installsuffix cgo
docker-binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) \
		-tags '$(TAGS)' \
		-ldflags '$(LDFLAGS)' \
		github.com/vapor-ware/ksync/cmd/radar

.PHONY: docker-build
docker-build: docker-binary
	docker build --rm -t ${IMAGE} rootfs
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

.PHONY: docker-push
docker-push:
	gcloud docker -- push ${IMAGE}
	gcloud docker -- push ${MUTABLE_IMAGE}
