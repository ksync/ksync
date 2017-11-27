
IMAGE_BASE      := gcr.io/elated-embassy-152022/ksync/ksync
MUTABLE_VERSION := canary
DOCKER_VERSION  := git-$(shell git rev-parse --short HEAD)
IMAGE           := ${IMAGE_BASE}:${DOCKER_VERSION}
MUTABLE_IMAGE   := ${IMAGE_BASE}:${MUTABLE_VERSION}

#CMD       ?= bin/ksync --log-level=debug init --upgrade && stern --namespace=kube-system --selector=app=radar
#CMD       ?= bin/ksync --log-level=debug init --upgrade && bin/ksync --log-level=debug watch
CMD ?= bin/ksync --log-level=debug watch

GO        ?= go
TAGS      :=
LDFLAGS   :=
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

SHELL=/bin/bash

GOOS=linux
GOARCH=amd64

.PHONY: all
all: build docker-binary docker-build

.PHONY: push
push: docker-push

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
	# ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) all push && $(CMD) && stern --namespace=kube-system --selector=app=radar"
	ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) build && $(CMD)"

HAS_DEP := $(shell command -v dep)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure

.PHONY: docker-binary
docker-binary: BINDIR = $(CURDIR)/docker/bin
docker-binary: GOFLAGS += -installsuffix cgo
docker-binary: docker-binary-radar

docker-binary-%:
	time GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		$(GO) build $(GOFLAGS) \
			-tags '$(TAGS)' \
			-ldflags '$(LDFLAGS)' \
			-o $(BINDIR)/$*_$(GOOS)_$(GOARCH) \
			github.com/vapor-ware/ksync/cmd/$*

.PHONY: docker-build
docker-build:
	docker build --rm -t ${IMAGE} -f docker/Dockerfile ./
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

.PHONY: docker-push
docker-push:
	gcloud docker -- push ${IMAGE}
	gcloud docker -- push ${MUTABLE_IMAGE}

.PHONY: test
test:
	kubectl apply -f testdata/k8s/config/testing.yaml
	go test -v ./...

HAS_LINT := $(shell command -v gometalinter)

.PHONY: lint
lint:
ifndef HAS_LINT
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
endif
	gometalinter ./... \
		--vendor \
		--skip "testdata" \
		--exclude "[a-zA-Z]*_test.go" \
		--disable=megacheck \
		--deadline=240s
	gometalinter ./...\
		--vendor \
		--skip "testdata" \
		--exclude "[a-zA-Z]*_test.go" \
		--disable-all \
		--enable=megacheck \
		--deadline=240s

HAS_STERN := $(shell command -v stern)

.PHONY: radar-logs
radar-logs:
ifndef HAS_STERN
	@printf "Install stern: https://github.com/wercker/stern/releases"; exit 1
endif
	stern --namespace=kube-system --selector=app=radar
