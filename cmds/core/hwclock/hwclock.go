// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hwclock reads or changes the hardware clock (RTC) in UTC format.
//
// Synopsis:
//
//	hwclock [-w]
//
// Description:
//
//	It prints the current hwclock time in UTC if called without any flags.
//	It sets the hwclock to the system clock in UTC if called with -w.
//
// Options:
//
//	-w: set hwclock to system clock in UTC
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/u-root/u-root/pkg/rtc"
)

var write = flag.Bool("w", false, "Set hwclock from system clock in UTC")

func main() {
	flag.Parse()

	r, err := rtc.OpenRTC()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	if *write {
		tu := time.Now().UTC()
		if err := r.Set(tu); err != nil {
			log.Fatal(err)
		}
	}

	t, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Print local time. Match the format of util-linux' hwclock.
	fmt.Println(t.Local().Format("Mon 2 Jan 2006 15:04:05 AM MST"))
}
