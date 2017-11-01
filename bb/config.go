// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func guessgoarch() {
	if arch := os.Getenv("GOARCH"); arch != "" {
		config.Arch = filepath.Clean(arch)
	} else {
		config.Arch = runtime.GOARCH
	}
}

func guessgoroot() {
	if root := os.Getenv("GOROOT"); root != "" {
		config.Goroot = filepath.Clean(root)
	} else {
		config.Goroot = runtime.GOROOT()
	}
	config.Gosrcroot = filepath.Dir(config.Goroot)
	log.Printf("Using %q as GOROOT", config.Goroot)
}

func guessgopath() {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		config.Gopath = gopath
		return
	}
	log.Fatalf("You have to set GOPATH, which is typically ~/go")
}

func doConfig() {
	var err error
	flag.BoolVar(&config.Debug, "d", false, "Debugging")
	flag.Parse()
	if config.Debug {
		debug = debugPrint
	}
	if config.Cwd, err = os.Getwd(); err != nil {
		log.Fatalf("Getwd: %v", err)
	}

	guessgoroot()
	guessgopath()
	guessgoarch()
	config.Gosrcroot = filepath.Dir(config.Goroot)
	config.Goos = "linux"
	config.TempDir, err = ioutil.TempDir("", "u-root")
	config.Go = ""
	if err != nil {
		log.Fatalf("%v", err)
	}
	config.Bbsh = filepath.Join(config.Cwd, "bbsh")
	os.RemoveAll(config.Bbsh)
}
