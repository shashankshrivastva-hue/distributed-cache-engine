.PHONY: build run test clean

build:
	go build -o bin/cache-server cmd/cache-server/main.go

test:
	go test -v ./...

run: build
	./bin/cache-server
