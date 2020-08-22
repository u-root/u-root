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
	"time"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bootcmd"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/boot/netboot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	ifName      = "^e.*"
	noLoad      = flag.Bool("no-load", false, "get DHCP response, print chosen boot configuration, but do not download + exec it")
	noExec      = flag.Bool("no-exec", false, "download boot configuration, but do not exec it")
	noNetConfig = flag.Bool("no-net-config", false, "get DHCP response, but do not apply the network config it to the kernel interface")
	verbose     = flag.Bool("v", false, "Verbose output")
)

const (
	dhcpTimeout = 5 * time.Second
	dhcpTries   = 3
)

// NetbootImages requests DHCP on every ifaceNames interface, and parses
// netboot images from the DHCP leases. Returns bootable OSes.
func NetbootImages(ifaceNames string) ([]boot.OSImage, error) {
	filteredIfs, err := dhclient.Interfaces(ifaceNames)
	if err != nil {
		return nil, err
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
	r := dhclient.SendRequests(ctx, filteredIfs, true, true, c, 30*time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case result, ok := <-r:
			if !ok {
				return nil, fmt.Errorf("nothing bootable found, all interfaces are configured or timed out")
			}
			iname := result.Interface.Attrs().Name
			if result.Err != nil {
				log.Printf("Could not configure %s for %s: %v", iname, result.Protocol, result.Err)
				continue
			}

			if *noNetConfig {
				log.Printf("Skipping configuring %s with lease %s", iname, result.Lease)
			} else if err := result.Lease.Configure(); err != nil {
				log.Printf("Failed to configure lease %s: %v", result.Lease, err)
				// Boot further regardless of lease configuration result.
				//
				// If lease failed, fall back to use locally configured
				// ip/ipv6 address.
			}

			// Don't use the other context, as it's for the DHCP timeout.
			imgs, err := netboot.BootImages(context.Background(), ulog.Log, curl.DefaultSchemes, result.Lease)
			if err != nil {
				log.Printf("Failed to boot lease %v: %v", result.Lease, err)
				continue
			}
			return imgs, nil
		}
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		log.Fatalf("Only one regexp-style argument is allowed, e.g.: " + ifName)
	}
	if len(flag.Args()) > 0 {
		ifName = flag.Args()[0]
	}

	images, err := NetbootImages(ifName)
	if err != nil {
		log.Printf("Netboot failed: %v", err)
	}

	menuEntries := menu.OSImages(*verbose, images...)
	menuEntries = append(menuEntries, menu.Reboot{})
	menuEntries = append(menuEntries, menu.StartShell{})

	// Boot does not return.
	bootcmd.ShowMenuAndBoot(menuEntries, nil, *noLoad, *noExec)
}
