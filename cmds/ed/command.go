// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Options:
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Command(f Editor, c string, startLine, endLine int) error {
	var err error
	if len(c) == 0 {
		_, err = f.Write(os.Stdout, f.Dot(), f.Dot())
		f.Move(endLine + 1)
		return err
	}

	a := c[1:]
	debug("Process %c, args %v", c[0], a)
	switch c[0] {
	case 'q', 'e':
		if f.IsDirty() {
			f.Dirty(false)
			return fmt.Errorf("f was dirty, no longer is, try again")
		}
	}
	switch c[0] {
	case 'd':
		f.Replace([]byte{}, startLine, endLine)
	case 'q':
		os.Exit(1)
	case 'e':
		startLine, endLine = f.Range()
		startLine = 0
		fallthrough
	case 'r':
		fname := strings.TrimLeft(a, " \t")
		debug("read %v @ %v, %v", f, startLine, endLine)
		var r io.Reader
		r, err = os.Open(fname)
		debug("%v: r is %v, err %v", fname, r, err)
		if err == nil {
			_, err = f.Read(r, startLine, endLine)
		}
	case 's':
		o := strings.SplitN(a[1:], a[0:1], 3)
		if o[1] == "" {
			o[1] = "" //f.pat
		}
		debug("after split o is %v", o)
		err = f.Sub(o[0], o[1], o[2], startLine, endLine)
	case 'w':
		fname := strings.TrimLeft(a, " \t")
		debug("NOT WRITINGT TO %v", fname)
		_, err = f.Write(os.Stdout, startLine, endLine)
	case 'p':
		fmt.Printf("what the shit")
		_, err = f.Print(os.Stdout, startLine, endLine)
		fmt.Printf("shiw %v", err)
	default:
		err = fmt.Errorf("%c: unknown command", c[0])
	}
	return err
}

func DoCommand(f Editor, l string) error {
	var err error
	var startLine, endLine int
	startLine = f.Dot()
	if len(l) == 0 {
		_, err = f.Write(os.Stdout, f.Dot(), f.Dot())
		f.Move(f.Dot() + 1)
		return err
	}
	switch {
	case l[0] == '.':
		debug(".\n")
		startLine = f.Dot()
		l = l[1:]
	case l[0] == '$':
		debug("$\n")
		_, endLine = f.Range()
		f.Move(endLine)
		startLine = f.Dot()
		l = l[1:]
	case startsearch.FindString(l) != "":
		pat := startsearch.FindString(l)
		l = l[len(pat):]
		debug("/\n")
		fail("Pattern search: not yet")
	case num.FindString(l) != "":
		debug("num\n")
		n := num.FindString(l)
		if startLine, err = strconv.Atoi(n); err != nil {
			return err
		}
		f.Move(startLine)
		l = l[len(n):]
	}
	debug("cmd before endsearch is %v", l)
	endLine = f.Dot()
	if len(l) > 0 && l[0] == ',' {
		l = l[1:]
		if len(l) < 1 {
			return fmt.Errorf("line ended at a ,?")
		}
		switch {
		case l[0] == '.':
			debug(".\n")
			endLine = f.Dot()
			l = l[1:]
		case l[0] == '$':
			debug("$\n")
			_, endLine = f.Range()
			l = l[1:]
		case startsearch.FindString(l) != "":
			debug("/\n")
			fail("Pattern search: not yet")
		case num.FindString(l) != "":
			debug("num\n")
			n := num.FindString(l)
			if endLine, err = strconv.Atoi(n); err != nil {
				return err
			}
			l = l[len(n):]
		}
	}
	endLine++
	debug("l before call to f.Command() is %v", l)
	err = Command(f, l, startLine, endLine)
	debug("Comand is done: f is %v", f)
	return err
}
