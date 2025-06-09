.PHONY: test test-coverage
ROOT_DIR = $(shell pwd)
GOPATH:=$(shell go env GOPATH)

#################################################################################
# RUN COMMANDS
#################################################################################
run:
	go mod vendor
	ENV_FILES_DIR=./build/secrets
	go run ./cmd/lowerthirds-api/main.go
	rm -rf vendor

run-cgi:
	go mod vendor
	ENV_FILES_DIR=./build/secrets
	go build -o ./build/public/lowerthirds.fcgi ./cgi/lowerthirds-api/main.go
	rm -rf vendor

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean up coverage files
clean:
	rm -f coverage.out coverage.html
