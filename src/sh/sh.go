// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
sh reads in a line at a time and runs it. 
prompt is '% '
*/

package main

import (
	"path/filepath"
	"os/exec"
	"fmt"
	"os"
	"strings"
	"bufio"
)

var urpath = "/go/bin:/buildbin:/bin:/usr/local/bin:"

func main() {
	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%% ")
	for scanner.Scan() {
		cmd := scanner.Text()
		argv := strings.Split(cmd, " ")
		globargv := []string{}
		for _, v := range argv[1:] {
			if globs, err := filepath.Glob(v); err == nil {
				globargv = append(globargv, globs...)
			} else {
				globargv = append(globargv, v)
			}
		}
		run := exec.Command(argv[0], globargv[:]...)
		run.Stdin = os.Stdin
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		if err := run.Start(); err != nil {
			fmt.Printf("%v: Path %v\n", err, os.Getenv("PATH"))
		} else if err := run.Wait(); err != nil {
			fmt.Printf("wait: %v:\n", err)
		}
		fmt.Printf("%% ")
	}
}
