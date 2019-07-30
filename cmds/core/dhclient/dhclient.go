// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dhclient sets up DHCP.
//
// Synopsis:
//     dhclient [OPTIONS...]
//
// Options:
//     -timeout:  lease timeout in seconds
//     -renewals: number of DHCP renewals before exiting
//     -verbose:  verbose output
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/u-root/u-root/pkg/dhclient"
)

var (
	ifName  = "^e.*"
	timeout = flag.Int("timeout", 15, "Lease timeout in seconds")
	retry   = flag.Int("retry", 5, "Max number of attempts for DHCP clients to send requests. -1 means infinity")
	verbose = flag.Bool("v", false, "Verbose output")
	ipv4    = flag.Bool("ipv4", true, "use IPV4")
	ipv6    = flag.Bool("ipv6", true, "use IPV6")
)

func main() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		log.Fatalf("only one re")
	}

	if len(flag.Args()) > 0 {
		ifName = flag.Args()[0]
	}

	c := dhclient.Config{
		Timeout: time.Duration(*timeout) * time.Second,
		Retries: *retry,
	}
	if *verbose {
		c.LogLevel = dhclient.LogSummary
	}
	if err := dhclient.ConfigureAll(context.Background(), ifName, c, *ipv4, *ipv6); err != nil {
		log.Fatalf("ConfigureAll(%s) failed: %v", ifName, err)
	} else {
		log.Printf("Configured all interfaces.")
	}
}
