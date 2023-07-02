.PHONY: run air all

DIR = ./tmp

drop: 
	cd $(DIR) && rm *

run:
	go run ./cmd/generate/main.go -dir $(DIR)

air:
	go run ./cmd/fetch/main.go -dir $(DIR)

all: air run
