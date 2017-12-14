VERSION = $(git describe --always --long --dirty)
GOOS = linux
GOARCH = amd64
OUTPUT = slab
all: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X main.Version=$(VERSION)" -o $(OUTPUT)