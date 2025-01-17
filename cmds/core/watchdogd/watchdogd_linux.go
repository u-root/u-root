// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// watchdogd is a background daemon for petting the watchdog.
//
// Synopsis:
//
//	watchdogd run [OPTIONS]
//	    Run the watchdogd in a child process (does not daemonize).
//	watchtdog pid [tinygo only]
//	    Print the PID of the running watchdog.
//	watchdogd stop
//	    Send a signal to arm the running watchdog.
//	watchdogd continue
//	    Send a signal to disarm the running watchdog.
//	watchdogd arm
//	    Send a signal to arm the running watchdog.
//	watchdogd disarm
//	    Send a signal to disarm the running watchdog.
//
// Options:
//
//	--dev DEV: Device (default /dev/watchdog)
//	--timeout: Duration before timing out (default -1)
//	--pre_timeout: Duration for pretimeout (default -1)
//	--keep_alive: Duration between issuing keepalive (default 10)
//	--monitors: comma separated list of monitors, ex: oops
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/watchdog"
	"github.com/u-root/u-root/pkg/watchdogd"
)

func runCommand() error {
	args := os.Args[1:]
	if len(args) == 0 {
		watchdog.Usage()
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]

	switch cmd {
	case "run":
		daemonOpts := &watchdogd.DaemonOpts{}
		fs := daemonOpts.InitFlags()
		monitor := fs.String("monitor", "oops", "comma separated list of monitors")
		fs.Parse(args)

		if fs.NArg() != 0 {
			watchdog.Usage()
		}

		daemonOpts.Monitors = []func() error{}
		for _, m := range strings.Split(*monitor, ",") {
			if m == "oops" {
				daemonOpts.Monitors = append(daemonOpts.Monitors, watchdogd.MonitorOops)
			} else {
				return fmt.Errorf("unrecognized monitor: %v", m)
			}
		}

		return watchdogd.Run(context.Background(), daemonOpts)

	default:
		if len(args) != 0 {
			watchdog.Usage()
		}
		d, err := watchdogd.New()
		if err != nil {
			return fmt.Errorf("could not dial watchdog daemon: %w", err)
		}
		f, ok := map[string]func() error{
			"stop":     d.Stop,
			"continue": d.Continue,
			"arm":      d.Arm,
			"disarm":   d.Disarm,
		}[cmd]
		if !ok {
			return fmt.Errorf("unrecognized command %q", cmd)
		}
		return f()
	}
}

func main() {
	if err := runCommand(); err != nil {
		log.Fatal(err)
	}
}
