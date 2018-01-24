
BINARY_VERSION  ?= "corrupted-version"
GIT_COMMIT      ?= $(shell git rev-parse --short HEAD)

DATE            := $(shell (which gdate > /dev/null && echo "gdate") || echo "date")
BUILD_DATE      := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
GO_VERSION      := $(shell go version | awk '{ print $$3 }')

IMAGE_BASE      := vaporio/ksync
MUTABLE_VERSION := canary
DOCKER_VERSION  := git-${GIT_COMMIT}
export IMAGE    := ${IMAGE_BASE}:${DOCKER_VERSION}
MUTABLE_IMAGE   := ${IMAGE_BASE}:${MUTABLE_VERSION}
ifdef CIRCLE_TAG
IMAGE_TAG       := ${IMAGE_BASE}:${CIRCLE_TAG}
endif

CMD       ?= bin/ksync --log-level=debug watch

GO        ?= go
TAGS      :=
LDFLAGS   := -w \
	-X github.com/vapor-ware/ksync/pkg/ksync.GitCommit=${GIT_COMMIT} \
	-X github.com/vapor-ware/ksync/pkg/ksync.BuildDate=${BUILD_DATE} \
	-X github.com/vapor-ware/ksync/pkg/ksync.VersionString=${BINARY_VERSION} \
	-X github.com/vapor-ware/ksync/pkg/ksync.GoVersion=${GO_VERSION} \
	-X github.com/vapor-ware/ksync/pkg/ksync.GitTag=${CIRCLE_TAG} \
	-X github.com/vapor-ware/ksync/pkg/radar.GitCommit=${GIT_COMMIT} \
	-X github.com/vapor-ware/ksync/pkg/radar.BuildDate=${BUILD_DATE} \
	-X github.com/vapor-ware/ksync/pkg/radar.VersionString=${BINARY_VERSION} \
	-X github.com/vapor-ware/ksync/pkg/radar.GoVersion=${GO_VERSION} \
	-X github.com/vapor-ware/ksync/pkg/radar.GitTag=${CIRCLE_TAG} \
	${LDFLAGS}

GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

SHELL=/bin/bash -o pipefail

GOOS=linux
GOARCH=amd64
PATH+=:/home/circleci/google-cloud-sdk/bin

.PHONY: all
all: build docker-binary docker-build

.PHONY: push
push: docker-push

.PHONY: build
build: build-proto build-cmd

.PHONY: build-ci
build-ci:
	gox --ldflags "${LDFLAGS}" --parallel=10 --output="bin/{{ .Dir }}_{{ .OS }}_{{ .Arch }}" ${OPTS} ./cmd/...

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
	ag -l --ignore "pkg/proto" | entr -dr /bin/sh -c "$(MAKE) build && $(CMD)"

HAS_DEP := $(shell command -v dep)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure

.PHONY: update-radar
update-radar: docker-binary-radar docker-build docker-push update-radar-image

.PHONY: update-radar-image
update-radar-image:
	bin/ksync* --log-level=debug --image=${IMAGE} init --upgrade --skip-checks

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
			github.com/vapor-ware/ksync/cmd/$*

.PHONY: docker-build
docker-build:
	docker build --rm -t ${IMAGE} -f docker/Dockerfile ./
	docker tag ${IMAGE} ${MUTABLE_IMAGE}
ifdef CIRCLE_TAG
	docker tag ${IMAGE} ${IMAGE_TAG}
endif

.PHONY: docker-push
docker-push:
	docker push ${IMAGE}

.PHONY: docker-tag-release
docker-tag-release:
	docker pull ${IMAGE}
ifdef CIRCLE_TAG
	docker tag ${IMAGE} ${IMAGE_TAG}
	docker push ${IMAGE_TAG}
endif

.PHONY: test
test:
	kubectl apply -f testdata/k8s/config/testing.yaml
	go test -v ./...

.PHONY: ci-test
ci-test:
	go test -v ./... 2>&1 | tee /tmp/${TEST_DIRECTORY}/test.out
	cat /tmp/${TEST_DIRECTORY}/test.out \
		| go-junit-report \
		> /tmp/${TEST_DIRECTORY}/report.xml

HAS_LINT := $(shell command -v gometalinter)

.PHONY: lint
lint: install-linter most-lint megacheck

.PHONY: install-linter
install-linter:
ifndef HAS_LINT
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
endif

.PHONY: most-lint
most-lint:
	gometalinter ./... \
		--vendor \
		--skip "testdata" \
		--exclude "[a-zA-Z]*_test.go" \
		--disable=megacheck \
		--tests \
		--sort=severity \
		--aggregate \
		--deadline=500s

.PHONY: megacheck
megacheck:
	gometalinter ./... \
		--vendor \
		--skip "testdata" \
		--exclude "[a-zA-Z]*_test.go" \
		--disable-all \
		--tests \
		--sort=severity \
		--aggregate \
		--enable=megacheck \
		--deadline=240s

HAS_STERN := $(shell command -v stern)

.PHONY: radar-logs
radar-logs:
ifndef HAS_STERN
	@printf "Install stern: https://github.com/wercker/stern/releases"; exit 1
endif
	stern --namespace=kube-system --selector=app=radar
