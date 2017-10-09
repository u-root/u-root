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
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
	"github.com/u-root/dhcp6"
	"github.com/vishvananda/netlink"
)

const (
	// slop is the slop in our lease time.
	slop          = 10 * time.Second
	linkUpAttempt = 30 * time.Second
)

var (
	ifName       = "^e.*"
	leasetimeout = flag.Int("timeout", 15, "Lease timeout in seconds")
	retry        = flag.Int("retry", -1, "Max number of attempts for DHCP clients to send requests. -1 means infinity")
	renewals     = flag.Int("renewals", -1, "Number of DHCP renewals before exiting. -1 means infinity")
	verbose      = flag.Bool("verbose", false, "Verbose output")
	ipv4         = flag.Bool("ipv4", true, "use IPV4")
	ipv6         = flag.Bool("ipv6", true, "use IPV6")
	test         = flag.Bool("test", false, "Test mode")
	debug        = func(string, ...interface{}) {}
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

func dhclient4(iface netlink.Link, numRenewals int, timeout time.Duration, retry int) error {
	mac := iface.Attrs().HardwareAddr
	conn, err := dhcp4client.NewPacketSock(iface.Attrs().Index)
	if err != nil {
		return fmt.Errorf("client connection generation: %v", err)
	}

	client, err := dhcp4client.New(dhcp4client.HardwareAddr(mac), dhcp4client.Connection(conn), dhcp4client.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	var packet dhcp4.Packet
	needsRequest := true
	for i := 0; numRenewals < 0 || i < numRenewals+1; i++ {
		debug("Start getting or renewing DHCPv4 lease")

		var success bool
		for i := 0; i < retry || retry < 0; i++ {
			if i > 0 {
				if needsRequest {
					debug("Resending DHCPv4 request...\n")
				} else {
					debug("Resending DHCPv4 renewal")
				}
			}

			if needsRequest {
				success, packet, err = client.Request()
			} else {
				success, packet, err = client.Renew(packet)
			}
			if err != nil {
				if err0, ok := err.(net.Error); ok && err0.Timeout() {
					log.Printf("%s: timeout contacting DHCP server", mac)
				} else {
					log.Printf("%s: error: %v", mac, err)
				}
			} else {
				// Client needs renew after no matter what state it is now.
				needsRequest = false
				break
			}
		}

		debug("Success on %s: %v\n", mac, success)
		debug("Packet: %v\n", packet)
		debug("Lease is %v seconds\n", packet.Secs())

		if !success {
			return fmt.Errorf("%s: we didn't successfully get a DHCP lease", mac)
		}
		debug("IP Received: %v\n", packet.YIAddr().String())

		// We got here because we got a good packet.
		o := packet.ParseOptions()

		netmask, ok := o[dhcp4.OptionSubnetMask]
		if ok {
			debug("OptionSubnetMask is %v\n", netmask)
		} else {
			// If they did not offer a subnet mask, we
			// choose the most restrictive option, namely,
			// our IP address.  This could happen on,
			// e.g., a point to point link.
			netmask = packet.YIAddr()
			debug("No OptionSubnetMask; default to %v\n", netmask)
		}

		dst := &netlink.Addr{IPNet: &net.IPNet{IP: packet.YIAddr(), Mask: netmask}, Label: ""}
		// Add the address to the iface.
		if *test == false {
			if err := netlink.AddrReplace(iface, dst); err != nil {
				if os.IsExist(err) {
					return fmt.Errorf("add/replace %v to %v: %v", dst, iface, err)
				}
			}

			if gwData, ok := o[dhcp4.OptionRouter]; ok {
				debug("router %v", gwData)
				routerName := net.IP(gwData).String()
				debug("routerName %v", routerName)
				r := &netlink.Route{
					Dst:       &net.IPNet{IP: packet.GIAddr()},
					LinkIndex: iface.Attrs().Index,
					Gw:        net.IP(gwData),
				}

				if err := netlink.RouteReplace(r); err != nil {
					return fmt.Errorf("%s: add %s: %v", iface.Attrs().Name, r.String(), routerName)
				}
			}
		}
		if binary.BigEndian.Uint16(packet.Secs()) == 0 {
			debug("%v: server returned infinite lease.", iface.Attrs().Name)
			break
		}

		// We can not assume the server will give us any grace time. So
		// sleep for just a tiny bit less than the minimum.
		time.Sleep(timeout - slop)
	}
	return nil
}

// dhcp6 support in go is hard to find. This function represents our best current
// guess based on reading and testing.
func dhclient6(iface netlink.Link, numRenewals int, timeout time.Duration, retry int) error {
	conn, err := dhcp6.NewPacketSock(iface.Attrs().Index)
	if err != nil {
		return fmt.Errorf("client connection generation: %v", err)
	}
	client := dhcp6.New(iface.Attrs().HardwareAddr, conn, timeout, retry)

	for i := 0; numRenewals < 0 || i < numRenewals+1; i++ {
		debug("Start getting or renewing DHCPv6 lease")
		iaAddrs, packet, err := client.Solicit()
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		debug("Packet: %+v\n\n", packet)
		debug("IAAddrs: %v\n", iaAddrs)

		if *test == false {
			dst := &netlink.Addr{
				IPNet:       &net.IPNet{IP: iaAddrs[0].IP},
				PreferedLft: int(iaAddrs[0].PreferredLifetime.Seconds()),
				ValidLft:    int(iaAddrs[0].ValidLifetime.Seconds()),
				Label:       "",
			}

			if err := netlink.AddrReplace(iface, dst); err != nil {
				if os.IsExist(err) {
					return fmt.Errorf("add/replace %v to %v: %v", dst, iface, err)
				}
			}
		}

		time.Sleep(timeout - slop)
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
			iface, err := ifup(ifname)
			if err != nil {
				done <- err
				return
			}
			if *ipv4 {
				wg.Add(1)
				done <- dhclient4(iface, *renewals, timeout, *retry)
				wg.Done()
			}
			if *ipv6 {
				wg.Add(1)
				done <- dhclient6(iface, *renewals, timeout, *retry)
				wg.Done()
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		nif++
	}

	if nif == 0 {
		fmt.Printf("No interfaces match %v\n", ifName)
	}
}
