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
	"log"
	"os"
	"os/exec"
)

var env = map[string]string{
	"PATH": "/go/bin:/bin:/buildbin:/usr/local/bin:",
	"LD_LIBRARY_PATH" : "/usr/local/lib",
	"GOROOT": "/go",
	"GOBIN": "/bin",
	"GOPATH": "/",
	"CGO_ENABLED": "0",
}

func main() {
	log.Printf("Welcome to u-root 4")
	envs := []string{}
	for k, v := range env {
		os.Setenv(k, v)
		envs =  append(envs, k + "=" + v)
	}
	log.Printf("envs %v", envs)
	cmd := exec.Command("go", "install", "-x", "installcommand")
	cmd.Env = envs
	cmd.Dir = "/"

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", cmd)
	err := cmd.Run()
	if err != nil {
		log.Printf("%v\n", err)
	}


	cmd = exec.Command("/buildbin/sh")
	cmd.Env = envs
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", cmd)
	err = cmd.Run()
	if err != nil {
		log.Printf("%v\n", err)
	}
	log.Printf("init: /bin/sh returned!\n")
	for {
	}
}
