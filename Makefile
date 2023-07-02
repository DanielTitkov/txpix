.PHONY: run air all build build-linux build-windows build-darwin

DIR = ./tmp

drop: 
	cd $(DIR) && rm *

run:
	go run ./cmd/generate/main.go -dir $(DIR)

air:
	go run ./cmd/fetch/main.go -dir $(DIR)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/txpix-linux ./cmd/generate/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/txpix-windows.exe ./cmd/generate/main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/txpix-darwin ./cmd/generate/main.go

build: build-linux build-windows build-darwin

all: air run
