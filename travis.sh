#!/bin/bash
if [ -z "${GOPATH}" ]; then
        export GOPATH=/home/travis/gopath
fi
set -e

echo "Check vendored dependencies"
 (dep status)

echo "Build u-root"
 (go build u-root.go)

echo "-----------------------> Initial bb test"
 (./u-root -build=bb)
mv /tmp/initramfs.linux_amd64.cpio /tmp/i2

# Test for reproducible initramfs in busybox mode
echo "-----------------------> Second bb test"
 (./u-root -build=bb)

echo "-----------------------> cmp bb test output (test reproducibility)"
cmp /tmp/initramfs.linux_amd64.cpio /tmp/i2
 which go

# Test all architectures we care about. At some point we may just
# grow the build matrix.
echo "-----------------------> ARM64 test build"
 (GOARCH=arm64 ./u-root -build=bb)
echo "-----------------------> ppc64le test build"
 (GOARCH=ppc64le ./u-root -build=bb)

echo "-----------------------> First ramfs test"
 (./u-root -build=source --tmpdir=/tmp/u-root)

echo "-----------------------> build all tools"
 (cd cmds && CGO_ENABLED=0 go build -a -installsuffix uroot -ldflags '-s' ./...)

echo "-----------------------> What got built? ls -l cmds/*"
 ls -l cmds/*

echo "-----------------------> go test"
 (cd cmds && CGO_ENABLED=0 go test -a -installsuffix uroot -ldflags '-s' ./...)

echo "-----------------------> test -cover"
 (cd cmds && CGO_ENABLED=0 go test -cover ./...)
 (cd pkg && CGO_ENABLED=0 go test -cover ./...)

echo "-----------------------> go vet"
  (go tool vet cmds pkg)

# is it go-gettable?
echo "-----------------------> test go-gettable"
 (go get github.com/u-root/u-root)

echo "Did it blend"
