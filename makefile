VERSION = $(shell git rev-parse --short HEAD)
GOOS = linux
GOARCH = amd64
OUTPUT = slab
BINARY_UNIX = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))/release/linux/slab
BINARY_DARWIN = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))/release/darwin/slab

all: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X main.Version=$(VERSION)" -o $(OUTPUT)

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)"  -o $(BINARY_UNIX) -v

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)"  -o $(BINARY_DARWIN) -v