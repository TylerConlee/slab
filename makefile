VERSION = $(git describe --always --long --dirty)
GOOS = linux
GOARCH = amd64
OUTPUT = slab
all: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X main.Version=$(VERSION)" -o $(OUTPUT)
deploy:
	scp slab ubuntu@35.160.9.184:slab
	ssh ubuntu@35.160.9.184
	supervisord restart
	exit