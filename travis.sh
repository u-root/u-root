#!/bin/bash
 (cd bb && go build . && ./bb)
 which go
 (cd scripts && go run ramfs.go -d -tmpdir=/tmp/u-root -removedir=false)
 GOBIN=/tmp/u-root/ubin GOROOT=/tmp/u-root/go GOPATH=/tmp/u-root CGO_ENABLED=0 /tmp/u-root/go/bin/go build -x github.com/u-root/u-root/cmds/ip
 (cd cmds && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' ./...)
 ls -l cmds/*
 (cd cmds && CGO_ENABLED=0 go test -a -installsuffix cgo -ldflags '-s' ./...)
 (go test -cover ./...)
 go tool vet cmds uroot netlink scripts
 sudo date
 echo "Did it blend"
