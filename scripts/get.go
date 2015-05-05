package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"
	"text/template"
)

const (
	copylist = `
		{{.go}}/go/src
		{{.go}}/go/VERSION.cache
		{{.go}}/go/misc
		{{.go}}/go/pkg/include
		{{.go}}/go/bin/{{.kernel}}_{{.arch}}/go
		{{.go}}/go/pkg/tool/{{.kernel}}_{{.arch}}/{{.letter}}g
		{{.go}}/go/pkg/tool/{{.kernel}}_{{.arch}}/{{.letter}}l
		{{.go}}/go/pkg/tool/{{.kernel}}_{{.arch}}/asm
		{{.go}}/go/pkg/tool/{{.kernel}}_{{.arch}}/old{{.letter}}a
`
)
var (
	fail int
	t = template.Must(template.New("filelist").Parse(copylist))
)

func cp(in, out string) {
	b, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", in, err)
		fail++
		return
	}
	err = ioutil.WriteFile(out, b, 0444)
	if err != nil {
	     	fmt.Fprintf(os.Stderr, "%v: %v\n", out, err)
		fail++
	}
}

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

func main() {
	flag.Parse()
	arch := getenv("GOARCH", "amd64")
	goos := getenv("GOROOT", "/")
	d, err := ioutil.TempDir("", "u-root")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Copying for goos %v, arch %v, dir %v\n", goos, arch, d)
	
}

//#!/bin/bash
//# This is becoming more of a buildroot script. 
//# If it is a one-time thing, do it here.
//# The simplest thing is to build a go via
//# mount --bind your-go-src /go
//# cd /go
//# export CGO_ENABLED=0
//# cd src && make.bash
//# This gives you a go with the right baked-in paths.
//# This script assumes you have done that; if not,
//# have your arg be the path-to-go
//
//# Also, to further compress things, you can
//# cd /go/src/cmd/go
//# go install -tags cmd_go_bootstrap
//# the go_bootstrap will land in /go/tools/pkg/OS_ARCH as go_bootstrap
//# in the long term we'll use this, as it makes u-root workable.
//# This shrinks the go command by 50% or so.
//# You can later recreate the full command once booted:
//# cd /go/src/cmd/go
//# go install 
//
//# I can't believe I have to do this.
//# There are more compact forms (e.g. {6a,6c,6g,6l} but this
//# simple-minded format works with simple shells.
//# go/pkg used to contain binaries, and now contains .h files.
//# Hence the move to cpio. However, pulling the cpio into a ramfs
//# dramatically shortens compile times, so this is good.
//(
//find $1/go/src
//find $1/go/VERSION.cache
//find $1/go/misc
//find $1/go/pkg/include
//# go figure. It installs to somewhere else now.
//find $1/go/bin/linux_arm/go
//find $1/go/pkg/tool/linux_arm/5g
//find $1/go/pkg/tool/linux_arm/5l
//find $1/go/pkg/tool/linux_arm/asm
//find $1/go/pkg/tool/linux_arm/old5a
//) |
//(cpio --no-absolute-filenames -o) > armgo.cpio
//
//mkdir -p dev etc usr/lib lib64 tmp bin
//cp /etc/localtime etc
//
//sudo rm -f dev/null
//sudo mknod dev/null c 1 3
