MODULE  := library-api
IMAGE   ?= library-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build test migrate docker-build version

build:
	go build -ldflags "-X $(MODULE)/internal/version.Version=$(VERSION)" -o server ./cmd/api

test:
	go test ./...

migrate:
	go run ./cmd/api --migrate-only

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

version:
	@echo $(VERSION)
