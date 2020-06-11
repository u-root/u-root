// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package netboot provides a one-stop shop for netboot parsing needs.
//
// netboot can take a URL from a DHCP lease and try to detect iPXE scripts and
// PXE scripts.
//
// TODO: detect multiboot and Linux kernels without configuration (URL points
// to a single kernel file).
//
// TODO: detect iSCSI root paths.
package netboot

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/netboot/ipxe"
	"github.com/u-root/u-root/pkg/boot/netboot/pxe"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/vishvananda/netlink"
)

// BootImages figure out a ranked order of images to boot from the given DHCP lease.
//
// Tries, in order:
//
// - to detect an iPXE script beginning with #!ipxe,
//
// - to detect a pxelinux.0, in which case we will ignore the pxelinux.0 and
//   try to parse pxelinux.cfg/<files>.
//
// TODO: detect straight up multiboot and bzImage Linux kernel files rather
// than just configuration scripts.
func BootImages(ctx context.Context, l ulog.Logger, s curl.Schemes, lease dhclient.Lease) ([]boot.OSImage, error) {
	uri, err := lease.Boot()
	if err != nil {
		return nil, err
	}
	l.Printf("Boot URI: %s", uri)

	// IP only makes sense for v4 anyway, because the PXE probing of files
	// uses a MAC address and an IPv4 address to look at files.
	var ip net.IP
	if p4, ok := lease.(*dhclient.Packet4); ok {
		ip = p4.Lease().IP
	}
	return getBootImages(ctx, l, s, uri, lease.Link().Attrs().HardwareAddr, ip), nil
}

// getBootImages attempts to parse the file at uri as an ipxe config and returns
// the ipxe boot image. Otherwise falls back to pxe and uses the uri directory,
// ip, and mac address to search for pxe configs.
func getBootImages(ctx context.Context, l ulog.Logger, schemes curl.Schemes, uri *url.URL, mac net.HardwareAddr, ip net.IP) []boot.OSImage {
	var images []boot.OSImage

	// Attempt to read the given boot path as an ipxe config file.
	ipc, err := ipxe.ParseConfig(ctx, l, uri, schemes)
	if err != nil {
		l.Printf("Parsing boot files as iPXE failed, trying other formats...: %v", err)
	}
	if ipc != nil {
		images = append(images, ipc)
	}

	// Fallback to pxe boot.
	wd := &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   path.Dir(uri.Path),
	}
	pxeImages, err := pxe.ParseConfig(ctx, wd, mac, ip, schemes)
	if err != nil {
		l.Printf("Failed to try parsing pxelinux config: %v", err)
	}
	return append(images, pxeImages...)
}

type IPXEParser struct {
	Log     ulog.Logger
	Schemes curl.Schemes
}

func (IPXEParser) Name() string { return "iPXE" }

func (i *IPXEParser) Parse(ctx context.Context, lease dhclient.Lease) ([]boot.OSImage, error) {
	uri, err := lease.Boot()
	if err != nil {
		return nil, err
	}
	i.Log.Printf("Boot URI: %s", uri)

	// Attempt to read the given boot path as an ipxe config file.
	ipc, err := ipxe.ParseConfig(ctx, i.Log, uri, i.Schemes)
	if err != nil {
		return nil, err
	}
	if ipc != nil {
		return []boot.OSImage{ipc}, nil
	}
	return nil, nil
}

type PXEParser struct {
	Log     ulog.Logger
	Schemes curl.Schemes
}

func (PXEParser) Name() string { return "pxelinux" }

func (p *PXEParser) Parse(ctx context.Context, lease dhclient.Lease) ([]boot.OSImage, error) {
	uri, err := lease.Boot()
	if err != nil {
		return nil, err
	}
	p.Log.Printf("Boot URI: %s", uri)

	// IP only makes sense for v4 anyway, because the PXE probing of files
	// uses a MAC address and an IPv4 address to look at files.
	var ip net.IP
	if p4, ok := lease.(*dhclient.Packet4); ok {
		ip = p4.Lease().IP
	}

	mac := lease.Link().Attrs().HardwareAddr
	wd := &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   path.Dir(uri.Path),
	}
	pxeImages, err := pxe.ParseConfig(ctx, wd, mac, ip, p.Schemes)
	if err != nil {
		return nil, err
	}
	return pxeImages, nil
}

type BootImageParser interface {
	Name() string
	Parse(context.Context, dhclient.Lease) ([]boot.OSImage, error)
}

// DHCPAndParse requests DHCP (v4 + v6) on every ifs given, and parses netboot
// images from the DHCP leases. Returns bootable OSes.
func DHCPAndParse(ctx context.Context, l ulog.Logger, ifs []netlink.Link, c dhclient.Config, parsers []BootImageParser, noNetConfig bool) ([]boot.OSImage, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	l.Printf("Hit Ctrl-C to interrupt DHCP...")

	r := dhclient.SendRequests(ctx, ifs, true, true, c)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-interrupt:
			return nil, fmt.Errorf("netboot interrupted by interrupt signal (Ctrl-C)")

		case result, ok := <-r:
			if !ok {
				return nil, fmt.Errorf("nothing bootable found, all interfaces are configured or timed out")
			}

			iname := result.Interface.Attrs().Name
			// Result is either a Lease, or an Error on a specific interface.
			if result.Err != nil {
				l.Printf("Could not configure %s for %s: %v", iname, result.Protocol, result.Err)
				continue
			}

			if noNetConfig {
				l.Printf("Skipping configuring %s with lease %s", iname, result.Lease)
			} else if err := result.Lease.Configure(); err != nil {
				l.Printf("Failed to configure lease %s: %v", result.Lease, err)

				// Boot further regardless of lease configuration result.
				//
				// If lease failed, fall back to use locally configured
				// ip/ipv6 address.
			}

			var images []boot.OSImage
			for _, parser := range parsers {
				imgs, err := parser.Parse(ctx, result.Lease)
				if err != nil {
					l.Printf("Could not interpret lease %s as %s: %v", result.Lease, parser.Name(), err)
				}
				images = append(images, imgs...)
			}
			if len(images) > 0 {
				return images, nil
			}

			l.Printf("Failed to boot lease %v: no boot images found")
		}
	}
}
