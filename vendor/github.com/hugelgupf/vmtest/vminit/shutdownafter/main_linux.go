// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command shutdownafter runs a command given in args and shuts down.
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func run() error {
	args := flag.Args()
	if len(args) == 0 {
		return nil
	}
	c := exec.Command(args[0], args[1:]...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Printf("Failed: %v", err)
	}

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to shutdown: %v", err)
	}
}
