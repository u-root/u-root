// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delay for the specified amount of time.
//
// Synopsis:
//
//	sleep DURATION
//
// Description:
//
//	If no units are given, the duration is assumed to be measured in
//	seconds, otherwise any format parsed by Go's `time.ParseDuration` is
//	accepted.
//
// Examples:
//
//	sleep 2.5
//	sleep 300ms
//	sleep 2h45m
//
// Bugs:
//
//	When sleep is first run, it must be compiled from source which creates a
//	delay significantly longer than anticipated.
package main

import (
	"errors"
	"flag"
	"log"
	"time"
)

var errDuration = errors.New("invalid duration")

func parseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		d, err = time.ParseDuration(s + "s")
	}
	if err != nil || d < 0 {
		return time.Duration(0), errDuration
	}
	return d, nil
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Incorrect number of arguments")
	}

	d, err := parseDuration(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(d)
}
