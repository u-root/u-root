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
	"path/filepath"
	"strings"
	"syscall"
	"uroot"
)

const PATH = "/bin:/buildbin:/usr/local/bin"

type dir struct {
	name string
	mode os.FileMode
}

type dev struct {
	name  string
	mode  os.FileMode
	magic int
	howmany int
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

	dirs = []dir{
		{name: "/proc", mode: os.FileMode(0555)},
		{name: "/buildbin", mode: os.FileMode(0777)},
		{name: "/bin", mode: os.FileMode(0777)},
		{name: "/tmp", mode: os.FileMode(0777)},
		{name: "/env", mode: os.FileMode(0777)},
		{name: "/etc", mode: os.FileMode(0777)},
		{name: "/tcz", mode: os.FileMode(0777)},
		{name: "/dev", mode: os.FileMode(0777)},
		{name: "/lib", mode: os.FileMode(0777)},
		{name: "/usr/lib", mode: os.FileMode(0777)},
		{name: "/go/pkg/linux_amd64", mode: os.FileMode(0777)},
	}
	devs = []dev{
		// chicken and egg: these need to be there before you start. So, sadly,
		// we will always need dev.cpio. 
		//{name: "/dev/null", mode: os.FileMode(0660) | 020000, magic: 0x0103},
		//{name: "/dev/console", mode: os.FileMode(0660) | 020000, magic: 0x0501},
	}
	namespace = []mount{
		{source: "proc", target: "/proc", fstype: "proc", flags: syscall.MS_MGC_VAL | syscall.MS_RDONLY, opts: ""},
	}
)

func main() {
	log.Printf("Welcome to u-root")
	// Pick some reasonable values in the (unlikely!) even that Uname fails.
	uname := "linux"
	mach := "x86_64"
	// There are three possible places for go:
	// The first is in /go/bin/$OS_$ARCH
	// The second is in /go/bin [why they still use this path is anyone's guess]
	// The third is in /go/pkg/tool/$OS_$ARCH
	if u, err := uroot.Uname(); err != nil {
		log.Printf("uroot.Utsname fails: %v, so assume %v_%v\n", uname, mach)
	} else {
		// Sadly, go and the OS disagree on case.
		uname = strings.ToLower(u.Sysname)
		mach = strings.ToLower(u.Machine)
		// Yes, we really have to do this stupid thing.
		if mach[0:3] == "arm" {
			mach = "arm"
		}
	}
	env["PATH"] = fmt.Sprintf("/go/bin/%s_%s:/go/bin:/go/pkg/tool/%s_%s:%v", uname, mach, uname, mach, PATH)
	envs := []string{}
	for k, v := range env {
		os.Setenv(k, v)
		envs = append(envs, k+"="+v)
	}

	for _, m := range dirs {
		if err := os.MkdirAll(m.name, m.mode); err != nil {
			log.Printf("mkdir :%s: mode %o: %v\n", m.name, m.mode, err)
			continue
		}
	}

	for _, d := range devs {
		syscall.Unlink(d.name)
		if err := syscall.Mknod(d.name, uint32(d.mode), d.magic); err != nil {
			log.Printf("mknod :%s: mode %o: magic: %v: %v\n", d.name, d.mode, d.magic, err)
			continue
		}
	}

	for _, m := range namespace {
		if err := syscall.Mount(m.source, m.target, m.fstype, m.flags, m.opts); err != nil {
			log.Printf("Mount :%s: on :%s: type :%s: flags %x: %v\n", m.source, m.target, m.fstype, m.flags, m.opts, err)
		}

	}

	// only in case of emergency.
	if false {
		if err := filepath.Walk("/", func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf(" WALK FAIL%v: %v\n", name, err)
				// That's ok, sometimes things are not there.
				return nil
			}
			fmt.Printf("%v\n", name)
			return nil
		}); err != nil {
			log.Printf("WALK fails %v\n", err)
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
