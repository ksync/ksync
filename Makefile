
IMAGE_BASE      := gcr.io/elated-embassy-152022/ksync/
MUTABLE_VERSION := canary
DOCKER_VERSION  := git-$(shell git rev-parse --short HEAD)

CMD       ?= bin/ksync --log-level=debug init --upgrade

GO        ?= go
TAGS      :=
LDFLAGS   :=
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

SHELL=/bin/bash

.PHONY: all
all: build docker-build-radar docker-build-mirror

.PHONY: push
push: docker-push-radar docker-push-mirror

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
	ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) all push && $(CMD) && stern --namespace=kube-system --selector=app=radar"

HAS_DEP := $(shell command -v dep)

.PHONY: bootstrapmake
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure

.PHONY: docker-binary
docker-binary: BINDIR = $(CURDIR)/radar/bin
docker-binary: GOFLAGS += -a -installsuffix cgo
docker-binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) \
		-tags '$(TAGS)' \
		-ldflags '$(LDFLAGS)' \
		github.com/vapor-ware/ksync/cmd/radar

docker-build-%: IMAGE = ${IMAGE_BASE}$*:${DOCKER_VERSION}
docker-build-%: MUTABLE_IMAGE = ${IMAGE_BASE}$*:${MUTABLE_VERSION}
docker-build-%:
	docker build --rm -t ${IMAGE} $*
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

docker-push-%: IMAGE = ${IMAGE_BASE}$*:${DOCKER_VERSION}
docker-push-%: MUTABLE_IMAGE = ${IMAGE_BASE}$*:${MUTABLE_VERSION}
docker-push-%:
	gcloud docker -- push ${IMAGE}
	gcloud docker -- push ${MUTABLE_IMAGE}

docker-build-radar: docker-binary
