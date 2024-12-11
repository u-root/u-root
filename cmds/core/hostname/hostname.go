// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9 || linux

// hostname prints or changes the system's hostname.
//
// Synopsis:
//
//	hostname [HOSTNAME]
//
// Author:
//
//	Beletti <rhiguita@gmail.com>
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func run(stdout io.Writer, args []string) error {
	switch len(args) {
	case 2:
		if err := Sethostname(args[1]); err != nil {
			return fmt.Errorf("could not set hostname: %w", err)
		}
		return nil
	case 1:
		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("could not obtain hostname: %w", err)
		}
		_, err = fmt.Fprintln(stdout, hostname)
		return err
	default:
		return fmt.Errorf("usage: hostname [HOSTNAME]")
	}
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
