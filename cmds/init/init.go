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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	verbose = flag.Bool("v", false, "print all build commands")
	test    = flag.Bool("test", false, "Test mode: don't try to set control tty")
	debug   = func(string, ...interface{}) {}
)

func main() {
	a := []string{"build"}
	flag.Parse()
	log.Printf("Welcome to u-root")
	util.Rootfs()

	if *verbose {
		debug = log.Printf
		a = append(a, "-x")
	}

	// populate buildbin

	// In earlier versions we just had src/cmds. Due to the Go rules it seems we need to
	// embed the URL of the repo everywhere. Yuck.
	c, err := filepath.Glob("/src/github.com/u-root/u-root/cmds/[a-z]*")
	if err != nil || len(c) == 0 {
		log.Printf("In a break with tradition, you seem to have NO u-root commands: %v", err)
	}
	o, err := filepath.Glob("/src/*/*/*")
	if err != nil {
		log.Printf("Your filepath glob for other commands seems busted: %v", err)
	}
	c = append(c, o...)
	for _, v := range c {
		name := filepath.Base(v)
		if name == "installcommand" || name == "init" {
			continue
		} else {
			destPath := filepath.Join("/buildbin", name)
			source := "/buildbin/installcommand"
			if err := os.Symlink(source, destPath); err != nil {
				log.Printf("Symlink %v -> %v failed; %v", source, destPath, err)
			}
		}
	}

	envs := os.Environ()
	debug("envs %v", envs)
	os.Setenv("GOBIN", "/buildbin")
	a = append(a, "-o", "/buildbin/installcommand", filepath.Join(util.CmdsPath, "installcommand"))
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

	// Before entering an interactive shell, decrease the loglevel because
	// spamming non-critical logs onto the shell frustrates users. The logs
	// are still accessible through dmesg.
	const sysLogActionConsoleLevel = 8
	const kernNotice = 5 // Only messages more severe than "notice" are printed.
	if _, _, err := syscall.Syscall(syscall.SYS_SYSLOG, sysLogActionConsoleLevel, 0, kernNotice); err != 0 {
		log.Print("Could not set log level")
	}

	// install /env.
	os.Setenv("GOBIN", "/ubin")
	envs = append(envs, "GOBIN=/ubin")
	for _, e := range envs {
		nv := strings.SplitN(e, "=", 2)
		if len(nv) < 2 {
			nv = append(nv, "")
		}
		n := filepath.Join("/env", nv[0])
		if err := ioutil.WriteFile(n, []byte(nv[1]), 0666); err != nil {
			log.Printf("%v: %v", n, err)
		}
	}

	var profile string
	// Some systems wipe out all the environment variables we so carefully craft.
	// There is a way out -- we can put them into /etc/profile.d/uroot if we want.
	// The PATH variable has to change, however.
	path := fmt.Sprintf("%v:%v:%v:%v", util.GoBin(), util.PATHHEAD, "$PATH", util.PATHTAIL)
	for k, v := range util.Env {
		// We're doing the hacky way for now. We can clean this up later.
		if k == "PATH" {
			profile += "export PATH=" + path + "\n"
		} else {
			profile += "export " + k + "=" + v + "\n"
		}
	}

	// The IFS lets us force a rehash every time we type a command, so that when we
	// build uroot commands we don't keep rebuilding them.
	profile += "IFS=`hash -r`\n"

	// IF the profile is used, THEN when the user logs in they will need a
	// private tmpfs. There's no good way to do this on linux. The closest
	// we can get for now is to mount a tmpfs of /go/pkg/%s_%s :-( Same
	// applies to ubin. Each user should have their own.
	profile += fmt.Sprintf("sudo mount -t tmpfs none /go/pkg/%s_%s\n", runtime.GOOS, runtime.GOARCH)
	profile += fmt.Sprintf("sudo mount -t tmpfs none /ubin\n")
	profile += fmt.Sprintf("sudo mount -t tmpfs none /pkg\n")

	// Now here's some good fun. We've set environment variables we want to see used.
	// But on some systems the environment variable we create is completely ignored.
	// Oh, is that you again, tinycore? Well.
	// So we can save the day by writing the profile string to /etc/profile.d/uroot.sh
	// mode, usually, 644.
	// Only bother doing this is /etc/profile.d exists and is a directory.
	if fi, err := os.Stat("/etc/profile.d"); err == nil && fi.IsDir() {
		if err := ioutil.WriteFile("/etc/profile.d/uroot.sh", []byte(profile), 0644); err != nil {
			log.Printf("Trying to write uroot profile failed: %v", err)
		}
	}

	// Start background build.
	if isBgBuildEnabled() {
		go startBgBuild()
	}

	// There may be an inito if we are building on
	// an existing initramfs. So, first, try to
	// run inito and then run our shell
	// inito is always first and we set default flags for it.
	cloneFlags := uintptr(syscall.CLONE_NEWPID)
	cmdList := []string{"/inito", "/buildbin/uinit", "/buildbin/rush"}
	noCmdFound := true
	for _, v := range cmdList {
		if _, err := os.Stat(v); !os.IsNotExist(err) {
			noCmdFound = false
			cmd = exec.Command(v)
			cmd.Env = envs
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if *test {
				cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: cloneFlags}
			} else {
				cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true, Cloneflags: cloneFlags}
			}
			debug("Run %v", cmd)
			if err := cmd.Run(); err != nil {
				log.Print(err)
			}
		}
		// only the first init needs its own PID space.
		cloneFlags = 0
	}

	if noCmdFound {
		log.Printf("init: No suitable executable found in %+v", cmdList)
	}

	log.Printf("init: All commands exited")
	log.Printf("init: Syncing filesystems")
	syscall.Sync()
	log.Printf("init: Exiting...")
}
