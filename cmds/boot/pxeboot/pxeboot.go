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
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/netboot"
	"github.com/u-root/u-root/pkg/urlfetch"
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

			if err := result.Lease.Configure(); err != nil {
				log.Printf("Failed to configure lease %s: %v", result.Lease, err)
				continue
			}
			img, err := netboot.BootImage(urlfetch.DefaultSchemes, result.Lease)
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

func main() {
	flag.Parse()

	if err := Netboot("eth0"); err != nil {
		log.Fatal(err)
	}
}
