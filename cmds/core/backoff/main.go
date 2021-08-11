// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run a command, repeatedly, until it succeeds or we are out of time
//
// Synopsis:
//	backoff [-t duration-string] command [args...]
//
// Description:
//	backoff will run the command until it succeeds or a timeout has passed.
//	The default timeout is 30s.
//	If no args are given, it will just return.
//
// Example:
//	$ backoff echo hi
//	hi
//	$
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var timeout = flag.String("t", "", "Timeout for command")

func runit(timeout string, c string, a ...string) error {
	ctx := context.Background()
	if len(timeout) != 0 {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}
		cx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		ctx = cx
	}
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	f := func() error {
		cmd := exec.Command(c, a...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return cmd.Run()
	}

	return backoff.Retry(f, b)
}

func main() {
	flag.Parse()

	var args []string
	a := flag.Args()
	if len(a) == 0 {
		return
	}
	if len(a) > 1 {
		args = a[1:]
	}
	if err := runit(*timeout, a[0], args[:]...); err != nil {
		log.Fatal(err)
	}
}
