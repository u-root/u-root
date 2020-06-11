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
	"log"
	"time"

	"github.com/u-root/u-root/pkg/boot/bootcmd"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/boot/netboot"
	"github.com/u-root/u-root/pkg/boot/netboot/iscsi"
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

func main() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		log.Fatalf("Only one regexp-style argument is allowed, e.g.: " + ifName)
	}

	if len(flag.Args()) > 0 {
		ifName = flag.Args()[0]
	}

	conf := dhclient.Config{
		Timeout: dhcpTimeout,
		Retries: dhcpTries,
	}
	if *verbose {
		conf.LogLevel = dhclient.LogSummary
	}

	filteredIfs, err := dhclient.Interfaces(ifName)
	if err != nil {
		log.Fatalf("Netboot interfaces failed: %v", err)
	}

	parsers := []netboot.BootImageParser{
		&netboot.IPXEParser{Log: ulog.Log, Schemes: curl.DefaultSchemes},
		&netboot.PXEParser{Log: ulog.Log, Schemes: curl.DefaultSchemes},
		&iscsi.ISCSIBoot{
			Log:        ulog.Log,
			CreateIBFT: true,
			DiskParsers: []iscsi.DiskParser{
				iscsi.ESXIBoot{},
			},
		},
	}
	entries, err := netboot.DHCPAndParse(context.Background(), ulog.Log, filteredIfs, conf, parsers, *noNetConfig)
	if err != nil {
		log.Printf("Netboot failed: %v", err)
	}

	menuEntries := menu.OSImages(*verbose, images...)
	menuEntries = append(menuEntries, menu.Reboot{})
	menuEntries = append(menuEntries, menu.StartShell{})

	// Boot does not return.
	bootcmd.ShowMenuAndBoot(menuEntries, nil, *noLoad, *noExec)
}
