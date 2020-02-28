#!/bin/sh
# you should have golangci-lint installed:  https://github.com/golangci/golangci-lint
# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.23.6
golangci-lint run --no-config --issues-exit-code=1 --enable-all --disable=gocyclo --disable=nakedret --disable=gochecknoglobals --tests=false --disable=goimports --disable=wsl