BIN := "./bin"
DOCKER_IMG="calendar:develop"
PROTOPATH:=$(GOPATH)/src/github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/api

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	mkdir -p pkg/calendar
	protoc --proto_path=$(GOPATH)/src --go_out=. --go-grpc_out=. $(PROTOPATH)/EventService.proto

build:
	go build -v -o $(BIN)/calendar -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)/scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN)/sender -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN)/calendar -config ./configs/calendar/config.toml
	$(BIN)/scheduler -config ./configs/sheduler/config.toml
	$(BIN)/sender -config ./configs/sender/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint generate
