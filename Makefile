GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

format:
	$(GOFMT) ./...

build/arm64:
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/linux-arm64/entropyd cmd/entropyd/main.go

build/amd64:
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/linux-amd64/entropyd cmd/entropyd/main.go

build: build/arm64 build/amd64

clean:
	@rm -rf bin

all: dir format build
