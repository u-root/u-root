# binctr

[![Build Status](https://travis-ci.org/genuinetools/binctr.svg?branch=master)](https://travis-ci.org/genuinetools/binctr)
[![Go Report Card](https://goreportcard.com/badge/github.com/genuinetools/binctr)](https://goreportcard.com/report/github.com/genuinetools/binctr)
[![GoDoc](https://godoc.org/github.com/genuinetools/binctr?status.svg)](https://godoc.org/github.com/genuinetools/binctr)

Create fully static, including rootfs embedded, binaries that pop you directly
into a container. **Can be run by an unprivileged user.**

Check out the blog post: [blog.jessfraz.com/post/getting-towards-real-sandbox-containers](https://blog.jessfraz.com/post/getting-towards-real-sandbox-containers/).

This is based off a crazy idea from [@crosbymichael](https://github.com/crosbymichael)
who first embedded an image in a binary :D

**HISTORY:** This project used to use a POC fork of libcontainer until [@cyphar](https://github.com/cyphar)
got rootless containers into upstream! Woohoo!
Check out the original thread on the 
[mailing list](https://groups.google.com/a/opencontainers.org/forum/#!topic/dev/yutVaSLcqWI).

**Table of Contents**

<!-- toc -->

  * [Checking out this repo](#checking-out-this-repo)
  * [Building](#building)
  * [Running](#running)
- [Cool things](#cool-things)

<!-- tocstop -->

### Checking out this repo

```console
$ git clone git@github.com:genuinetools/binctr.git
```

### Building

You will need `libapparmor-dev` and `libseccomp-dev`.

Most importantly you need userns in your kernel (`CONFIG_USER_NS=y`)
or else this won't even work.

```console
# building the alpine example
$ make alpine
Static container created at: ./alpine

# building the busybox example
$ make busybox
Static container created at: ./busybox

# building the cl-k8s example
$ make cl-k8s
Static container created at: ./cl-k8s
```

### Running

```console
$ ./alpine
$ ./busybox
$ ./cl-k8s
```

## Cool things

The binary spawned does NOT need to oversee the container process if you
run in detached mode with a PID file. You can have it watched by the user mode
systemd so that this binary is really just the launcher :)