.PHONY: build clean default

build: clean
	@go build -o bin/memgrdpeek ./main.go

clean:
	@rm -rf ./bin/*

default: build
