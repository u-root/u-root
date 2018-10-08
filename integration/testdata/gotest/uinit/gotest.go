// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/u-root/u-root/pkg/sh"
)

// Mount a vfat volume and run the tests within.
func main() {
	sh.RunOrDie("mkdir", "/testdata")
	sh.RunOrDie("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")

	// Gather list of tests.
	files, err := ioutil.ReadDir("/testdata/tests")
	if err != nil {
		log.Fatal(err)
	}
	tests := []string{}
	for _, f := range files {
		tests = append(tests, f.Name())
	}

	// Sort tests.
	sort.Strings(tests)

	// We are using TAP-style test output. See: https://testanything.org/
	// One unfortunate design in TAP is "ok" is a subset of "not ok", so we
	// prepend each line with "TAP: " and search for for "TAP: ok".
	log.Printf("TAP: 1..%d", len(tests))

	// Run tests.
	for i, t := range tests {
		runMsg := fmt.Sprintf("TAP: # running %d - %s", i, t)
		passMsg := fmt.Sprintf("TAP: ok %d - %s", i, t)
		failMsg := fmt.Sprintf("TAP: not ok %d - %s", i, t)
		log.Println(runMsg)

		ctx, cancel := context.WithTimeout(context.Background(), 1900*time.Millisecond)
		defer cancel()
		cmd := exec.CommandContext(ctx, filepath.Join("/testdata/tests", t))
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err := cmd.Run()

		if ctx.Err() == context.DeadlineExceeded {
			log.Println("TAP: # timed out")
			log.Println(failMsg)
		} else if err == nil {
			log.Println(passMsg)
		} else {
			log.Println(err)
			log.Println(failMsg)
		}
	}
}
