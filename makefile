all: 
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(git describe --always --long)" -o slab