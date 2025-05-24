.PHONY: test
ROOT_DIR = $(shell pwd)
GOPATH:=$(shell go env GOPATH)

#################################################################################
# RUN COMMANDS
#################################################################################
run:
	go mod vendor
	ENV_FILES_DIR=./build/secrets go run ./cmd/lowerthirds-api/main.go
	rm -rf vendor

test:
	go test ./...
