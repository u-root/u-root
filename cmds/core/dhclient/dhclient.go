// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !freebsd && !windows

// dhclient sets up network config using DHCP.
//
// Synopsis:
//
//	dhclient [OPTIONS...]
//
// Options:
//
//	-timeout:  lease timeout in seconds
//	-renewals: number of DHCP renewals before exiting
//	-verbose:  verbose output
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

var (
	ifName   = "^e.*"
	timeout  = flag.Int("timeout", 15, "Lease timeout in seconds")
	retry    = flag.Int("retry", 5, "Max number of attempts for DHCP clients to send requests. -1 means infinity")
	dryRun   = flag.Bool("dry-run", false, "Just make the DHCP requests, but don't configure interfaces")
	verbose  = flag.Bool("v", false, "Verbose output (print message summary for each DHCP message sent/received)")
	vverbose = flag.Bool("vv", false, "Really verbose output (print all message options for each DHCP message sent/received)")
	ipv4     = flag.Bool("ipv4", true, "use IPV4")
	ipv6     = flag.Bool("ipv6", true, "use IPV6")

	v6Port   = flag.Int("v6-port", dhcpv6.DefaultServerPort, "DHCPv6 server port to send to")
	v6Server = flag.String("v6-server", "ff02::1:2", "DHCPv6 server address to send to (multicast or unicast)")

	v4Port = flag.Int("v4-port", dhcpv4.ServerPort, "DHCPv4 server port to send to")
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

	c := dhclient.Config{
		Timeout: packetTimeout,
		Retries: *retry,
		V4ServerAddr: &net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: *v4Port,
		},
		V6ServerAddr: &net.UDPAddr{
			IP:   net.ParseIP(*v6Server),
			Port: *v6Port,
		},
	}
	if *verbose {
		c.LogLevel = dhclient.LogSummary
	}
	if *vverbose {
		c.LogLevel = dhclient.LogDebug
	}
	r := dhclient.SendRequests(context.Background(), ifs, *ipv4, *ipv6, c, 30*time.Second)

	for result := range r {
		if result.Err != nil {
			log.Printf("Could not configure %s for %s: %v", result.Interface.Attrs().Name, result.Protocol, result.Err)
		} else if *dryRun {
			log.Printf("Dry run: would have configured %s with %s", result.Interface.Attrs().Name, result.Lease)
		} else if err := result.Lease.Configure(); err != nil {
			log.Printf("Could not configure %s for %s: %v", result.Interface.Attrs().Name, result.Protocol, err)
		} else {
			log.Printf("Configured %s with %s", result.Interface.Attrs().Name, result.Lease)
		}
	}
	log.Printf("Finished trying to configure all interfaces.")
}
