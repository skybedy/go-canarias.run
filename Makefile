SHELL := /bin/bash

.PHONY: run test build

run:
	GOCACHE=/tmp/go-build go run .

test:
	GOCACHE=/tmp/go-build go test ./...

build:
	GOCACHE=/tmp/go-build go build ./...
