# Circleci Images

Each folder contains the Dockerfile for running integration tests with a
different architectures.


## Build and Run

The following examples assume two variables:

- `$GOARCH`: the architecture
- `$VERSION`: a "vX.Y.Z" string, a version which has not been used yet
  (ideally this follows semantic versioning).

See previously built images at:

    https://hub.docker.com/r/uroottest/test-image-$GOARCH/tags

Build and run a new image:

    sudo docker build . -t uroottest/test-image-$GOARCH:$VERSION
    sudo docker run --rm -it uroottest/test-image-$GOARCH:$VERSION


## Push

Push:

    # Ping Ryan O'Leary (on slack or via email) for push access.
    docker push uroottest/test-image-$GOARCH:$VERSION

More instructions:

    https://circleci.com/docs/2.0/custom-images/


## Alternative Method

With this new trick, there's no need to install Docker! Simply push a tag with
the following name:

- `test-image-$GOARCH-$VERSION`

The new Docker image takes about 40 minutes to build and propagate.
