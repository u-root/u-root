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
	"crypto/rand"
	"flag"
	"log"
	"regexp"
	"time"

	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

var (
	ifName  = "^e.*"
	timeout = flag.Int("timeout", 15, "Lease timeout in seconds")
	retry   = flag.Int("retry", 5, "Max number of attempts for DHCP clients to send requests. -1 means infinity")
	verbose = flag.Bool("verbose", false, "Verbose output")
	ipv4    = flag.Bool("ipv4", true, "use IPV4")
	ipv6    = flag.Bool("ipv6", true, "use IPV6")
	test    = flag.Bool("test", false, "Test mode")
	debug   = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	// if we boot quickly enough, the random number generator
	// may not be ready, and the dhcp package panics in that case.
	// Worse, /dev/urandom, which the Go package falls back to,
	// might not be there. Still worse, the Go package is "sticky"
	// in that once it decides to use /dev/urandom, it won't go back,
	// even if the system call would subsequently work.
	// You're screwed. Exit.
	// Wouldn't it be nice if we could just do the blocking system
	// call? But that comes with its own giant set of headaches.
	// Maybe we'll end up in a loop, sleeping, and just running
	// ourselves.
	if n, err := rand.Read([]byte{0}); err != nil || n != 1 {
		log.Fatalf("We're sorry, the random number generator is not up. Please file a ticket")
	}

	if len(flag.Args()) > 1 {
		log.Fatalf("only one re")
	}

	if len(flag.Args()) > 0 {
		ifName = flag.Args()[0]
	}

	ifRE := regexp.MustCompilePOSIX(ifName)

	ifnames, err := netlink.LinkList()
	if err != nil {
		log.Fatalf("Can't get list of link names: %v", err)
	}

	var filteredIfs []netlink.Link
	for _, iface := range ifnames {
		if ifRE.MatchString(iface.Attrs().Name) {
			filteredIfs = append(filteredIfs, iface)
		}
	}

	if len(filteredIfs) == 0 {
		log.Fatalf("No interfaces match %s", ifName)
	}

	configureAll(filteredIfs)
}

func configureAll(ifs []netlink.Link) {
	packetTimeout := time.Duration(*timeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), packetTimeout*time.Duration(*retry))
	defer cancel()

	r := dhclient.SendRequests(ctx, ifs, packetTimeout, *retry, *ipv4, *ipv6)

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
