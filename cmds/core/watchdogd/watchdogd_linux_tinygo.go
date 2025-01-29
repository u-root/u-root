// Copyright 2021-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && tinygo

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

func run(args []string) error {
	if len(args) == 0 {
		watchdog.Usage()
		return watchdogd.ErrNoCommandSpecified
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "run":
		daemonOpts := &watchdogd.DaemonOpts{}
		fs := daemonOpts.InitFlags()
		monitor := fs.String("monitors", "oops", "comma separated list of monitors")
		fs.Parse(args)

		if fs.NArg() != 0 {
			watchdog.Usage()
			return watchdogd.ErrMissingArgument
		}

		daemonOpts.Monitors = []func() error{}
		for _, m := range strings.Split(*monitor, ",") {
			if m == "oops" {
				daemonOpts.Monitors = append(daemonOpts.Monitors, watchdogd.MonitorOops)
			} else {
				return fmt.Errorf("%w: %v", watchdogd.ErrInvalidMonitor, m)
			}
		}

		return watchdogd.Run(context.Background(), daemonOpts)

	case "pid":
		if len(args) != 0 {
			watchdog.Usage()
			return watchdogd.ErrTooManyArguments
		}

		daemon, err := watchdogd.New()
		if err != nil {
			return fmt.Errorf("could not find watchdog daemon: %w", err)
		}

		fmt.Println(daemon.Pid)

	default:
		if len(args) != 0 {
			watchdog.Usage()
			return watchdogd.ErrTooManyArguments
		}

		daemon, err := watchdogd.New()
		if err != nil {
			return fmt.Errorf("could not dial watchdog daemon: %w", err)
		}
		f, ok := map[string]func() error{
			"stop":     daemon.Stop,
			"continue": daemon.Continue,
			"arm":      daemon.Arm,
			"disarm":   daemon.Disarm,
		}[cmd]
		if !ok {
			return watchdogd.ErrNoCommandSpecified
		}
		return f()
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
