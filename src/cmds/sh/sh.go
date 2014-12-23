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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type builtin func(c *Command) error

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

func runit(c *Command) error {
	if b, ok := builtins[c.cmd]; ok {
		if err := b(c); err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		run := exec.Command(c.cmd, c.argv[:]...)
		run.Stdin = c.in
		run.Stdout = c.out
		run.Stderr = c.err
		if err := run.Start(); err != nil {
			return errors.New(fmt.Sprintf("%v: Path %v\n", err, os.Getenv("PATH")))
		} else if err := run.Wait(); err != nil {
			return errors.New(fmt.Sprintf("wait: %v:\n", err))
		}
	}
	return nil

}

func OpenRead(c *Command, r io.Reader, fd int) (io.Reader, error) {
	if c.fdmap[fd] != "" {
		return os.Open(c.fdmap[fd])
	}
	return r, nil
}
func OpenWrite(c *Command, w io.Writer, fd int) (io.Writer, error) {
	if c.fdmap[fd] != "" {
		return os.Create(c.fdmap[fd])
	}
	return w, nil
}
func command(c *Command) error {
	// IO defaults.
	var err error
	if c.in, err = OpenRead(c, os.Stdin, 0); err != nil {
		return err
	}
	if c.out, err = OpenWrite(c, os.Stdout, 1); err != nil {
		return err
	}
	if c.err, err = OpenWrite(c, os.Stderr, 2); err != nil {
		return err
	}

	globargv := []string{}
	for _, v := range c.args {
		if v.mod == "ENV" {
			e := v.val
			if !path.IsAbs(v.val) {
				e = path.Join(envDir, e)
			}
			b, err := ioutil.ReadFile(e)
			if err != nil {
				return err
			}
			// It goes in as one argument. Not sure if this is what we want
			// but it gets very weird to start splitting it on spaces. Or maybe not?
			globargv = append(globargv, string(b))
		} else if globs, err := filepath.Glob(v.val); err == nil && len(globs) > 0 {
			globargv = append(globargv, globs...)
		} else {
			globargv = append(globargv, v.val)
		}
	}

	c.cmd = globargv[0]
	c.argv = globargv[1:]
	// for now, bg will just happen in background.
	if c.bg {
		go func() {
			if err := runit(c); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
		}()
	} else {
		err := runit(c)
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
