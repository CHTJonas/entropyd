SHELL := bash
.ONESHELL:

GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

format:
	$(GOFMT) ./...

build/armv7:
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/entropyd-armv7 cmd/entropyd/main.go cmd/entropyd/version.go

build/arm64:
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/entropyd-arm64 cmd/entropyd/main.go cmd/entropyd/version.go

build/386:
	export GOOS=linux
	export GOARCH=386
	$(GOBUILD) -o bin/entropyd-386 cmd/entropyd/main.go cmd/entropyd/version.go

build/amd64:
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/entropyd-amd64 cmd/entropyd/main.go cmd/entropyd/version.go

build: build/armv7 build/arm64 build/386 build/amd64

clean:
	@rm -rf bin

all: dir format build
