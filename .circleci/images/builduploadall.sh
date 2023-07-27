#!/bin/bash
set -ex

VERSION=v4.5.0

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

# Tamago has slightly different requirements; until we are sure why,
# do a slightly custom build

(
	cd test-image-tamago
	docker build . --build-arg UID=1000 --build-arg GID=1000 -t uroottest/test-image-tamago:$VERSION
	docker push uroottest/test-image-tamago:$VERSION
)
