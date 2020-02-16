#!/bin/sh

set -eox

# To easily cross-compile binaries
GO111MODULE=off go get github.com/mitchellh/gox

VERSION=${DRONE_TAG:-latest}
GIT_COMMIT=$(git rev-list -1 HEAD || echo 'dirrrty')

CURRENT_ARCH="$(go env GOOS)/$(go env GOARCH)"

BUILD_ARCHS=${BUILD_ARCHS:-$CURRENT_ARCH}

mkdir -p output
CGO_ENABLED=0 gox -output="output/trieres_{{.OS}}_{{.Arch}}" \
  -osarch="${BUILD_ARCHS}" \
  -ldflags "-s -w  -X github.com/jakolehm/trieres/main.Version=${VERSION}" \
  github.com/jakolehm/trieres/
