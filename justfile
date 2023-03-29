
#set-variables:
##!/usr/bin/env bash
#VERSION := $(shell cat VERSION)
#echo $VERSION
##BUILD_TIME := $(date -u +"%Y-%m-%dT%H:%M:%SZ")
##COMMIT := git rev-parse --short HEAD
##LDFLAGS := -ldflags "-X=main.version=$(VERSION) -X=main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(COMMIT)"

# build the binary file
build:
    go build $(LDFLAGS) -o ./grimmson