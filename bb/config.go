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
	"strings"
)

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

// TODO: put this in the uroot package
// It's annoying asking them to set lots of things. So let's try to figure it out.
func guessgoroot() {
	config.Goroot = os.Getenv("GOROOT")
	if config.Goroot != "" {
		log.Printf("Using %v as GOROOT from environment variable", config.Goroot)
		config.Gosrcroot = filepath.Dir(config.Goroot)
		return
	}
	log.Print("Goroot is not set, trying to find a go binary")
	p := os.Getenv("PATH")
	paths := strings.Split(p, ":")
	for _, v := range paths {
		g := filepath.Join(v, "go")
		if _, err := os.Stat(g); err == nil {
			config.Goroot = filepath.Dir(filepath.Dir(v))
			config.Gosrcroot = filepath.Dir(config.Goroot)
			log.Printf("Guessing that goroot is %v", config.Goroot)
			return
		}
	}
	log.Printf("GOROOT is not set and can't find a go binary in %v", p)
	config.Fail = true
}

func guessgopath() {
	defer func() {
		config.Gosrcroot = filepath.Dir(config.Goroot)
	}()
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		config.Gopath = filepath.Clean(gopath)
		return
	}
	// We need to change the guess logic but that will have to wait.
	log.Fatal("GOPATH was not set")
	// It's a good chance they're running this from the u-root source directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("GOPATH was not set and I can't get the wd: %v", err)
		config.Fail = true
		return
	}
	// walk up the cwd until we find a u-root entry. See if cmds/init/init.go exists.
	for c := cwd; c != "/"; c = filepath.Dir(c) {
		if filepath.Base(c) != "u-root" {
			continue
		}
		check := filepath.Join(c, "cmds/init/init.go")
		if _, err := os.Stat(check); err != nil {
			//log.Printf("Could not stat %v", check)
			continue
		}
		config.Gopath = c
		log.Printf("Guessing %v as GOPATH", c)
		os.Setenv("GOPATH", c)
		return
	}
	config.Fail = true
	log.Printf("GOPATH was not set, and I can't see a u-root-like name in %v", cwd)
	return
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
	config.Arch = getenv("GOARCH", "amd64")
	if config.Fail {
		os.Exit(1)
	}
	config.Gosrcroot = filepath.Dir(config.Goroot)
	config.Goos = "linux"
	config.TempDir, err = ioutil.TempDir("", "u-root")
	config.Go = ""
	if err != nil {
		log.Fatalf("%v", err)
	}
	config.Bbsh = filepath.Join(config.Cwd, "bbsh")
	os.RemoveAll(config.Bbsh)
	config.Args = flag.Args()
	if len(config.Args) == 0 {
		config.Args = defaultCmd
	}

}
