// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/u-root/dhcp4/dhcp4client"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/dhcp6client"
	"github.com/u-root/u-root/pkg/ipxe"
	"github.com/u-root/u-root/pkg/pxe"
	"github.com/vishvananda/netlink"
)

var (
	verbose = flag.Bool("v", true, "print all kinds of things out, more than Chris wants")
	dryRun  = flag.Bool("dry-run", false, "download kernel, but don't kexec it")
	debug   = func(string, ...interface{}) {}
)

const (
	kDhcpTimeout = 30 * time.Second
	kDhcpTries   = 5
)

func leaseAndConfigureV4(iface netlink.Link, leased chan bool) (*url.URL, net.IP, error) {
	client, err := dhcp4client.New(iface,
		dhcp4client.WithTimeout(kDhcpTimeout),
		dhcp4client.WithRetry(kDhcpTries))
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Attempting to get DHCPv4 lease on %s", iface.Attrs().Name)
	p, err := client.Request()
	if err != nil {
		return nil, nil, err
	}
	packet := dhclient.NewPacket4(p)
	uri, err := packet.Boot()
	if err != nil {
		log.Printf("Got DHCPv4 lease, but no valid PXE information.")
		return nil, nil, err
	}

	_, ok := <-leased
	if !ok {
		return nil, nil, errors.New("another DHCP request offer was acepted")
	}
	close(leased)

	log.Printf("Got DHCPv4 lease on %s", iface.Attrs().Name)
	return uri, packet.Lease().IP, dhclient.Configure4(iface, packet.P)
}

func leaseAndConfigureV6(iface netlink.Link, leased chan bool) (*url.URL, net.IP, error) {
	client, err := dhcp6client.New(iface,
		dhcp6client.WithTimeout(kDhcpTimeout),
		dhcp6client.WithRetry(kDhcpTries))
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Attempting to get DHCPv6 lease on %s", iface.Attrs().Name)
	iana, p, err := client.RapidSolicit()
	if err != nil {
		return nil, nil, err
	}

	packet := dhclient.NewPacket6(p, iana)
	uri, _, err := packet.Boot()
	if err != nil {
		log.Printf("Got DHCPv6 lease, but no valid PXE information.")
		return nil, nil, err
	}

	_, ok := <-leased
	if !ok {
		return nil, nil, errors.New("another DHCP request offer was acepted")
	}
	close(leased)

	log.Printf("Got DHCPv6 lease on %s", iface.Attrs().Name)

	return uri, packet.Lease().IP, dhclient.Configure6(iface, p, iana)
}

func leaseAndConfigure(iface netlink.Link) (*url.URL, net.IP, error) {
	// leased has a message if no lease had been accepted. It is closed otherwise.
	leased := make(chan bool, 1)
	leased <- true

	type dhcpResponse struct {
		u   *url.URL
		ip  net.IP
		err error
	}

	response := make(chan dhcpResponse)

	leaseFunctions := []func(iface netlink.Link, leased chan bool) (*url.URL, net.IP, error){
		leaseAndConfigureV4,
		leaseAndConfigureV6,
	}

	for _, f := range leaseFunctions {
		f := f
		go func() {
			u, ip, err := f(iface, leased)
			response <- dhcpResponse{
				u:   u,
				ip:  ip,
				err: err,
			}
		}()
	}

	errCount := 0

	for r := range response {
		// If we got a lease, we are all set.
		// Note that one go routine will leak, but eventually timeout.
		if r.err != nil {
			return r.u, r.ip, nil
		}
		errCount++
		// All attempts failed, report to caller.
		if errCount == len(leaseFunctions) {
			return nil, nil, errors.New("unable to get DHCP lease on IPv4 nor IPv6")
		}
	}
	return nil, nil, errors.New("BUG: unreachable code")
}

// getBootImage attempts to parse the file at uri as an ipxe config and returns
// the ipxe boot image. Otherwise falls back to pxe and uses the uri directory,
// ip, and mac address to search for pxe configs.
func getBootImage(uri *url.URL, mac net.HardwareAddr, ip net.IP) (*boot.LinuxImage, error) {
	// Attempt to read the given boot path as an ipxe config file.
	if ipc, err := ipxe.NewConfig(uri); err == nil {
		log.Printf("Got configuration: %s", ipc.BootImage)
		return ipc.BootImage, nil
	}

	// Fallback to pxe boot.
	wd := &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   path.Dir(uri.Path),
	}

	pc := pxe.NewConfig(wd)
	if err := pc.FindConfigFile(mac, ip); err != nil {
		return nil, fmt.Errorf("failed to parse pxelinux config: %v", err)
	}

	label := pc.Entries[pc.DefaultEntry]
	log.Printf("Got configuration: %s", label)
	return label, nil
}

func Netboot() error {
	ifs, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, iface := range ifs {
		// TODO: Do 'em all in parallel.
		if iface.Attrs().Name != "eth0" {
			continue
		}

		uri, ip, err := leaseAndConfigure(iface)
		if err != nil {
			log.Printf("error while attempting DHCP on interface %v: %v", iface.Attrs().Name, err)
			continue
		}

		log.Printf("Boot URI: %s", uri)

		img, err := getBootImage(uri, iface.Attrs().HardwareAddr, ip)
		if err != nil {
			return err
		}

		if *dryRun {
			img.ExecutionInfo(log.New(os.Stderr, "", log.LstdFlags))
		} else if err := img.Execute(); err != nil {
			log.Printf("Kexec error: %v", err)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	if err := Netboot(); err != nil {
		log.Fatal(err)
	}
}
