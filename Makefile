GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

format:
	$(GOFMT) ./...

build:
	$(GOBUILD) -o bin/entropyd cmd/entropyd/main.go

clean:
	@rm -rf bin

all: dir format build
