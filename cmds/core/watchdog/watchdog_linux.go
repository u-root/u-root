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
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/watchdog"
)

var dev = flag.String("dev", "/dev/watchdog", "device")
var errUsage = errors.New("usage error")

func usage() {
	flag.Usage()
	fmt.Print(`watchdog keepalive
	Pet the watchdog. This resets the time left back to the timeout.
watchdog set[pre]timeout DURATION
	Set the watchdog timeout or pretimeout.
watchdog get[pre]timeout
	Print the watchdog timeout or pretimeout.
watchdog gettimeleft
	Print the amount of time left.
`)
}

func runCommand(dev string, args ...string) error {
	if len(args) < 1 {
		return errUsage
	}

	wd, err := watchdog.Open(dev)
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

func main() {
	flag.Parse()
	if err := runCommand(*dev, flag.Args()...); err != nil {
		if errors.Is(err, errUsage) {
			usage()
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
