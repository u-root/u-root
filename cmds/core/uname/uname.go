// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// Print build information about the kernel and machine.
//
// Synopsis:
//
//	uname [-asnrvmd]
//
// Options:
//
//	-a: print everything
//	-s: print the kernel name
//	-n: print the network node name
//	-r: print the kernel release
//	-v: print the kernel version
//	-m: print the machine hardware name
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	all       = flag.Bool("a", false, "print everything")
	kernel    = flag.Bool("s", false, "print the kernel name")
	node      = flag.Bool("n", false, "print the network node name")
	release   = flag.Bool("r", false, "print the kernel release")
	version   = flag.Bool("v", false, "print the kernel version")
	machine   = flag.Bool("m", false, "print the machine hardware name")
	processor = flag.Bool("p", false, "print the machine hardware name")
)

type params struct {
	kernel  bool
	node    bool
	release bool
	version bool
	machine bool
}

func handleFlags(u *unix.Utsname, p params) string {
	sysname, nodename := unix.ByteSliceToString(u.Sysname[:]), unix.ByteSliceToString(u.Nodename[:])
	release, version := unix.ByteSliceToString(u.Release[:]), unix.ByteSliceToString(u.Version[:])
	machine := unix.ByteSliceToString(u.Machine[:])
	info := make([]string, 0, 5)

	if p.kernel {
		info = append(info, sysname)
	}
	if p.node {
		info = append(info, nodename)
	}
	if p.release {
		info = append(info, release)
	}
	if p.version {
		info = append(info, version)
	}
	if p.machine {
		info = append(info, machine)
	}

	return strings.Join(info, " ")
}

func run(stdout io.Writer, p params) error {
	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		return err
	}
	info := handleFlags(&u, p)
	_, err := fmt.Fprintln(stdout, info)
	return err
}

func parseParams(all, kernel, node, release, version, machine, processor bool) params {
	p := params{
		kernel:  kernel || all,
		node:    node || all,
		release: release || all,
		version: version || all,
		machine: machine || processor || all,
	}

	if !p.kernel && !p.node && !p.release && !p.version && !p.machine {
		p.kernel = true
	}

	return p
}

func main() {
	flag.Parse()
	p := parseParams(*all, *kernel, *node, *release, *version, *machine, *processor)
	if err := run(os.Stdout, p); err != nil {
		log.Fatal(err)
	}
}
