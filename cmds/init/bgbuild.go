// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Build commands in background processes. This feature certainly makes a
// non-busybox shell feel more realtime.
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

// Commands are built approximately in order from smallest to largest length of
// the command name. So, two letter commands like `ls` and `cd` will be built
// before `mknod` and `mount`. Generally, shorter commands are used more often
// (that is why they were made short) and are more likely to be used first,
// thus they should be built first.
type cmdSlice []os.FileInfo

func (p cmdSlice) Len() int {
	return len(p)
}

func (p cmdSlice) Less(i, j int) bool {
	return len(p[i].Name()) < len(p[j].Name())
}

func (p cmdSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func worker(cmds chan string) {
	for cmdName := range cmds {
		args := []string{"--onlybuild", "--noforce", "--lowpri", cmdName}
		cmd := exec.Command("installcommand", args...)
		if err := cmd.Start(); err != nil {
			log.Println("Cannot start:", err)
			continue
		}
		if err := cmd.Wait(); err != nil {
			log.Println("installcommand error:", err)
		}
	}
}

func startBgBuild() {
	cmds := make(chan string)

	// Start a worker for each CPU.
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(cmds)
	}

	// Create a slice of commands and order them by the aformentioned
	// heuristic.
	fis, err := ioutil.ReadDir("/buildbin")
	if err != nil {
		log.Print(err)
		return
	}
	sort.Sort(cmdSlice(fis))

	// Send the command by name to the workers.
	for _, fi := range fis {
		cmds <- fi.Name()
	}
	close(cmds)
}

func cmdlineContainsFlag(flag string) bool {
	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		// Procfs not mounted?
		return false
	}
	args := strings.Split(string(cmdline), " ")
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func isBgBuildEnabled() bool {
	return !cmdlineContainsFlag("uroot.nobgbuild")
}
