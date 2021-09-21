# Circleci Images

Each folder contains the Dockerfile for running integration tests with a
different architecture.


## Build and Run

The following examples assume two variables:

- `$GOARCH`: the architecture
- `$VERSION`: a "vX.Y.Z" string, a version which has not been used yet
  (ideally this follows semantic versioning).

See previously built images at:

    https://hub.docker.com/r/uroottest/test-image-$GOARCH/tags

Note: If `id | grep docker` contains no matches, your user is not a member of
the docker group, so you will need to run each of the following `docker`
commands with `sudo`.

Build and run a new image:

    cd test-image-$GOARCH
    docker build . -t uroottest/test-image-$GOARCH:$VERSION
    docker run --rm -it uroottest/test-image-$GOARCH:$VERSION


## Push

Push:

    # Ping Ryan O'Leary (on slack or via email) for push access.
    docker login
    docker push uroottest/test-image-$GOARCH:$VERSION

Remember to update the image version in `.circleci/config.yml`.

More instructions:

    https://circleci.com/docs/2.0/custom-images/
