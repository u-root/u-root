// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ipxe"
	"github.com/u-root/u-root/pkg/pxe"
	"github.com/vishvananda/netlink"
)

var (
	dryRun = flag.Bool("dry-run", false, "download kernel, but don't kexec it")
)

const (
	dhcpTimeout = 15 * time.Second
	dhcpTries   = 3
)

// Netboot boots all interfaces matched by the regex in ifaceNames.
func Netboot(ifaceNames string) error {
	ifs, err := netlink.LinkList()
	if err != nil {
		return err
	}

	var filteredIfs []netlink.Link
	ifregex := regexp.MustCompilePOSIX(ifaceNames)
	for _, iface := range ifs {
		if ifregex.MatchString(iface.Attrs().Name) {
			filteredIfs = append(filteredIfs, iface)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), dhcpTries*dhcpTimeout)
	defer cancel()

	r := dhclient.SendRequests(ctx, filteredIfs, dhcpTimeout, dhcpTries, true, true)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case result, ok := <-r:
			if !ok {
				log.Printf("Configured all interfaces.")
				return fmt.Errorf("nothing bootable found")
			}
			if result.Err != nil {
				continue
			}

			img, err := Boot(result.Lease)
			if err != nil {
				log.Printf("Failed to boot lease %v: %v", result.Lease, err)
				continue
			}

			// Cancel other DHCP requests in flight.
			cancel()
			log.Printf("Got configuration: %s", img)

			if *dryRun {
				img.ExecutionInfo(log.New(os.Stderr, "", log.LstdFlags))
				return nil
			}
			if err := img.Execute(); err != nil {
				return fmt.Errorf("kexec of %v failed: %v", img, err)
			}

			// Kexec should either return an error or not return.
			panic("unreachable")
		}
	}
}

// getBootImage attempts to parse the file at uri as an ipxe config and returns
// the ipxe boot image. Otherwise falls back to pxe and uses the uri directory,
// ip, and mac address to search for pxe configs.
func getBootImage(uri *url.URL, mac net.HardwareAddr, ip net.IP) (*boot.LinuxImage, error) {
	// Attempt to read the given boot path as an ipxe config file.
	if ipc, err := ipxe.NewConfig(uri); err == nil {
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
	return label, nil
}

func Boot(lease dhclient.Lease) (*boot.LinuxImage, error) {
	if err := lease.Configure(); err != nil {
		return nil, err
	}

	uri, err := lease.Boot()
	if err != nil {
		return nil, err
	}
	log.Printf("Boot URI: %s", uri)

	// IP only makes sense for v4 anyway, because the PXE probing of files
	// uses a MAC address and an IPv4 address to look at files.
	var ip net.IP
	if p4, ok := lease.(*dhclient.Packet4); ok {
		ip = p4.Lease().IP
	}
	return getBootImage(uri, lease.Link().Attrs().HardwareAddr, ip)
}

func main() {
	flag.Parse()

	if err := Netboot("eth0"); err != nil {
		log.Fatal(err)
	}
}
