.PHONY: build dev test lint clean

BIN := bin/server

build:
	go build -o $(BIN) ./cmd/server

dev:
	go run ./cmd/server

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ tmp/
