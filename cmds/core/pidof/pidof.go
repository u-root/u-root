// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var errNotFound = errors.New("pid name not found")

type proc struct {
	comm string
	pid  string
}

func run(stdout io.Writer, procPath string, single, quiet bool, args []string) error {
	procs, err := collect(procPath)
	if err != nil {
		return err
	}

	var pids []string

	for _, proc := range procs {
		for _, arg := range args {
			if proc.comm == arg {
				if single {
					if !quiet {
						fmt.Fprintln(stdout, proc.pid)
					}
					return nil
				}
				pids = append(pids, proc.pid)
			}
		}
	}

	if len(pids) == 0 {
		return errNotFound
	}

	if !quiet {
		fmt.Fprintln(stdout, strings.Join(pids, " "))
	}
	return nil
}

func main() {
	single := flag.Bool("s", false, "single shot - this instructs the program to only return one pid")
	quiet := flag.Bool("q", false, "do not display matched PIDs to standard out")
	flag.Parse()
	if err := run(os.Stdout, procPath, *single, *quiet, flag.Args()); err != nil {
		if errors.Is(err, errNotFound) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
