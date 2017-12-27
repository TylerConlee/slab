VERSION = $(shell git rev-parse --short HEAD)
GOOS = linux
GOARCH = amd64
OUTPUT = slab
all: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X main.Version=$(VERSION)" -o $(OUTPUT)
deploy:
	ssh -o StrictHostKeyChecking=no ubuntu@35.160.9.184 supervisorctl stop slab
	scp -o StrictHostKeyChecking=no slab ubuntu@35.160.9.184:slab
	ssh -o StrictHostKeyChecking=no ubuntu@35.160.9.184 supervisorctl start slab
