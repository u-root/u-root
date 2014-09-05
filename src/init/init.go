// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// assumptions
// we've been booted into a ramfs with all this stuff unpacked and ready.
// we don't need a loop device mount because it's all there.
// So we run /go/bin/go build installcommand
// and then exec /buildbin/sh

package main

import (
	"os"
	"os/exec"
	"log"
)

var env = []string{
    	"PATH=/go/bin:/buildbin:/bin:/usr/local/bin:",
	"GOROOT=/go",
	"GOBIN=/bin",
	}

func main() {
	log.Printf("Welcome to u-root")
	cmd := exec.Command("/buildbin/sh")
	cmd.Env = env
	os.Setenv("GOROOT", "/go")
	os.Setenv("GOPATH", "/")
	os.Setenv("GOBIN", "/bin")
	os.Setenv("PATH", "/go/bin:/buildbin:/bin:/usr/local/bin")

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", cmd)
	err := cmd.Run()
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("init: /bin/sh returned!\n")
		for {
		}
}
