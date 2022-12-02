SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags --always --dirty)
GO=$(shell which go)
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build -trimpath -mod=readonly -ldflags "-X main.version=$(VER:v%=%) -s -w -buildid="

dir:
	@if [ ! -d bin ]; then mkdir -p bin; fi

mod:
	@$(GOMOD) download

format:
	@$(GOFMT) ./...

build/linux/armv7: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/entropyd-linux-$(VER:v%=%)-armv7 cmd/entropyd/*.go

build/linux/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/entropyd-linux-$(VER:v%=%)-arm64 cmd/entropyd/*.go

build/linux/386: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=386
	$(GOBUILD) -o bin/entropyd-linux-$(VER:v%=%)-386 cmd/entropyd/*.go

build/linux/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/entropyd-linux-$(VER:v%=%)-amd64 cmd/entropyd/*.go

build/linux: build/linux/armv7 build/linux/arm64 build/linux/386 build/linux/amd64

build: build/linux

clean:
	@rm -rf bin

all: format build
