# Makefile for stackgen project

APP_NAME=server
DOCKER_IMAGE=stackgen:latest
GO_FILES=$(shell find . -type f -name '*.go')

.PHONY: all build docker-build run fmt test clean

all: build

build:
	go build -o $(APP_NAME) ./cmd/server/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

run:
	./$(APP_NAME)

fmt:
	gofmt -w $(GO_FILES)

test:
	go test ./...

clean:
	rm -f $(APP_NAME)
