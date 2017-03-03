// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Kill kills processes.
//
// Synopsis:
//     kill -l
//     kill [<-s | --signal | -> <isgname|signum>] pid [pid...]
//
// Options:
//     -l:                       list the signal names
//     -name, --signal name, -s: name is the message to send. On some systems
//                               this is a string, on others a number. It is
//                               optional and an OS-dependent value will be
//                               used if it is not set. pid is a list of at
//                               least one pid.
package main

import (
	"fmt"
	"os"
)

const eUsage = "Usage: kill -l | kill [<-s | --signal | -> <signame|signum>] pid [pid...]"

func usage() {
	die(eUsage)
}

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	op := os.Args[1]
	pids := os.Args[2:]
	if op[0] != '-' {
		op = defaultSignal
		pids = os.Args[1:]
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
		if len(os.Args) > 2 {
			usage()
		}
		fmt.Print(siglist())
		os.Exit(1)
	}

	// N.B. Be careful if you want to change this. It has to continue to work if
	// the signal is an arbitrary string. We take the path that the signal
	// has to start with a -, or might be preceded by -s or --string.

	if op == "-s" || op == "--signal" {
		if len(os.Args) < 3 {
			usage()
		}
		op = os.Args[2]
		pids = os.Args[3:]
	} else {
		op = op[1:]
	}

	s, ok := signums[op]
	if !ok {
		die("%v is not a valid signal", op)
	}

	if len(pids) < 1 {
		usage()
	}
	if err := kill(s, pids...); err != nil {
		die("Some processes could not be killed: %v", err)
	}
}
