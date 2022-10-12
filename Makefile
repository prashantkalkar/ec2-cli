.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

test: vet
	go test -v ./...
.PHONY:test

build: test
	go build .
.PHONY:build

dist:
	GOOS=darwin GOARCH=amd64 go build -o ./darwin/amd64/ec2-cli .
	GOOS=linux GOARCH=amd64 go build -o ./linux/amd64/ec2-cli .
.PHONY:dist