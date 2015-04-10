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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"uroot"
)

const PATH = "/bin:/buildbin:/usr/local/bin"

type dir struct {
	name string
	mode   os.FileMode
}

type mount struct {
	source string
	target string
	fstype string
	flags  uintptr
	opts   string
}

var (
	env = map[string]string{
		"LD_LIBRARY_PATH": "/usr/local/lib",
		"GOROOT":          "/go",
		"GOPATH":          "/",
		"CGO_ENABLED":     "0",
	}

	dirs = []dir {
		{name: "/proc", mode: os.FileMode(0555),},
		{name: "/buildbin", mode: os.FileMode(0777),},
		{name: "/bin", mode: os.FileMode(0777),},
		{name: "/tmp", mode: os.FileMode(0777),},
		{name: "/env", mode: os.FileMode(0777),},
		{name: "/etc", mode: os.FileMode(0777),},
		{name: "/tcz", mode: os.FileMode(0777),},
		{name: "/dev", mode: os.FileMode(0777),},
		{name: "/lib", mode: os.FileMode(0777),},
		{name: "/usr/lib", mode: os.FileMode(0777),},
		{name: "/go/pkg/linux_amd64", mode: os.FileMode(0777),},
	}
	namespace = []mount{
		{source: "proc", target: "/proc", fstype: "proc", flags: syscall.MS_MGC_VAL | syscall.MS_RDONLY, opts: "",},
	}
)

func main() {
	log.Printf("Welcome to u-root")
	// Pick some reasonable values in the (unlikely!) even that Uname fails.
	uname := "linux"
	mach := "x86_64"
	// There are two possible places for go:
	// The first is in /go/bin
	// The second is in /go/pkg/tool/$OS_$ARCH
	if u, err := uroot.Uname(); err != nil {
		log.Printf("uroot.Utsname fails: %v, so assume %v_%v\n", uname, mach)
	} else {
		// Sadly, go and the OS disagree on case.
		uname = strings.ToLower(u.Sysname)
	}
	env["PATH"] = fmt.Sprintf("/go/bin:/go/pkg/tool/%s_%s:%v", uname, mach, PATH)
	envs := []string{}
	for k, v := range env {
		os.Setenv(k, v)
		envs = append(envs, k+"="+v)
	}

	for _, m := range dirs {
		if err := os.MkdirAll(m.name, m.mode); err != nil {
			log.Printf("mkdir :%s: mode %o: %v\n", m.name, m.mode)
			continue
		}
	}

	for _, m := range namespace {
		if err := syscall.Mount(m.source, m.target, m.fstype, m.flags, m.opts); err != nil {
			log.Printf("Mount :%s: on :%s: type :%s: flags %x: %v\n", m.source, m.target, m.fstype, m.flags, m.opts, err)
		}

	}
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
	os.Setenv("GOBIN", "/bin")
	envs = append(envs, "GOBIN=/bin")
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

	cmd = exec.Command("/buildbin/sh")
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
	log.Printf("init: /bin/sh returned!\n")
}
