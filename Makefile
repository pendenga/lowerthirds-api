.PHONY: test
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


test:
	go test ./...
