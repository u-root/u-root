#!/usr/bin/env bash

# because things are never simple.
# See https://github.com/codecov/example-go#caveat-multiple-files

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

# build systemboot using u-root
go get -u github.com/u-root/u-root
"${GOPATH}/bin/u-root" -build=bb core localboot netboot uinit
