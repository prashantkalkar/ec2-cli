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
	GOOS=darwin GOARCH=amd64 go build -o ./ec2-cli-darwin-amd64 .

	GOOS=linux GOARCH=amd64 go build -o ./linux/amd64/ec2-cli .
	GOOS=linux GOARCH=amd64 go build -o ./ec2-cli-linux-amd64 .

	GOOS=darwin GOARCH=arm64 go build -o ./darwin/arm64/ec2-cli .
	GOOS=darwin GOARCH=arm64 go build -o ./ec2-cli-darwin-arm64 .
.PHONY:dist