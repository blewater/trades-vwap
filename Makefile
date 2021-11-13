DOCKER := $(shell which docker)
VERSION = $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  = $(shell git log -1 --format='%H')

build:
	go build -o build/vwap ./cmd

# escape analysis (stack or heap allocations)
escape:
	go build -gcflags '-m -m' -o build/vwap ./cmd

run:
	go run main.go -d

clean:
	rm ./build/*

lint:
	golangci-lint run ./...

# a stricter go formatter
# go get mvdan.cc/gofumpt
fmt:
	gofumpt -w **/*.go

# deterministic imports orderer and formatter
# go get github.com/daixiang0/gci
imp:
	gci -w **/*.go

install:
	go install ./cmd

build-docker:
	$(MAKE) -C docker/ all

test:
	go test ./...

github-ci:
	$(MAKE) test

.PHONY: build test build-docker license
