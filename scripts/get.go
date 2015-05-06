package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	copylist = `{{.Goroot}}/go/src
{{.Goroot}}/go/VERSION.cache
{{.Goroot}}/go/misc
{{.Goroot}}/go/pkg/include
{{.Goroot}}/go/bin/{{.Goos}}_{{.Arch}}/go
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}g
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}l
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/asm
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/old{{.Letter}}a
`
)
var (
	fail int
	t = template.Must(template.New("filelist").Parse(copylist))
	letter = map[string]string{
		"amd64": "6",
		"arm": "5",
		"ppc": "9",
		}
)

func cp(in, out string) error {
	b, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", in, err)
		fail++
		return err
	}
	err = ioutil.WriteFile(out, b, 0444)
	if err != nil {
	     	fmt.Fprintf(os.Stderr, "%v: %v\n", out, err)
		fail++
	}
	return nil
}

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

func main() {
	type config struct {
		Goroot string
		Arch string
		Goos string
		Letter string
	}
	var a config
	flag.Parse()
	var err error
	a.Arch = getenv("GOARCH", "amd64")
	a.Goroot = getenv("GOROOT", "/")
	a.Goos = "linux"
	f, err := ioutil.TempFile("", "u-root")
	fmt.Fprintf(os.Stderr, "f is %v\n", f.Name())
	a.Letter = letter[a.Arch]
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Copying for Goos %v, Arch %v\n", a.Goos, a.Arch)
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("cpio", "--verbose", "-o")
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = f
	fmt.Fprintf(os.Stderr, "Run %v", cmd)
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	n := strings.Split(b.String(), "\n")
	fmt.Fprintf(os.Stderr, "Strings :%v:\n", n)
	for _, v := range n {
		fmt.Fprintf(os.Stderr, "%v\n", v)
		err := filepath.Walk(v, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%v: %v\n", path, err)
				return err
			}
			fmt.Fprintf(w, "%v\n", path)
			fmt.Printf("%v\n", path)
			return err
		})
		fmt.Printf("WALKED %v\n", v)
		if err != nil {
			fmt.Printf("%s: %v\n", v, err)
		}
	}
	w.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	// It worked, so move the file to where it needs to be.
	cpioName := fmt.Sprintf("%v_%vgo.cpio", a.Goos, a.Arch)
	if err := cp(f.Name(), cpioName); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Fprintf(os.Stderr, "Successfully created CPIO %v\n", cpioName)
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
