// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// assumptions
// we've been booted into a ramfs with all this stuff unpacked and ready.
// we don't need a loop device mount because it's all there.
// So we run /go/bin/go build installcommand
// and then exec /buildbin/sh

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/u-root/u-root/uroot"
)

var (
	verbose   = flag.Bool("v", false, "print all build commands")
	ludicrous = flag.Bool("ludicrous", false, "print out information about symlink creation")
	debug     = func(string, ...interface{}) {}
)

func main() {
	a := []string{"build"}
	flag.Parse()
	log.Printf("Welcome to u-root")
	uroot.Rootfs()

	if *verbose {
		debug = log.Printf
		a = append(a, "-x")
	}

	// populate buildbin

	// In earlier versions we just had src/cmds. Due to the Go rules it seems we need to
	// embed the URL of the repo everywhere. Yuck.
	if commands, err := ioutil.ReadDir(path.Join("/src", uroot.CmdsPath)); err == nil {
		for _, v := range commands {
			name := v.Name()
			if name == "installcommand" || name == "init" {
				continue
			} else {
				destPath := path.Join("/buildbin", name)
				source := "/buildbin/installcommand"
				if err := os.Symlink(source, destPath); err != nil {
					log.Printf("Symlink %v -> %v failed; %v", source, destPath, err)
				}
				if *ludicrous {
					log.Printf("Symlink %v -> %v", source, destPath)
				}
			}
		}
	} else {
		log.Fatalf("Can't read %v; %v", "/src", err)
	}
	envs := uroot.Envs
	debug("envs %v", envs)
	os.Setenv("GOBIN", "/buildbin")
	a = append(a, "-o", "/buildbin/installcommand", path.Join(uroot.CmdsPath, "installcommand"))
	cmd := exec.Command("go", a...)
	installenvs := envs
	installenvs = append(envs, "GOBIN=/buildbin")
	cmd.Env = installenvs
	cmd.Dir = "/"

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	debug("Run %v", cmd)
	if err := cmd.Run(); err != nil {
		log.Printf("%v\n", err)
	}

	// install /env.
	os.Setenv("GOBIN", "/ubin")
	envs = append(envs, "GOBIN=/ubin")
	for _, e := range envs {
		nv := strings.SplitN(e, "=", 2)
		if len(nv) < 2 {
			nv = append(nv, "")
		}
		n := path.Join("/env", nv[0])
		if err := ioutil.WriteFile(n, []byte(nv[1]), 0666); err != nil {
			log.Printf("%v: %v", n, err)
		}
	}

	// Now here's some good fun. We've set environment variables we want to see used.
	// But on some systems the environment variable we create is completely ignored.
	// Oh, is that you again, tinycore? Well.
	// So we can save the day by writing the uroot.profile string to /etc/profile.d/uroot.sh
	// mode, usually, 644.
	// Only bother doing this is /etc/profile.d exists and is a directory.
	if fi, err := os.Stat("/etc/profile.d"); err == nil && fi.IsDir() {
		if err := ioutil.WriteFile("/etc/profile.d/uroot.sh", []byte(uroot.Profile), 0644); err != nil {
			log.Printf("Trying to write uroot profile failed: %v", err)
		}
	}

	// There may be an inito if we are building on
	// an existing initramfs. So, first, try to
	// run inito and then run our shell
	// Perhaps we should stat inito first.
	// inito is always first and we set default flags for it.
	cloneFlags := uintptr(syscall.CLONE_NEWPID)
	for _, v := range []string{"/inito", "/buildbin/rush"} {
		cmd = exec.Command(v)
		cmd.Env = envs
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		// TODO: figure out why we get EPERM when we use this.
		//cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true, Cloneflags: cloneFlags}
		cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: cloneFlags}
		debug("Run %v", cmd)
		if err := cmd.Run(); err != nil {
			log.Printf("%v\n", err)
		}
		// only the first init needs its own PID space.
		cloneFlags = 0
	}
	log.Printf("init: All commands exited")
}
