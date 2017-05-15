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
	"os"
	"time"

	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
	"github.com/vishvananda/netlink"
)

const (
	defaultIfname = "eth0"
	// slop is the slop in our lease time.
	slop = 10 * time.Second
)

var (
	leasetimeout = flag.Int("timeout", 600, "Lease timeout in seconds")
	renewals     = flag.Int("renewals", -1, "Number of DHCP renewals before exiting")
	verbose      = flag.Bool("verbose", false, "Verbose output")
	ipv6         = flag.Bool("ipv6", false, "use IPV6")
	ipv4         = flag.Bool("ipv4", true, "use IPV4")
	debug        = func(string, ...interface{}) {}
)

func ifup(ifname string) (netlink.Link, error) {
	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, fmt.Errorf("%s: netlink.LinkByName failed: %v", ifname, err)
	}

	if err := netlink.LinkSetUp(iface); err != nil {
		return nil, fmt.Errorf("%v: %v can't make it up: %v", ifname, iface, err)
	}

	return iface, nil

}

func dhclient4(iface netlink.Link, numRenewals int, timeout time.Duration) error {
	mac := iface.Attrs().HardwareAddr
	conn, err := dhcp4client.NewPacketSock(iface.Attrs().Index)
	if err != nil {
		return fmt.Errorf("client conection generation: %v", err)
	}

	client, err := dhcp4client.New(dhcp4client.HardwareAddr(mac), dhcp4client.Connection(conn), dhcp4client.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	var packet dhcp4.Packet
	for i := 0; numRenewals < 0 || i < numRenewals+1; i++ {
		debug("Start getting or renewing lease")

		var success bool
		if packet == nil {
			success, packet, err = client.Request()
		} else {
			success, packet, err = client.Renew(packet)
		}
		if err != nil {
			networkError, ok := err.(*net.OpError)
			if ok && networkError.Timeout() {
				return fmt.Errorf("%s: could not find DHCP server", mac)
			}
			return fmt.Errorf("%s: error: %v", mac, err)
		}

		debug("Success on %s: %v\n", mac, success)
		debug("Packet: %v\n", packet)
		debug("Lease is %v seconds\n", packet.Secs())

		if !success {
			return fmt.Errorf("%s: we didn't sucessfully get a DHCP lease.", mac)
		}
		debug("IP Received: %v\n", packet.YIAddr().String())

		// We got here because we got a good packet.
		o := packet.ParseOptions()

		netmask, ok := o[dhcp4.OptionSubnetMask]
		if ok {
			log.Printf("OptionSubnetMask is %v\n", netmask)
		} else {
			// If they did not offer a subnet mask, we
			// choose the most restrictive option, namely,
			// our IP address.  This could happen on,
			// e.g., a point to point link.
			netmask = packet.YIAddr()
			log.Printf("No OptionSubnetMask; default to %v\n", netmask)
		}

		dst := &netlink.Addr{IPNet: &net.IPNet{IP: packet.YIAddr(), Mask: netmask}, Label: ""}
		// Replace (or add, if it presents) the address to the iface.
		if err := netlink.AddrReplace(iface, dst); err != nil {
			if os.IsExist(err) {
				return fmt.Errorf("add/replace %v to %v: %v", dst, iface, err)
			}
		}

		if gwData, ok := o[dhcp4.OptionRouter]; ok {
			log.Printf("router %v", gwData)
			routerName := net.IP(gwData).String()
			debug("routerName %v", routerName)
			r := &netlink.Route{
				LinkIndex: iface.Attrs().Index,
				Gw:        net.IP(gwData),
			}

			if err := netlink.RouteReplace(r); err != nil {
				return fmt.Errorf("%s: add %s: %v", iface.Attrs().Name, r.String(), routerName)
			}
		}

		// We can not assume the server will give us any grace time. So
		// sleep for just a tiny bit less than the minimum.
		time.Sleep(timeout - slop)
	}
	return nil
}

func dhclient6(iface netlink.Link, numRenewals int, timeout time.Duration) error {
	return fmt.Errorf("v6 is not done")
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

	iList := []string{defaultIfname}
	if len(flag.Args()) > 0 {
		iList = flag.Args()
	}

	timeout := time.Duration(*leasetimeout) * time.Second
	// if timeout is < slop, it's too short.
	if timeout < slop {
		timeout = 2 * slop
		log.Printf("increased lease timeout to %s", timeout)
	}

	done := make(chan error)
	for _, i := range iList {
		go func(ifname string) {
			iface, err := ifup(ifname)
			if err != nil {
				done <- err
				return
			}
			if *ipv4 {
				done <- dhclient4(iface, *renewals, timeout)
			}
			if *ipv6 {
				done <- dhclient6(iface, *renewals, timeout)
			}
		}(i)
	}

	// Wait for all goroutines to finish.
	for range iList {
		if err := <-done; err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
