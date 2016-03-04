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
	"github.com/u-root/u-root/uroot"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	log.Printf("Welcome to u-root")
	uroot.Rootfs()
	// populate buildbin

	if commands, err := ioutil.ReadDir("/src/cmds"); err == nil {
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
			}
		}
	} else {
		log.Fatalf("Can't read %v; %v", "/src", err)
	}
	envs := uroot.Envs
	log.Printf("envs %v", envs)
	os.Setenv("GOBIN", "/buildbin")
	cmd := exec.Command("go", "install", "-x", path.Join("cmds", "installcommand"))
	installenvs := envs
	installenvs = append(envs, "GOBIN=/buildbin")
	cmd.Env = installenvs
	cmd.Dir = "/"

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v", cmd)
	err := cmd.Run()
	if err != nil {
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

	// There may be an init.orig if we are building on
	// an existing initramfs. So, first, try to
	// run init.orig and then run our shell
	// Perhaps we should stat init.orig first.
	for _, v := range []string{"/init.orig", "/buildbin/rush"} {
		cmd = exec.Command(v)
		cmd.Env = envs
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		// TODO: figure out why we get EPERM when we use this.
		//cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true,}
		log.Printf("Run %v", cmd)
		err = cmd.Run()
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
	log.Printf("init: All commands exited")
}
