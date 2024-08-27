// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// watchdog interacts with /dev/watchdog.
//
// Synopsis:
//
//	watchdog [--dev=DEV] keepalive
//	    Pet the watchdog. This resets the time left back to the timeout.
//	watchdog [--dev=DEV] set[pre]timeout DURATION
//	    Set the watchdog timeout or pretimeout
//	watchdog [--dev=DEV] get[pre]timeout
//	    Print the watchdog timeout or pretimeout
//	watchdog [--dev=DEV] gettimeleft
//	    Print the amount of time left.
//
// Options:
//
//	--dev DEV: Device (default /dev/watchdog)
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/watchdog"
)

var errUsage = errors.New("usage error")

type cmd struct {
	dev string
}

func (c cmd) run(args []string) error {
	if len(args) < 1 {
		return errUsage
	}

	wd, err := watchdog.Open(c.dev)
	if err != nil {
		return err
	}
	defer func() {
		if err := wd.Close(); err != nil {
			log.Printf("Failed to close watchdog: %v", err)
		}
	}()

	switch args[0] {
	case "keepalive":
		if len(args) != 1 {
			return errUsage
		}
		if err := wd.KeepAlive(); err != nil {
			return err
		}
	case "settimeout":
		if len(args) != 2 {
			return errUsage
		}
		d, err := time.ParseDuration(args[1])
		if err != nil {
			return err
		}
		if err := wd.SetTimeout(d); err != nil {
			return err
		}
	case "setpretimeout":
		if len(args) != 2 {
			return errUsage
		}
		d, err := time.ParseDuration(args[1])
		if err != nil {
			return err
		}
		if err := wd.SetPreTimeout(d); err != nil {
			return err
		}
	case "gettimeout":
		if len(args) != 1 {
			return errUsage
		}
		i, err := wd.Timeout()
		if err != nil {
			return err
		}
		fmt.Println(i)
	case "getpretimeout":
		if len(args) != 1 {
			return errUsage
		}
		i, err := wd.PreTimeout()
		if err != nil {
			return err
		}
		fmt.Println(i)
	case "gettimeleft":
		if len(args) != 1 {
			return errUsage
		}
		i, err := wd.TimeLeft()
		if err != nil {
			return err
		}
		fmt.Println(i)
	default:
		return fmt.Errorf("unrecognized command: %q", args[0])
	}
	return nil
}

func run(args []string) error {
	var c cmd
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&c.dev, "dev", "/dev/watchdog", "watchdog device")

	f.Usage = func() {
		fmt.Fprintf(f.Output(), `Usage:
watchdog [--dev=DEV] keepalive
	Pet the watchdog. This resets the time left back to the timeout.
watchdog [--dev=DEV] set[pre]timeout DURATION
	Set the watchdog timeout or pretimeout.
watchdog [--dev=DEV] get[pre]timeout
	Print the watchdog timeout or pretimeout.
watchdog [--dev=DEV] gettimeleft
	Print the amount of time left.

Options:
`)
		f.PrintDefaults()
	}

	f.Parse(args[1:])
	err := c.run(f.Args())
	if err != nil {
		if errors.Is(err, errUsage) {
			f.Usage()
		}
		return err
	}
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
