CUR_DIR = $(CURDIR)
all: check-style test

## Runs golangci-lint
.PHONY: check-style
check-style:
	golangci-lint run ./...

## Builds project
.PHONY: build
build:
	go build .

## Runs tests
.PHONY: test
test:
	go test ./...

## Runs benchmarks
.PHONY: bench
bench:
	go test -bench=.