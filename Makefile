APP_VERSION = $(shell git describe --abbrev=0 --tags)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date -u "+%Y%m%d-%H%M")
VERSION_PKG = github.com/InjectiveLabs/dexterm/version

all:

install: export GO111MODULE=on
install: export GOPROXY=direct
install: export VERSION_FLAGS="-X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) -X $(VERSION_PKG).BuildDate=$(BUILD_DATE)"
install:
	go install \
		-ldflags $(VERSION_FLAGS) \
		github.com/InjectiveLabs/dexterm

.PHONY: install image push gen lint test mock cover

lint: export GO111MODULE=on
lint: export GOPROXY=direct
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --no-config --issues-exit-code=1 --enable-all --disable=gocyclo --disable=nakedret --disable=gochecknoglobals --tests=false --disable=goimports --disable=wsl

build-release: export DOCKER_BUILDKIT=1
build-release: export VERSION_FLAGS="-X $(VERSION_PKG).AppVersion=$(APP_VERSION) -X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) -X $(VERSION_PKG).BuildDate=$(BUILD_DATE)"
build-release:
	docker build \
		--build-arg LDFLAGS=$(VERSION_FLAGS) \
		--build-arg PKG=github.com/InjectiveLabs/dexterm \
		--ssh=default -t dexterm-release -f Dockerfile.release .

prepare-release:
	mkdir -p dist/dexterm_linux_amd64/
	mkdir -p dist/dexterm_darwin_amd64/
	mkdir -p dist/dexterm_windows_amd64/
	#
	docker create --name tmp_dexterm dexterm-release bash
	#
	docker cp tmp_dexterm:/root/go/bin/dexterm-linux-amd64 dist/dexterm_linux_amd64/dexterm
	docker cp tmp_dexterm:/root/go/bin/dexterm-darwin-amd64 dist/dexterm_darwin_amd64/dexterm
	docker cp tmp_dexterm:/root/go/bin/dexterm-windows-amd64 dist/dexterm_windows_amd64/dexterm.exe
	#
	docker rm tmp_dexterm

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser --rm-dist
