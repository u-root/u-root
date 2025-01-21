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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/watchdog"
	"github.com/u-root/u-root/pkg/watchdogd"
)

func usage() {
	fmt.Print(`watchdogd run [--dev DEV] [--timeout N] [--pre_timeout N] [--keep_alive N] [--monitors STRING]
	Run the watchdogd daemon in a child process (does not daemonize).
watchdogd stop
	Send a signal to arm the running watchdogd.
watchdogd continue
	Send a signal to disarm the running watchdogd.
watchdogd arm
	Send a signal to arm the running watchdogd.
watchdogd disarm
	Send a signal to disarm the running watchdogd.
`)
	os.Exit(1)
}

func runCommand() error {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
	}
	cmd, args := args[0], args[1:]

	switch cmd {
	case "run":
		fs := flag.NewFlagSet("run", flag.PanicOnError)
		var (
			dev        = fs.String("dev", watchdog.Dev, "device")
			timeout    = fs.Duration("timeout", -1, "duration before timing out")
			preTimeout = fs.Duration("pre_timeout", -1, "duration for pretimeout")
			keepAlive  = fs.Duration("keep_alive", 5*time.Second, "duration between issuing keepalive")
			monitors   = fs.String("monitors", "", "comma seperated list of monitors, ex: oops")
			uds        = fs.String("uds", "/tmp/watchdogd", "unix domain socket path for the daemon")
		)
		fs.Parse(args)
		if fs.NArg() != 0 {
			usage()
		}
		if *timeout == -1 {
			timeout = nil
		}
		if *preTimeout == -1 {
			preTimeout = nil
		}

		monitorFuncs := []func() error{}
		for _, m := range strings.Split(*monitors, ",") {
			if m == "oops" {
				monitorFuncs = append(monitorFuncs, watchdogd.MonitorOops)
			} else {
				return fmt.Errorf("unrecognized monitor: %v", m)
			}
		}

		return watchdogd.Run(context.Background(), &watchdogd.DaemonOpts{
			Dev:        *dev,
			Timeout:    timeout,
			PreTimeout: preTimeout,
			KeepAlive:  *keepAlive,
			Monitors:   monitorFuncs,
			UDS:        *uds,
		})
	default:
		if len(args) != 0 {
			usage()
		}
		d, err := watchdogd.NewClient()
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
