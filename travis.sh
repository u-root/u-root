#!/bin/bash
export GOPATH=/home/travis/gopath
set -e
 (cd bb && go build . && ./bb)
mv /tmp/initramfs.linux_amd64.cpio /tmp/i2
# echo Test for reproducible initramfds in busybox mode
 (cd bb && go build . && ./bb)
cmp /tmp/initramfs.linux_amd64.cpio /tmp/i2
 which go
 (cd scripts && go run ramfs.go -tmpdir=/tmp/u-root -removedir=false)
 GOBIN=/tmp/u-root/ubin GOROOT=/tmp/u-root/go GOPATH=/tmp/u-root CGO_ENABLED=0 /tmp/u-root/go/bin/go build  github.com/u-root/u-root/cmds/ip
 (cd cmds && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' ./...)
 ls -l cmds/*
 (cd cmds && CGO_ENABLED=0 go test -a -installsuffix cgo -ldflags '-s' ./...)
 (cd cmds && CGO_ENABLED=0 go test -cover ./...)
 go tool vet cmds uroot pkg
 go tool vet scripts/ramfs.go
 sudo date
 echo "Did it blend"
