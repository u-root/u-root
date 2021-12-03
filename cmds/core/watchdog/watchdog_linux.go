// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// watchdog interacts with /dev/watchdog.
//
// Synopsis:
//     watchdog [--dev=DEV] keepalive
//         Pet the watchdog. This resets the time left back to the timeout.
//     watchdog [--dev=DEV] set[pre]timeout DURATION
//         Set the watchdog timeout or pretimeout
//     watchdog [--dev=DEV] get[pre]timeout
//         Print the watchdog timeout or pretimeout
//     watchdog [--dev=DEV] gettimeleft
//         Print the amount of time left.
//
// Options:
//     --dev DEV: Device (default /dev/watchdog)
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/watchdog"
)

var dev = flag.String("dev", "/dev/watchdog", "device")

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
	os.Exit(1)
}

func runCommand() error {
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	wd, err := watchdog.Open(*dev)
	if err != nil {
		return err
	}
	defer func() {
		if err := wd.Close(); err != nil {
			log.Printf("Failed to close watchdog: %v", err)
		}
	}()

	switch flag.Arg(0) {
	case "keepalive":
		if flag.NArg() != 1 {
			usage()
		}
		if err := wd.KeepAlive(); err != nil {
			return err
		}
	case "settimeout":
		if flag.NArg() != 2 {
			usage()
		}
		d, err := time.ParseDuration(flag.Arg(1))
		if err != nil {
			return err
		}
		if err := wd.SetTimeout(d); err != nil {
			return err
		}
	case "setpretimeout":
		if flag.NArg() != 2 {
			usage()
		}
		d, err := time.ParseDuration(flag.Arg(1))
		if err != nil {
			return err
		}
		if err := wd.SetPreTimeout(d); err != nil {
			return err
		}
	case "gettimeout":
		if flag.NArg() != 1 {
			usage()
		}
		i, err := wd.Timeout()
		if err != nil {
			return err
		}
		fmt.Println(i)
	case "getpretimeout":
		if flag.NArg() != 1 {
			usage()
		}
		i, err := wd.PreTimeout()
		if err != nil {
			return err
		}
		fmt.Println(i)
	case "gettimeleft":
		if flag.NArg() != 1 {
			usage()
		}
		i, err := wd.TimeLeft()
		if err != nil {
			return err
		}
		fmt.Println(i)
	default:
		return fmt.Errorf("unrecognized command: %q", flag.Arg(0))
	}
	return nil
}

func main() {
	if err := runCommand(); err != nil {
		log.Fatal(err)
	}
}
