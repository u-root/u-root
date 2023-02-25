// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Rush is an interactive shell similar to sh.
//
// Description:
//
//	Prompt is '% '.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type builtin func(c *Command) error

// TODO: probably have one builtin map and use it for both types?
var (
	urpath = "/go/bin:/ubin:/buildbin:/bbin:/bin:/usr/local/bin:"

	builtins map[string]builtin
)

func addBuiltIn(name string, f builtin) error {
	if builtins == nil {
		builtins = make(map[string]builtin)
	}
	if _, ok := builtins[name]; ok {
		return fmt.Errorf("%v already a builtin", name)
	}
	builtins[name] = f
	return nil
}

func wire(cmds []*Command) error {
	for i, c := range cmds {
		// IO defaults.
		var err error
		if c.Stdin == nil {
			if c.Stdin, err = openRead(c, os.Stdin, 0); err != nil {
				return err
			}
		}
		if c.Link != "|" {
			if c.Stdout, err = openWrite(c, os.Stdout, 1); err != nil {
				return err
			}
		}
		if c.Stderr, err = openWrite(c, os.Stderr, 2); err != nil {
			return err
		}
		// The validation is such that "|" is not set on the last one.
		// Also, there won't be redirects and "|" inappropriately.
		if c.Link != "|" {
			continue
		}
		w, err := cmds[i+1].StdinPipe()
		if err != nil {
			return err
		}
		r, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		// Oh, yuck.
		// There seems to be no way to do the classic
		// inherited pipes thing in Go. Hard to believe.
		go func() {
			_, _ = io.Copy(w, r)
			w.Close()
		}()
	}
	return nil
}

func runit(c *Command) error {
	defer func() {
		for fd, f := range c.files {
			f.Close()
			delete(c.files, fd)
		}
	}()
	if b, ok := builtins[c.cmd]; ok {
		if err := b(c); err != nil {
			return err
		}
		return nil
	}
	return runone(c)
}

func openRead(c *Command, r io.Reader, fd int) (io.Reader, error) {
	if c.fdmap[fd] != "" {
		f, err := os.Open(c.fdmap[fd])
		c.files[fd] = f
		return f, err
	}
	return r, nil
}

func openWrite(c *Command, w io.Writer, fd int) (io.Writer, error) {
	if c.fdmap[fd] != "" {
		f, err := os.OpenFile(c.fdmap[fd], os.O_CREATE|os.O_WRONLY, 0o666)
		c.files[fd] = f
		return f, err
	}
	return w, nil
}

func doArgs(cmds []*Command) {
	for _, c := range cmds {
		globargv := []string{}
		for _, v := range c.Args {
			if v.mod == "ENV" {
				globargv = append(globargv, os.Getenv(v.val))
			} else if globs, err := filepath.Glob(v.val); err == nil && len(globs) > 0 {
				globargv = append(globargv, globs...)
			} else {
				globargv = append(globargv, v.val)
			}
		}

		c.cmd = globargv[0]
		c.argv = globargv[1:]
	}
}

// There seems to be no harm in creating a Cmd struct
// even for builtins, so for now, we do.
// It will, however, do a path lookup, which we really don't need,
// and we may change it later.
func commands(cmds []*Command) error {
	for _, c := range cmds {
		c.Cmd = exec.Command(c.cmd, c.argv[:]...)
		// this is a Very Special Case related to a Go issue.
		// we're not able to unshare correctly in builtin.
		// Not sure of the issue but this hack will have to do until
		// we understand it. Barf.
		if c.cmd == "builtin" {
			builtinAttr(c)
		}
	}
	return nil
}

func command(c *Command) error {
	// for now, bg will just happen in background.
	if c.BG {
		go func() {
			if err := runit(c); err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
		}()
	}
	return runit(c)
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	b := bufio.NewReader(os.Stdin)
	tty()
	fmt.Printf("%% ")
	for {
		foreground()
		cmds, status, err := getCommand(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		doArgs(cmds)
		if err := commands(cmds); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if err := wire(cmds); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		for _, c := range cmds {
			if err := command(c); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				if c.Link == "||" {
					continue
				}
				// yes, not needed, but useful so you know
				// what goes on here.
				if c.Link == "&&" {
					break
				}
				break
			} else {
				if c.Link == "||" {
					break
				}
			}
		}
		if status == "EOF" {
			break
		}
		fmt.Printf("%% ")
	}
}
