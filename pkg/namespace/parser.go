// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package namespace parses name space description files
// https://plan9.io/magic/man2html/6/namespace
package namespace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// Parse takes a namespace file and returns a collection
// of operations that build a name space in Plan 9.
//
// Think oci runtime spec, but don't think too much cause the JSON
// will make your head hurt.
//
// http://man.cat-v.org/plan_9/1/ns
func Parse(r io.Reader) (File, error) {
	scanner := bufio.NewScanner(r)

	cmds := []Modifier{}
	for scanner.Scan() {
		buf := scanner.Bytes()
		if len(buf) <= 0 {
			continue
		}
		r := buf[0]
		// Blank lines and lines with # as the first non–space character are ignored.
		if r == '#' || r == ' ' {
			continue
		}
		cmd, err := ParseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return cmds, nil
}

// ParseFlags parses flags as mount, bind would. Supports -a, -b, -c,
// -q, and -C, although the -q flag is not yet implemented. Other flags
// are silently ignored.
// Not using the os.Flag package here as that would break the
// name space description files.
// https://9p.io/magic/man2html/6/namespace
func ParseFlags(args []string) (mountflag, []string) {
	flag := REPL
	for i, arg := range args {
		// these args are passed trough strings.Fields which doesn't return empty strings
		// so this is ok.
		if arg[0] == '-' {
			args = append(args[:i], args[i+1:]...)
			for _, r := range arg {
				switch r {
				case 'a':
					flag |= AFTER
				case 'b':
					flag |= BEFORE
				case 'c':
					flag |= CREATE
				case 'q':
				// todo(sevki): support quiet flag
				case 'C':
					flag |= CACHE
				default:
				}
			}
		}
	}
	return flag, args
}

// ParseArgs parses os.Args (not flag.Args()!) to return a
// namespace.Modifier which is appropriate for the given
// command. The command is taken from args[0]; the following
// commands are supported: bind, mount, unmount, clear, cd, ., import.
//
// After extracting the command, ParseArgs will call ParseFlags to extract
// appropriate bind/mount flags, if any.
func ParseArgs(args []string) (Modifier, error) {
	arg := path.Base(args[0])
	args = args[1:]
	trap := syzcall(0)

	c := cmd{
		syscall: trap,
		flag:    REPL,
		args:    args,
	}
	switch arg {
	case "bind":
		c.syscall = BIND
	case "mount":
		c.syscall = MOUNT
	case "unmount":
		c.syscall = UNMOUNT
	case "clear":
		c.syscall = RFORK
	case "cd":
		c.syscall = CHDIR
	case ".":
		c.syscall = INCLUDE
	case "import":
		c.syscall = IMPORT
	default:
		return nil, fmt.Errorf("%q is an unknown operation", arg)
	}

	c.flag, c.args = ParseFlags(args)

	return c, nil
}

// ParseLine breaks a string into arguments and passes them to ParseArgs.
func ParseLine(line string) (Modifier, error) {
	return ParseArgs(strings.Fields(line))
}
