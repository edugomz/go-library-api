MODULE  := library-api
IMAGE   ?= library-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build test docker-build version

build:
	go build -ldflags "-X $(MODULE)/internal/version.Version=$(VERSION)" -o server ./cmd/api

test:
	go test ./...

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

version:
	@echo $(VERSION)
