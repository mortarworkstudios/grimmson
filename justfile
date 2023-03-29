
set-variables:
##!/usr/bin/env bash
VERSION := `cat VERSION`
BUILD_TIME := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
COMMIT := `git rev-parse --short HEAD`

# build the binary file
build:set-variables
    go build -o ./grimmson -ldflags "-X=main.version={{VERSION}} -X=main.buildTime={{BUILD_TIME}}  -X=main.gitCommit={{COMMIT}}"