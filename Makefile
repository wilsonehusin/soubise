GOROOT?=$(shell go env GOROOT)
GIT_TAG?=$(shell git describe --always --dirty --tags)
GIT_SHA?=$(shell git rev-parse --verify HEAD)
DEFAULT_SERVER?=https://test.soubise.org
CONTAINER_REGISTRY?=ghcr.io/wilsonehusin
GOX_FLAGS?=-osarch="darwin/amd64 linux/amd64 linux/arm darwin/arm64 windows/amd64"
OUT_DIR?=_output

GOTARGET=github.com/wilsonehusin/soubise
GOVERSION=$(shell go env GOVERSION)

REPO_ROOT=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SOUBISE_IMAGES=$(shell docker image ls --all --format "{{.Repository}}:{{.Tag}}" | grep soubise)

VERSION_FLAG=-X=$(GOTARGET)/internal/buildinfo.Version=$(GIT_TAG)
GIT_SHA_FLAG=-X=$(GOTARGET)/internal/buildinfo.GitSHA=$(GIT_SHA)
SERVER_FLAG=-X=$(GOTARGET)/internal/buildinfo.Server=$(DEFAULT_SERVER)
COMPILER_FLAG=-X=$(GOTARGET)/internal/buildinfo.Go=$(GOVERSION)

BUILD_FLAGS=-ldflags '$(VERSION_FLAG) $(GIT_SHA_FLAG) $(COMPILER_FLAG)'
TEST_BUILD_FLAGS=-ldflags '$(VERSION_FLAG) $(GIT_SHA_FLAG) $(COMPILER_FLAG) $(SERVER_FLAG)'
DOCKER_TAGS=--tag $(CONTAINER_REGISTRY)/soubise:latest --tag $(CONTAINER_REGISTRY)/soubise:$(GIT_TAG)

.PHONY: all
all: build container

.PHONY: testbuild
testbuild: GIT_TAG:=zz_test-$(GIT_TAG)
testbuild: BUILD_FLAGS:=$(TEST_BUILD_FLAGS)
testbuild: build

.PHONY: testmultibuild
testmultibuild: GIT_TAG:=zz_test-$(GIT_TAG)
testmultibuild: BUILD_FLAGS:=-ldflags '$(VERSION_FLAG) $(GIT_SHA_FLAG) $(COMPILER_FLAG) $(SERVER_FLAG)'
testmultibuild: multibuild

.PHONY: testcontainer
testcontainer: GIT_TAG:=zz_test-$(GIT_TAG)
testcontainer: testbuild
	docker build $(DOCKER_TAGS) .

.PHONY: build
build:
	go build -o $(OUT_DIR)/ $(BUILD_FLAGS) .

.PHONY: multibuild
multibuild:
	gox -output="$(OUT_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" $(GOX_FLAGS) $(BUILD_FLAGS)

.PHONY: container
container: build
	docker build $(DOCKER_TAGS) .

.PHONY: release
release: container
	docker push $(DOCKER_TAGS)

.PHONY: clean
clean:
	rm -rf $(REPO_ROOT)/$(OUT_DIR)
	@docker image rm $(SOUBISE_IMAGES) 2>/dev/null || true
