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
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4client"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/dhcp6client"
	"github.com/vishvananda/netlink"
)

const (
	// slop is the slop in our lease time.
	slop          = 10 * time.Second
	linkUpAttempt = 30 * time.Second
)

var (
	ifName         = "^e.*"
	leasetimeout   = flag.Int("timeout", 15, "Lease timeout in seconds")
	retry          = flag.Int("retry", 5, "Max number of attempts for DHCP clients to send requests. -1 means infinity")
	renewals       = flag.Int("renewals", 0, "Number of DHCP renewals before exiting. -1 means infinity")
	renewalTimeout = flag.Int("renewal timeout", 3600, "How long to wait before renewing in seconds")
	verbose        = flag.Bool("verbose", false, "Verbose output")
	ipv4           = flag.Bool("ipv4", true, "use IPV4")
	ipv6           = flag.Bool("ipv6", true, "use IPV6")
	test           = flag.Bool("test", false, "Test mode")
	debug          = func(string, ...interface{}) {}
)

func ifup(ifname string) (netlink.Link, error) {
	debug("Try bringing up %v", ifname)
	start := time.Now()
	for time.Since(start) < linkUpAttempt {
		// Note that it may seem odd to keep trying the
		// LinkByName operation, by consider that a hotplug
		// device such as USB ethernet can just vanish.
		iface, err := netlink.LinkByName(ifname)
		debug("LinkByName(%v) returns (%v, %v)", ifname, iface, err)
		if err != nil {
			return nil, fmt.Errorf("cannot get interface by name %v: %v", ifname, err)
		}

		if iface.Attrs().OperState == netlink.OperUp {
			debug("Link %v is up", ifname)
			return iface, nil
		}

		if err := netlink.LinkSetUp(iface); err != nil {
			return nil, fmt.Errorf("%v: %v can't make it up: %v", ifname, iface, err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("Link %v still down after %d seconds", ifname, linkUpAttempt)
}

func dhclient4(iface netlink.Link, timeout time.Duration, retry int, numRenewals int) error {
	ifa, err := net.InterfaceByIndex(iface.Attrs().Index)
	if err != nil {
		return err
	}

	client, err := dhcp4client.New(ifa /*, timeout, retry*/)
	if err != nil {
		return err
	}

	for i := 0; numRenewals < 0 || i <= numRenewals; i++ {
		var packet *dhcp4.Packet
		var err error
		if i == 0 {
			packet, err = client.Request()
		} else {
			time.Sleep(time.Duration(*renewalTimeout) * time.Second)
			packet, err = client.Renew(packet)
		}
		if err != nil {
			return err
		}

		if err := dhclient.Configure4(iface, packet); err != nil {
			return err
		}
	}
	return nil
}

func dhclient6(iface netlink.Link, timeout time.Duration, retry int, numRenewals int) error {
	client, err := dhcp6client.New(iface,
		dhcp6client.WithTimeout(timeout),
		dhcp6client.WithRetry(retry))
	if err != nil {
		return err
	}

	for i := 0; numRenewals < 0 || i <= numRenewals; i++ {
		if i != 0 {
			time.Sleep(time.Duration(*renewalTimeout) * time.Second)
		}
		iana, packet, err := client.RapidSolicit()
		if err != nil {
			return err
		}

		if err := dhclient.Configure6(iface, packet, iana); err != nil {
			return err
		}
	}
	return nil
}

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

	timeout := time.Duration(*leasetimeout) * time.Second
	// if timeout is < slop, it's too short.
	if timeout < slop {
		timeout = 2 * slop
		log.Printf("increased lease timeout to %s", timeout)
	}

	var wg sync.WaitGroup
	done := make(chan error)
	for _, i := range ifnames {
		if !ifRE.MatchString(i.Attrs().Name) {
			continue
		}
		wg.Add(1)
		go func(ifname string) {
			defer wg.Done()
			iface, err := dhclient.IfUp(ifname)
			if err != nil {
				done <- err
				return
			}
			if *ipv4 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					done <- dhclient4(iface, timeout, *retry, *renewals)
				}()
			}
			if *ipv6 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					done <- dhclient6(iface, timeout, *retry, *renewals)
				}()
			}
			debug("Done dhclient for %v", ifname)
		}(i.Attrs().Name)
	}

	go func() {
		wg.Wait()
		close(done)
	}()
	// Wait for all goroutines to finish.
	var nif int
	for err := range done {
		debug("err from done %v", err)
		if err != nil {
			log.Print(err)
		}
		nif++
	}

	if nif == 0 {
		log.Fatalf("No interfaces match %v\n", ifName)
	}
	fmt.Printf("%d dhclient attempts were sent", nif)
}
