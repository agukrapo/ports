.DEFAULT_GOAL := all

NAME := $(shell basename $(CURDIR))

all: build test format lint

clean:
	@echo "Cleaning ${NAME}..."
	@go clean -i ./...
	@rm -rf bin

build: clean
	@echo "Building ${NAME}..."
	@go mod download
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpc/upload.proto
	@go build -o ./bin/${NAME} ./cmd

test:
	@echo "Testing ${NAME}..."
	@go test ./... -cover -race -shuffle=on

format:
	@echo "Formatting ${NAME}..."
	@go mod tidy
	@gofumpt -l -w . #go install mvdan.cc/gofumpt@latest

lint:
	@echo "Linting ${NAME}..."
	@go vet ./...
	@golangci-lint run #https://golangci-lint.run/usage/install/
