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
	"github.com/vishvananda/netlink"
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

	filteredIfs, err := dhclient.Interfaces(ifName)
	if err != nil {
		log.Fatal(err)
	}

	configureAll(filteredIfs)
}

func configureAll(ifs []netlink.Link) {
	packetTimeout := time.Duration(*timeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), packetTimeout*time.Duration(1<<uint(*retry)))
	defer cancel()

	c := dhclient.Config{
		Timeout: packetTimeout,
		Retries: *retry,
	}
	if *verbose {
		c.LogLevel = dhclient.LogSummary
	}
	r := dhclient.SendRequests(ctx, ifs, *ipv4, *ipv6, c)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Done with dhclient: %v", ctx.Err())
			return

		case result, ok := <-r:
			if !ok {
				log.Printf("Configured all interfaces.")
				return
			}
			if result.Err != nil {
				log.Printf("Could not configure %s: %v", result.Interface.Attrs().Name, result.Err)
			} else if err := result.Lease.Configure(); err != nil {
				log.Printf("Could not configure %s: %v", result.Interface.Attrs().Name, err)
			} else {
				log.Printf("Configured %s with %s", result.Interface.Attrs().Name, result.Lease)
			}
		}
	}
}
