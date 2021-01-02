SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags)
GO=$(shell which go)
GOFMT=$(GO) fmt
GOBUILD=$(GO) build -ldflags "-X main.version=$(VER)"

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

format:
	$(GOFMT) ./...

build/armv7:
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/entropyd-$(VER:v%=%)-armv7 cmd/entropyd/main.go

build/arm64:
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/entropyd-$(VER:v%=%)-arm64 cmd/entropyd/main.go

build/386:
	export GOOS=linux
	export GOARCH=386
	$(GOBUILD) -o bin/entropyd-$(VER:v%=%)-386 cmd/entropyd/main.go

build/amd64:
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/entropyd-$(VER:v%=%)-amd64 cmd/entropyd/main.go

build: build/armv7 build/arm64 build/386 build/amd64

clean:
	@rm -rf bin

all: dir format build
