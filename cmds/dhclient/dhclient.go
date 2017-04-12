// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dhclient setups up DHCP.
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

func dhclient(ifname string, timeout time.Duration, iList []string, done chan error) {
	var err error

	// if timeout is < 10 seconds, it's too short.
	if timeout < slop {
		timeout = 2 * slop
	}

	n, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/address", ifname))
	if err != nil {
		done <- fmt.Errorf("cannot get mac for %v: %v", ifname, err)
		return
	}
	// This is truly amazing but /sys appends newlines to all this data.
	n = bytes.TrimSpace(n)
	m, err := net.ParseMAC(string(n))
	if err != nil {
		done <- fmt.Errorf("mac error: %v", err)
		return
	}

	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		done <- fmt.Errorf("%s: netlink.LinkByName failed: %v", ifname, err)
		return
	}

	c, err := dhcp4client.NewInetSock(dhcp4client.SetLocalAddr(net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 68}), dhcp4client.SetRemoteAddr(net.UDPAddr{IP: net.IPv4bcast, Port: 67}))
	if err != nil {
		done <- fmt.Errorf("client Conection Generation: %v", err)
		return
	}

	client, err := dhcp4client.New(dhcp4client.HardwareAddr(m), dhcp4client.Connection(c), dhcp4client.Timeout(timeout))
	if err != nil {
		done <- fmt.Errorf("error: %v", err)
		return
	}

	// we require at least one successful request.
	success, p, err := client.Request()

	if err != nil {
		networkError, ok := err.(*net.OpError)
		if ok && networkError.Timeout() {
			done <- fmt.Errorf("%s: did not find a DHCP Server", n)
			return
		}
		done <- fmt.Errorf("%s: Error:%v", n, err)
		return
	}

	debug("Success on %s:%v\n", n, success)
	debug("Packet:%v\n", p)
	debug("Lease is %v seconds\n", p.Secs())

	if !success {
		done <- fmt.Errorf("we did not sucessfully get a DHCP Lease?")
		return
	}
	log.Printf("IP Received:%v\n", p.YIAddr().String())

	for i := 0; i < *renewals; i++ {
		addr := p.YIAddr()
		// We got here because we got a good packet.
		o := p.ParseOptions()
		netmask, ok := o[dhcp4.OptionSubnetMask]
		if ok {
			fmt.Printf("OptionSubnetMask is %v\n", netmask)
		} else {
			// what do to?
			netmask = addr
		}
		dst := &netlink.Addr{IPNet: &net.IPNet{IP: p.YIAddr(), Mask: netmask}, Label: ""}
		// Add the address to the iface.
		if err := netlink.AddrAdd(iface, dst); err != nil {
			if os.IsExist(err) {
				done <- fmt.Errorf("add %v to %v: %v", dst, n, err)
				return
			}
		}

		if gwData, ok := o[dhcp4.OptionRouter]; ok {
			fmt.Printf("router %v\n", gwData)
			routerName := net.IP(gwData).String()
			debug("routerName %v", routerName)
			r := &netlink.Route{
				Dst:       &net.IPNet{IP: p.GIAddr(), Mask: netmask},
				LinkIndex: iface.Attrs().Index,
				Gw:        p.GIAddr(),
			}

			if err := netlink.RouteReplace(r); err != nil {
				done <- fmt.Errorf("%s: add %s: %v", ifname, r.String(), routerName)
				return
			}
		}

		// We can not assume the server will give us any grace time.
		// So sleep for just a tiny bit less than the minimum.
		time.Sleep(timeout - slop)
		debug("Start Renewing Lease")
		success, p, err = client.Renew(p)
		if err != nil {
			networkError, ok := err.(*net.OpError)
			if ok && networkError.Timeout() {
				done <- fmt.Errorf("Renewal Failed! Because it didn't find the DHCP server very Strange")
				return
			}
			done <- fmt.Errorf("error: %v", err)
		}

		if !success {
			done <- fmt.Errorf("We didn't sucessfully Renew a DHCP Lease?")
			return
		}
		debug("IP Received:%v\n", p.YIAddr().String())
	}
	done <- nil
	return
}

func main() {
	done := make(chan error)
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	iList := []string{defaultIface}
	if len(flag.Args()) > 0 {
		iList = flag.Args()
	}
	for _, s := range iList {
		go dhclient(s, time.Duration(*leasetimeout)*time.Second, iList, done)
	}

	// TODO; the goroutines should pretty much run forever,
	// and only send a message on done when they are finished.
	// We can keep our own counter (don't need a sync.WaitGroup)
	// given that they send us an exit message.
	for range iList {
		if err := <-done; err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
