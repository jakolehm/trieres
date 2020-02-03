.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover