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
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type builtin func(string, []string) error

var (
	urpath   = "/go/bin:/buildbin:/bin:/usr/local/bin:"
	builtins = make(map[string]builtin)
	// the environment dir is INTENDED to be per-user and bound in
	// a private name space at /env.
	envDir = "/env"
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

func command(c *Command) error {
	globargv := []string{}
	for _, v := range c.args[1:] {
		if v.mod == "ENV" {
			// Later, this will involve substitution.
			e := v.val
			if !path.IsAbs(v.val) {
				e = path.Join(envDir, e)
			}
			b, err := ioutil.ReadFile(e)
			if err != nil {
				return err
			}
			globargv = append(globargv, string(b))
		} else if globs, err := filepath.Glob(v.val); err == nil && len(globs) > 0 {
			globargv = append(globargv, globs...)
		} else {
			globargv = append(globargv, v.val)
		}
	}

	// for now, bg will just happen in background.
	if c.bg {
		go func() {
			if err := runit(c.args[0].val, globargv); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
		}()
	} else {
		err := runit(c.args[0].val, globargv)
		return err
	}
	return nil
}
func main() {
	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	b := bufio.NewReader(os.Stdin)
	fmt.Printf("%% ")
	for {
		cmds, status, err := getCommand(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		if len(cmds) > 1 {
			fmt.Fprintf(os.Stderr, "no compounds yet\n")
		}
		// Once we get to compounds this will be a lot more complex, of course.
		// And there's no redirection at this point, sorry.
		if len(cmds) > 0 {
			if err := command(cmds[0]); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
		if status == "EOF" {
			break
		}
		fmt.Printf("%% ")
	}
}
