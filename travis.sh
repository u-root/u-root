#!/bin/bash
if [ -z "${GOPATH}" ]; then
        export GOPATH=/home/travis/gopath
fi
set -e
echo "-----------------------> Initial bb test"
 (cd bb && go build . && ./bb)
mv /tmp/initramfs.linux_amd64.cpio /tmp/i2
# echo Test for reproducible initramfds in busybox mode
echo "-----------------------> echo Second bb test"
 (cd bb && go build . && ./bb)

# Test all architectures we care about. At some point we may just
# grow the build matrix.
echo "-----------------------> ARM64 test build"
(cd bb && go build . && GOARCH=arm64 ./bb)
echo "-----------------------> ppc64le test build"
(cd bb && go build . && GOARCH=ppc64le ./bb)

echo "-----------------------> cmp bb test output"
cmp /tmp/initramfs.linux_amd64.cpio /tmp/i2
 which go

echo "-----------------------> First ramfs test"
 (go run u-root.go -build=source -format=cpio -tmpdir=/tmp/u-root)

echo "-----------------------> build all tools"
 (cd cmds && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' ./...)

echo "-----------------------> What got built? ls -l cmds/*"
 ls -l cmds/*

echo "-----------------------> go test"
 (cd cmds && CGO_ENABLED=0 go test -a -installsuffix cgo -ldflags '-s' ./...)

echo "-----------------------> test -cover" 
 (cd cmds && CGO_ENABLED=0 go test -cover ./...)
 (cd pkg && CGO_ENABLED=0 go test -cover ./...)

echo "-----------------------> go vet"
 go tool vet cmds uroot pkg


# is it go-gettable?
echo "-----------------------> test go-gettable"
 go get github.com/u-root/u-root

echo "Did it blend"
