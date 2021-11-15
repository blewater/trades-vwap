DOCKER := $(shell which docker)
VERSION = $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  = $(shell git log -1 --format='%H')

build:
	go build -o build/vwap main.go

check-race:
	go build -race -o build/vwap main.go

# escape analysis (stack or heap allocations)
escape:
	go build -gcflags '-m -m' -o build/vwap ./cmd

run:
	go run main.go -d -w 4 -u wss://ws-feed-public.sandbox.exchange.coinbase.com

run-prod: build
	build/vwap -u wss://ws-feed.exchange.coinbase.com -w 5

clean:
	rm ./build/*
	rm workflow.test
	rm profile_cpu*

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
	$(DOCKER) build . --tag blewater/vwap

run-docker:
	$(DOCKER) run -it blewater/vwap

test:
	go test ./...

bench:
	go test ./workflow -run=xxx -bench=. -cpuprofile profile_cpu.out
	go tool pprof -svg profile_cpu.out > profile_cpu.svg

github-ci:
	$(MAKE) test

.PHONY: clean bench build check-race run run-prod lint imp fmt test github-ci build-docker run-docker
