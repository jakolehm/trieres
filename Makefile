.PHONY: test clean
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover


BUILD_ARCHS ?= $(shell go env GOOS)/$(shell go env GOARCH)
# User drone tag as the build verion if given in env
VERSION ?= $(or ${DRONE_TAG},latest)

build:
	mkdir -p output
	CGO_ENABLED=0 gox -output="output/trieres_{{.OS}}_{{.Arch}}" \
		-osarch="${BUILD_ARCHS}" \
  		-ldflags "-s -w -X main.Version=${VERSION}" \
  		github.com/jakolehm/trieres/

clean:
	rm -rf output/*