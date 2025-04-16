run: build
	@./bin/main

build:
	@go generate ./...
	@go build -o bin/main main.go

