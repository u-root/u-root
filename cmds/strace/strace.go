// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// strace is a simple single-process tracer.
// It starts the comand and lets the strace.Run() do all the work.
//
// Synopsis:
//     strace <command> [args...]
//
// Description:
//	trace a single process given a command name.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/strace"
)

var (
	cmdUsage = "Usage: strace <command> [args...]"
	debug    = flag.BoolP("debug", "d", false, "enable debug printing")
)

func usage() {
	log.Fatalf(cmdUsage)
}

func main() {
	flag.Parse()

	if *debug {
		strace.Debug = log.Printf
	}
	a := flag.Args()
	if len(a) < 1 {
		usage()
	}

	c := exec.Command(a[0], a[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr

	t, err := strace.New()
	if err != nil {
		log.Fatal(err)
	}

	go t.RunTracerFromCmd(c)

	for r := range t.Records {
		fmt.Printf("%s\n", r.String())
	}
}
