// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
sh reads in a line at a time and runs it. 
prompt is '% '
*/

package main

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"bufio"
)

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
		e := os.Environ()
		e = append(e, "GOROOT=/go")
		e = append(e, "GOPATH=/")
		e = append(e, "GOBIN=/bin")
		run := exec.Command(argv[0], argv[1:]...)
		run.Env = e
		out, err := run.CombinedOutput()
		if err != nil {
			fmt.Printf("%v: Path %v\n", err, os.Getenv("PATH"))
		}
		fmt.Printf("%s", out)
		fmt.Printf("%% ")
	}
}
