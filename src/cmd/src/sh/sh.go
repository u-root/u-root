// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
sh reads in a line at a time and runs it.
prompt is '% '
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type builtin func(string, []string) error

var (
	urpath   = "/go/bin:/buildbin:/bin:/usr/local/bin:"
	builtins = make(map[string]builtin)
)

func addBuiltIn(name string, f builtin) error {
	if _, ok := builtins[name]; ok {
		return errors.New(fmt.Sprintf("%v already a builtin", name))
	}
	builtins[name] = f
	return nil
}

func runit(cmd string, argv []string) error {
	if b, ok := builtins[cmd]; ok {
		if err := b(cmd, argv); err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		run := exec.Command(cmd, argv[:]...)
		run.Stdin = os.Stdin
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		if err := run.Start(); err != nil {
			return errors.New(fmt.Sprintf("%v: Path %v\n", err, os.Getenv("PATH")))
		} else if err := run.Wait(); err != nil {
			return errors.New(fmt.Sprintf("wait: %v:\n", err))
		}
	}
	return nil

}

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
		if len(cmd) == 0 {
			fmt.Printf("%% ")
			continue
		}
		globargv := []string{}
		for _, v := range argv[1:] {
			if globs, err := filepath.Glob(v); err == nil && len(globs) > 0 {
				globargv = append(globargv, globs...)
			} else {
				globargv = append(globargv, v)
			}
		}
		err := runit(argv[0], globargv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		fmt.Printf("%% ")
	}
}
