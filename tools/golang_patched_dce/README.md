# Golang Docker image with DCE patches

## Overview

This is a Docker image that contains patched Go stdlib, namely:

 * [ad1d54f.diff](ad1d54f.diff) - https://golang.org/cl/210284

## Building

```
$ docker build -t golang-patched-dce .
...
Successfully tagged golang-patched-dce:latest
```

## Usage

```
$ docker run --rm -it -u $(id -u):$(id -g) \
      -v $(go env GOPATH)/src:/go/src \
      -v $(go env GOCACHE):/go/.cache \
      -v $PWD:/out \
    golang-patched-dce sh -c 'go build github.com/u-root/u-root && \
                              ./u-root -build=bb -o /out/initramfs'
```
