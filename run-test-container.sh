#!/bin/bash

# Runs the container used for tests on CircleCI locally, with your u-root
# bind-mounted into the test path.
REPO="github.com/u-root/u-root"
docker run                                                             \
  -it --user=$(id -u)                                                  \
  --mount type=bind,source="$GOPATH/src/$REPO",target="/go/src/$REPO"  \
  -w "/go/src/$REPO" \
  uroottest/test-image-amd64:v3.2.6
