// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && linux && (amd64 || riscv64 || arm64)
// +build !tinygo
// +build linux
// +build amd64 riscv64 arm64

// traceopen traces a process trees file opens.
// It is intended to make building small containers easier, so that,
// e.g., gigantic containers (e.g. the 25G Nemo ML container)
// can be cut down.
//
// Synopsis:
//
//	traceopen <command> [args...]
//
// Description:
//
//	trace all file opens for a tree of processes.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/strace"
	"golang.org/x/sys/unix"
)

var errUsage = errors.New("usage: traceopen [-o <outputfile>] <command> [args...]")

type params struct {
	output string
}

// PrintTraces prints only trace entry events, and only
// for a subset of system calls, and only the file name
// and open mode.
func PrintFiles(w io.Writer) strace.EventCallback {
	return func(t strace.Task, record *strace.TraceRecord) error {
		if record.Event != strace.SyscallEnter {
			return nil
		}
		var n strace.Addr
		var mode uint
		var cwd string
		var err error

		e := record.Syscall
		switch e.Sysno {
		case syscall.SYS_OPEN:
			n = e.Args[0].Pointer()
			mode = e.Args[1].ModeT()
		case syscall.SYS_OPENAT:
			if e.Args[0].Int() != unix.AT_FDCWD {
				at := e.Args[0].Pointer()
				if cwd, err = strace.ReadString(t, at, unix.PathMax); err != nil {
					return fmt.Errorf("reading dir from pointer %#x:%w", at, err)
				}
			}

			n = e.Args[1].Pointer()
			mode = e.Args[2].ModeT()
		default:
			return nil
		}

		name, err := strace.ReadString(t, n, unix.PathMax)
		if err != nil {
			return fmt.Errorf("reading name from pointer %#x:%w", n, err)
		}

		fmt.Fprintf(w, "%q %#o\n", filepath.Join(cwd, name), mode)
		return nil
	}
}

func run(stdin io.Reader, stdout, stderr io.Writer, p params, args ...string) error {
	if len(args) < 1 {
		return errUsage
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = stdin, stdout, stderr

	if p.output != "" {
		f, err := os.Create(p.output)
		if err != nil {
			return fmt.Errorf("creating output file: %s: %w", p.output, err)
		}
		defer f.Close()
		c.Stderr = f
	}

	return strace.Trace(c, PrintFiles(stderr))
}

func main() {
	output := flag.String("o", "", "write output to file (if empty, stdout)")
	flag.Parse()

	if err := run(os.Stdin, os.Stdout, os.Stderr, params{output: *output}, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
