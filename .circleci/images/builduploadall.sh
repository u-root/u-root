#!/bin/bash
set -ex

VERSION=v4.4.0

VC=$(git diff --exit-code builduploadall.sh | grep --no-ignore-case -P "VERSION=+" | wc -l)
echo $vc
if [[ $VC != "3" ]]; then
    echo "Increment version before a push"
    exit 1
fi

for GOARCH in amd64 arm arm64; do
  (
    cd test-image-$GOARCH
    docker build . -t uroottest/test-image-$GOARCH:$VERSION
    docker push uroottest/test-image-$GOARCH:$VERSION
  )
done
