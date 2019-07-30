// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command pxeboot implements PXE-based booting.
//
// pxeboot combines a DHCP client with a TFTP/HTTP client to download files as
// well as pxelinux and iPXE configuration file parsing.
//
// PXE-based booting requests a DHCP lease, and looks at the BootFileName and
// ServerName options (which may be embedded in the original BOOTP message, or
// as option codes) to find something to boot.
//
// This BootFileName may point to
//
// - an iPXE script beginning with #!ipxe
//
// - a pxelinux.0, in which case we will ignore the pxelinux and try to parse
//   pxelinux.cfg/<files>
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"path"
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ipxe"
	"github.com/u-root/u-root/pkg/pxe"
)

var (
	noLoad  = flag.Bool("no-load", false, "get DHCP response, but don't load the kernel")
	dryRun  = flag.Bool("dry-run", false, "download kernel, but don't kexec it")
	verbose = flag.Bool("v", false, "Verbose output")
)

const (
	dhcpTimeout = 5 * time.Second
	dhcpTries   = 3
)

// Netboot boots all interfaces matched by the regex in ifaceNames.
func Netboot(ifaceNames string) error {
	filteredIfs, err := dhclient.Interfaces(ifaceNames)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), (1<<dhcpTries)*dhcpTimeout)
	defer cancel()

	c := dhclient.Config{
		Timeout: dhcpTimeout,
		Retries: dhcpTries,
	}
	if *verbose {
		c.LogLevel = dhclient.LogSummary
	}
	r := dhclient.SendRequests(ctx, filteredIfs, true, true, c)

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

			if *noLoad {
				return nil
			}
			if err := img.Load(*dryRun); err != nil {
				return fmt.Errorf("kexec load of %v failed: %v", img, err)
			}
			if *dryRun {
				return nil
			}
			if err := boot.Execute(); err != nil {
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
	ipc, err := ipxe.ParseConfig(uri)
	if err == nil {
		return ipc, nil
	}
	log.Printf("Falling back to pxe boot: %v", err)

	// Fallback to pxe boot.
	wd := &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   path.Dir(uri.Path),
	}
	pc, err := pxe.ParseConfig(wd, mac, ip)
	if err != nil {
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
