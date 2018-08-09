// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func sh(arg0 string, args ...string) {
	cmd := exec.Command(arg0, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Mount a vfat volume and run the tests within.
func main() {
	sh("mkdir", "/testdata")
	sh("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")

	// Gather list of tests.
	files, err := ioutil.ReadDir("/testdata/tests")
	if err != nil {
		log.Fatal(err)
	}
	tests := []string{}
	for _, f := range files {
		tests = append(tests, f.Name())
	}

	// Sort and run tests.
	sort.Strings(tests)
	for _, t := range tests {
		cmd := exec.Command(filepath.Join("/testdata/tests", t))
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		runMsg := fmt.Sprintf("#### %s RUN ####", t)
		passMsg := fmt.Sprintf("#### %s PASSED ####", t)
		failMsg := fmt.Sprintf("#### %s FAILED ####", t)

		log.Println(runMsg)
		if err := cmd.Run(); err == nil {
			log.Println(passMsg)
		} else {
			log.Println(err)
			log.Println(failMsg)
		}
	}
}
