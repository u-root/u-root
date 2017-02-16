// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Build commands in backgorund processes. This feature certainly makes a
// non-busybox shell feel more realtime.
package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
)

// Generally shorter commands (ls, cd, ...) are used more often, so they are
// compiled first. This avoids maintaining a list for compilation order.
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
		args := []string{"--noexec", "--noforce", "--lowpri", cmdName}
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
