// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Rush is an interactive shell similar to sh.
//
// Description:
//     Prompt is '% '.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

type builtin func(c *Command) error

// TODO: probably have one builtin map and use it for both types?
var (
	urpath   = "/go/bin:/ubin:/buildbin:/bbin:/bin:/usr/local/bin:"
	builtins = make(map[string]builtin)
	// Some builtins really want to be forked off, esp. in the busybox case.
	forkBuiltins = make(map[string]builtin)
	// the environment dir is INTENDED to be per-user and bound in
	// a private name space at /env.
	envDir = "/env"
)

func addBuiltIn(name string, f builtin) error {
	if _, ok := builtins[name]; ok {
		return fmt.Errorf("%v already a builtin", name)
	}
	builtins[name] = f
	return nil
}

func addForkBuiltIn(name string, f builtin) error {
	if _, ok := builtins[name]; ok {
		return fmt.Errorf("%v already a forkBuiltin", name)
	}
	forkBuiltins[name] = f
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
		if c.link != "|" {
			if c.Stdout, err = openWrite(c, os.Stdout, 1); err != nil {
				return err
			}
		}
		if c.Stderr, err = openWrite(c, os.Stderr, 2); err != nil {
			return err
		}
		// The validation is such that "|" is not set on the last one.
		// Also, there won't be redirects and "|" inappropriately.
		if c.link != "|" {
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
			io.Copy(w, r)
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
	} else {
		c.Cmd.SysProcAttr = &syscall.SysProcAttr{}
		if c.bg {
			c.Cmd.SysProcAttr.Setpgid = true
		} else {
			c.Cmd.SysProcAttr.Foreground = true
			c.Cmd.SysProcAttr.Ctty = int(ttyf.Fd())
		}
		if err := c.Start(); err != nil {
			return fmt.Errorf("%v: Path %v", err, os.Getenv("PATH"))
		}
		if err := c.Wait(); err != nil {
			return fmt.Errorf("wait: %v", err)
		}
	}
	return nil
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
		f, err := os.Create(c.fdmap[fd])
		c.files[fd] = f
		return f, err
	}
	return w, nil
}

func doArgs(cmds []*Command) error {
	for _, c := range cmds {
		globargv := []string{}
		for _, v := range c.args {
			if v.mod == "ENV" {
				e := v.val
				if !path.IsAbs(v.val) {
					e = filepath.Join(envDir, e)
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
	}
	return nil
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
			c.Cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
		}
	}
	return nil
}
func command(c *Command) error {
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
	b := bufio.NewReader(os.Stdin)

	// we use path.Base in case they type something like ./cmd
	if f, ok := forkBuiltins[path.Base(os.Args[0])]; ok {
		if err := f(&Command{cmd: os.Args[0], Cmd: &exec.Cmd{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}, argv: os.Args[1:]}); err != nil {
			log.Fatalf("%v", err)
		}
		os.Exit(0)
	}

	if len(os.Args) != 1 {
		fmt.Println("no scripts/args yet")
		os.Exit(1)
	}

	tty()
	fmt.Printf("%% ")
	for {
		foreground()
		cmds, status, err := getCommand(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		if err := doArgs(cmds); err != nil {
			fmt.Fprintf(os.Stderr, "args problem: %v\n", err)
			continue
		}
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
				if c.link == "||" {
					continue
				}
				// yes, not needed, but useful so you know
				// what goes on here.
				if c.link == "&&" {
					break
				}
				break
			} else {
				if c.link == "||" {
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
