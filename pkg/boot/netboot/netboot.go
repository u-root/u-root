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
	"fmt"
	"net"
	"net/url"
	"path"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/netboot/ipxe"
	"github.com/u-root/u-root/pkg/boot/netboot/pxe"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ulog"
)

// BootImage figures out the image to boot from the given DHCP lease.
//
// Tries, in order:
//
// - to detect an iPXE script beginning with #!ipxe,
//
// - to detect a pxelinux.0, in which case we will ignore the pxelinux and try
//   to parse pxelinux.cfg/<files>.
//
// TODO: detect straight up multiboot and bzImage Linux kernel files rather
// than just configuration scripts.
func BootImage(l ulog.Logger, s curl.Schemes, lease dhclient.Lease) (*boot.LinuxImage, error) {
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
	return getBootImage(l, s, uri, lease.Link().Attrs().HardwareAddr, ip)
}

// getBootImage attempts to parse the file at uri as an ipxe config and returns
// the ipxe boot image. Otherwise falls back to pxe and uses the uri directory,
// ip, and mac address to search for pxe configs.
func getBootImage(l ulog.Logger, schemes curl.Schemes, uri *url.URL, mac net.HardwareAddr, ip net.IP) (*boot.LinuxImage, error) {
	// Attempt to read the given boot path as an ipxe config file.
	ipc, err := ipxe.ParseConfig(l, uri, schemes)
	if err == nil {
		return ipc, nil
	}
	l.Printf("Falling back to pxe boot: %v", err)

	// Fallback to pxe boot.
	wd := &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   path.Dir(uri.Path),
	}
	pc, err := pxe.ParseConfig(wd, mac, ip, schemes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pxelinux config: %v", err)
	}

	label, ok := pc.Entries[pc.DefaultEntry]
	if !ok {
		return nil, fmt.Errorf("Could not find %q from entries %v", pc.DefaultEntry, pc.Entries)
	}
	return label, nil
}
