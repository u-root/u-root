// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !windows

// ntpdate uses NTP to adjust the system clock.
//
// Synopsis:
//
//	ntpdate [--config=/etc/ntp.conf] [--rtc] [--verbose] [server ...]
//
// Description:
//
//	ntpdate queries NTP server(s) for time and update system time.
//	If --rtc is specified, it updates the hardware clock as well.
//	Servers to query are obtained from /etc/ntp.conf and/or the command line.
//	By default --config is set to /etc/ntp.conf, config lookup can be disabled
//	by setting --config to an empty string.
//	If servers are specified on the command line, they are tried first.
//	time.google.com is used as the last resort.
//
// Options:
//
//	-w: set hwclock to system clock in UTC
package main

import (
	"flag"
	"log"

	"github.com/u-root/u-root/pkg/ntpdate"
)

var (
	config  = flag.String("config", ntpdate.DefaultNTPConfig, "NTP config file.")
	setRTC  = flag.Bool("rtc", false, "Set RTC time as well")
	verbose = flag.Bool("verbose", false, "Verbose output")
)

const (
	fallback = "time.google.com"
)

func main() {
	flag.Parse()
	if *verbose {
		ntpdate.Debug = log.Printf
	}
	server, offset, err := ntpdate.SetTime(flag.Args(), *config, fallback, *setRTC)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	plus := ""
	if offset > 0 {
		plus = "+"
	}
	log.Printf("adjust time server %s offset %s%f sec", server, plus, offset)
}
