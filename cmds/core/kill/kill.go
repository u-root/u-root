// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// Kill kills processes.
//
// Synopsis:
//
//	kill -l
//	kill [<-s | --signal | -> <isgname|signum>] pid [pid...]
//
// Options:
//
//	-l:                       list the signal names
//	-name, --signal name, -s: name is the message to send. On some systems
//	                          this is a string, on others a number. It is
//	                          optional and an OS-dependent value will be
//	                          used if it is not set. pid is a list of at
//	                          least one pid.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const eUsage = "Usage: kill -l | kill [<-s | --signal | -> <signame|signum>] pid [pid...]"

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s\n", eUsage)
}

func killProcess(w io.Writer, args ...string) error {
	if len(args) < 2 {
		usage(w)
		return nil
	}
	op := args[1]
	pids := args[2:]
	if op[0] != '-' {
		op = defaultSignal
		pids = args[1:]
	}
	// sadly, we can not use flag. Well, we could,
	// it would be pretty cheap to just put in every
	// possible flag as a switch, but it just gets
	// messy.
	//
	// kill is one of the old school commands (1971)
	// and arg processing was not really standard back then.
	// Also, note, the -l has no meaning on Plan 9 or Harvey
	// since signals on those systems are arbitrary strings.

	if op[0:2] == "-l" {
		if len(args) > 2 {
			usage(w)
			return nil
		}
		fmt.Fprintf(w, "%s\n", siglist())
		return nil
	}

	// N.B. Be careful if you want to change this. It has to continue to work if
	// the signal is an arbitrary string. We take the path that the signal
	// has to start with a -, or might be preceded by -s or --string.

	if op == "-s" || op == "--signal" {
		if len(args) < 3 {
			usage(w)
			return nil
		}
		op = args[2]
		pids = args[3:]
	} else {
		op = op[1:]
	}

	s, ok := signums[op]
	if !ok {
		return fmt.Errorf("%v is not a valid signal", op)
	}

	if len(pids) < 1 {
		usage(w)
		return nil
	}
	if err := kill(s, pids...); err != nil {
		return fmt.Errorf("some processes could not be killed: %w", err)
	}
	return nil
}

func main() {
	if err := killProcess(os.Stdout, os.Args...); err != nil {
		log.Fatal(err)
	}
}
