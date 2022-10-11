// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package netboot provides a one-stop shop for netboot parsing needs.
//
// netboot can take a URL from a DHCP lease and try to detect iPXE scripts and
// PXE scripts.
//
// TODO: detect iSCSI root paths.
package netboot

import (
	"context"
	"net"
	"net/url"
	"path"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/netboot/ipxe"
	"github.com/u-root/u-root/pkg/boot/netboot/pxe"
	"github.com/u-root/u-root/pkg/boot/netboot/simple"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ulog"
)

// BootImages figure out a ranked order of images to boot from the given DHCP lease.
//
// Tries, in order:
//
// - to detect an iPXE script beginning with #!ipxe,
//
//   - to detect a pxelinux.0, in which case we will ignore the pxelinux.0 and
//     try to parse pxelinux.cfg/<files>.
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

	// 1: Attempt to download the given url as is.
	//
	// 1.1: Try ipxe config file.
	ipc, err := ipxe.ParseConfig(ctx, l, uri, schemes)
	if err != nil {
		l.Printf("Parsing boot files as iPXE failed, trying other formats...: %v", err)
	}
	if ipc != nil {
		images = append(images, ipc)
	}

	// 1.2: Check if target is a simple file instead of config script
	if ipc == nil {
		l.Printf("Trying to parse file as a non config Image...")
		sImages, err := simple.FetchAndProbe(ctx, uri, schemes)
		if err != nil {
			l.Printf("failed to parse boot file as simple file: %v", err)
		}
		if sImages != nil {
			images = append(images, sImages...)
		}
	}

	// 2: Fallback to pxe boot.
	//
	// Look for pxelinux.cfg from parent directory of given url path.
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
