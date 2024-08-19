// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && ((linux && arm64) || (linux && amd64) || (linux && riscv64))

// syscallfilter runs a command with a possibly empty set of filters:
//
// Synopsis:
// syscallfilter [event]* [--] command [arguments]
//
// Description:
//
// The command name at least is mandatory. If no filters are specified,
// then the command runs as normal.
// Filters can be used to manage events.
// Events are specified as triples: eventRE,action,value
// The event is a regular expression, matching a system call name, at entry or exit,
// e.g. E.*read.* will match all variants of read at entry; or an strace
// event, e.g. NewChild will match strace events concerning new process creation.
// For more details, see the syscallfilter package.
//
// Should we wish to bar date from getting the time of day, for example:
// syscallfilter E.*time.*,error,-1
//
// And the result is:
// ./syscallfilter E.*timeof.*,error,-1 -- date
// 2022/01/13 11:50:47 Filtering ["date"]: -1
//
// You can deny echo from fulfilling its destiny:
// syscallfilter E.*write.*,error,-1 -- echo hi
// 2022/01/13 11:52:03 Filtering ["echo" "hi"]: -1
//
// You can enforce limits on population growth:
// syscallfilter NewChild,error,-1 -- bash -c date
// 2022/01/13 11:53:21 Filtering ["bash" "-c" "date"]: -1
//
// Note that the syscallfilter command can not currently tell if the error was a real error or
// filtered error. There is a case to be made that it should not be possible to tell,
// so, for now, we do not make it possible to distinguish real errors and fake errors.
package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/syscallfilter"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var logactions = flag.Bool("l", false, "Log actions output from the filter")

const cmdUsage = "Usage: syscallfilter [-l] [action... --] command [args]"

func main() {
	// TODO: fill this in from arguments.
	flag.Usage = util.Usage(flag.Usage, cmdUsage)
	flag.Parse()

	// By default, there are no actions, and this becomes just "run a program"
	var args, events []string
	for _, a := range flag.Args() {
		if a == "--" {
			events = args
			args = []string{}
			continue
		}
		args = append(args, a)
	}

	if len(args) == 0 {
		flag.Usage()
		log.Fatal(1)
	}
	c := syscallfilter.Command(args[0], args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := c.AddActions(events...); err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	if *logactions {
		c.Log = &b
	}
	// There's a bit of question here about printing errors.
	// Did the command get an error because it got an error,
	// for other reasons, or because we forced the error?
	// Anyway, for now, we'll print the error and hope the
	// user knows that they got an error because they asked
	// for an error.
	if err := c.Run(); err != nil {
		log.Printf("Error filtering %q with events %q: %v", c, events, err)
	}
	if *logactions {
		log.Printf("Actions Log:\n%v", b.String())
	}
}
