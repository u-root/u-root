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
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
	"github.com/vishvananda/netlink"
)

const (
	defaultIface = "eth0"
	// slop is the slop in our lease time.
	slop = 10 * time.Second
)

var (
	leasetimeout = flag.Int("timeout", 600, "Lease timeout in seconds")
	renewals     = flag.Int("renewals", 2, "Number of DHCP renewals before exiting")
	verbose      = flag.Bool("verbose", false, "Verbose output")
	debug        = func(string, ...interface{}) {}
)

func dhclient(ifname string, numRenewals int, timeout time.Duration) error {
	var err error

	// if timeout is < 10 seconds, it's too short.
	if timeout < slop {
		timeout = 2 * slop
		log.Printf("increased log timeout to %s", timeout)
	}

	n, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/address", ifname))
	if err != nil {
		return fmt.Errorf("cannot get mac for %v: %v", ifname, err)
	}

	// This is truly amazing but /sys appends newlines to all this data.
	n = bytes.TrimSpace(n)
	mac, err := net.ParseMAC(string(n))
	if err != nil {
		return fmt.Errorf("mac error: %v", err)
	}

	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		return fmt.Errorf("%s: netlink.LinkByName failed: %v", ifname, err)
	}

	conn, err := dhcp4client.NewInetSock()
	if err != nil {
		return fmt.Errorf("client conection generation: %v", err)
	}

	client, err := dhcp4client.New(dhcp4client.HardwareAddr(mac), dhcp4client.Connection(conn), dhcp4client.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	var packet dhcp4.Packet
	for i := 0; i < numRenewals+1; i++ {
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
			// Default subnet mask is the received IP?
			//
			// N.B.(chrisko): This doesn't feel right?
			netmask = packet.YIAddr()
		}

		dst := &netlink.Addr{IPNet: &net.IPNet{IP: packet.YIAddr(), Mask: netmask}, Label: ""}
		// Add the address to the iface.
		if err := netlink.AddrAdd(iface, dst); err != nil {
			if os.IsExist(err) {
				return fmt.Errorf("add %v to %v: %v", dst, n, err)
			}
		}

		if gwData, ok := o[dhcp4.OptionRouter]; ok {
			log.Printf("router %v", gwData)
			routerName := net.IP(gwData).String()
			debug("routerName %v", routerName)
			r := &netlink.Route{
				Dst:       &net.IPNet{IP: packet.GIAddr(), Mask: netmask},
				LinkIndex: iface.Attrs().Index,
				Gw:        packet.GIAddr(),
			}

			if err := netlink.RouteReplace(r); err != nil {
				return fmt.Errorf("%s: add %s: %v", ifname, r.String(), routerName)
			}
		}

		// We can not assume the server will give us any grace time. So
		// sleep for just a tiny bit less than the minimum.
		time.Sleep(timeout - slop)
	}
	return nil
}

func main() {
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	iList := []string{defaultIface}
	if len(flag.Args()) > 0 {
		iList = flag.Args()
	}

	done := make(chan error)
	for _, iface := range iList {
		go func(iface string) {
			done <- dhclient(iface, *renewals, time.Duration(*leasetimeout)*time.Second)
		}(iface)
	}

	// Wait for all goroutines to finish.
	for range iList {
		if err := <-done; err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
