GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

export GOARCH=amd64
export GOOS=linux

dir:
	@if [ ! -d $(CURDIR)/bin ] ; then mkdir -p $(CURDIR)/bin ; fi

format:
	$(GOFMT) main.go

build:
	$(GOBUILD) -o $(CURDIR)/bin/entropyd main.go

clean:
	@rm -rf $(CURDIR)/bin

all: dir format build
