// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

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
	"fmt"
	"log"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
)

var (
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

type opts struct {
	timeout  int
	retry    int
	dryRun   bool
	verbose  bool
	vverbose bool
	ipv4     bool
	ipv6     bool
	v6Port   int
	v6Server string
	v4Port   int
}

func run(opts *opts, args []string) error {
	var ifName string
	switch len(args) {
	case 0:
		ifName = "^e.*"
	case 1:
		ifName = args[0]
	default:
		return fmt.Errorf("more than one interface specified")
	}

	filteredIfs, err := dhclient.Interfaces(ifName)
	if err != nil {
		return err
	}

	if err := configureAll(filteredIfs, opts); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	opts := opts{timeout: *timeout, retry: *retry, dryRun: *dryRun, verbose: *verbose, vverbose: *vverbose, ipv4: *ipv4, ipv6: *ipv6, v6Port: *v6Port, v6Server: *v6Server, v4Port: *v4Port}
	if err := run(&opts, flag.Args()); err != nil {
		log.Fatalf("dhclient: %v", err)
	}
}

func configureAll(ifs []netlink.Link, opts *opts) error {
	packetTimeout := time.Duration(opts.timeout) * time.Second

	c := dhclient.Config{
		Timeout: packetTimeout,
		Retries: opts.retry,
		V4ServerAddr: &net.UDPAddr{
			IP:   net.IPv4bcast,
			Port: opts.v4Port,
		},
		V6ServerAddr: &net.UDPAddr{
			IP:   net.ParseIP(opts.v6Server),
			Port: opts.v6Port,
		},
	}
	if opts.verbose {
		c.LogLevel = dhclient.LogSummary
	}
	if opts.vverbose {
		c.LogLevel = dhclient.LogDebug
	}
	r := dhclient.SendRequests(context.Background(), ifs, opts.ipv4, opts.ipv6, c, 30*time.Second)

	for result := range r {
		if result.Err != nil {
			log.Printf("Could not configure %s for %s: %v", result.Interface.Attrs().Name, result.Protocol, result.Err)
		} else if opts.dryRun {
			log.Printf("Dry run: would have configured %s with %s", result.Interface.Attrs().Name, result.Lease)
		} else if err := result.Lease.Configure(); err != nil {
			log.Printf("Could not configure %s for %s: %v", result.Interface.Attrs().Name, result.Protocol, err)
		} else {
			log.Printf("Configured %s with %s", result.Interface.Attrs().Name, result.Lease)
		}
	}
	log.Printf("Finished trying to configure all interfaces.")
	return nil
}
