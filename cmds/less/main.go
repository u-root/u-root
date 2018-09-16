// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// less pages through files
//
// Synopsis:
//     less [OPTIONS] FILES...
//
// Options:
//     -profile FILE: Save profile in this file
//     -tabstop NUMBER: Number of spaces per tab
//
// Keybindings:
//     Control:
//
//     * q: Quit
//
//     Scrolling:
//
//     * j: Scroll down
//     * k: Scroll up
//     * g: Scroll to top
//     * G: Scroll to bottom
//     * Pgdn: Scroll down one screen full
//     * Pgup: Scroll up one screen full
//     * ^D: Scroll down one half screen full
//     * ^U: Scroll up one half screen full
//
//     Searching:
//
//     * /: Enter search regex (re2 syntax). Press enter to search.
//     * n: Jump down to next search result
//     * N: Jump up to previous search result
//
// Author:
//     Michael Pratt (github.com/prattmic)
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"syscall"

	"github.com/nsf/termbox-go"
)

var profile = flag.String("profile", "", "Save profile in this file")
var tabStop = flag.Int("tabstop", 8, "Number of spaces per tab")

func mmapFile(f *os.File) ([]byte, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return syscall.Mmap(int(f.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
}

func main() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: %s filename\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	name := flag.Arg(0)
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s: %v\n", name, err)
		os.Exit(1)
	}

	m, err := mmapFile(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to mmap file: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Munmap(m)

	err = termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init: %v\n", err)
		os.Exit(1)
	}
	defer termbox.Close()

	if *profile != "" {
		p, err := os.Create(*profile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create profile: %v\n", err)
			os.Exit(1)
		}
		defer p.Close()
		pprof.StartCPUProfile(p)
		defer pprof.StopCPUProfile()
	}

	l := NewLesser(bytes.NewReader(m), *tabStop)
	l.Run()
}
