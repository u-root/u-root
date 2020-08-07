#!/usr/bin/env bash

set -e

GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint
$GOPATH/bin/golangci-lint run
