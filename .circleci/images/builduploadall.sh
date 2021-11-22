#!/bin/bash
set -ex

VERSION=v4.4.0

for GOARCH in amd64 arm arm64; do
  (
    cd test-image-$GOARCH
    docker build . -t uroottest/test-image-$GOARCH:$VERSION
    docker push uroottest/test-image-$GOARCH:$VERSION
  )
done
