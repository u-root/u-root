#!/bin/bash

export PATH="/go/bin:/ubin:$PATH/buildbin:/bbin"

# The IFS lets us force a rehash every time we type a command, so that when we
# build uroot commands we don't keep rebuilding them.
IFS=`hash -r`

# TODO: why do they need this private tmpfs?
#
# IF the profile is used, THEN when the user logs in they will need a
# private tmpfs. There's no good way to do this on linux. The closest
# we can get for now is to mount a tmpfs of /go/pkg/%s_%s :-( Same
# applies to ubin. Each user should have their own.
sudo mount -t tmpfs none /go/pkg/linux_amd64
sudo mount -t tmpfs none /go/pkg/linux_arm64
sudo mount -t tmpfs none /go/pkg/linux_arm
sudo mount -t tmpfs none /ubin
sudo mount -t tmpfs none /pkg
