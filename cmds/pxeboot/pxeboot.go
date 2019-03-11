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
	"sync"
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
	dhcpTimeout = 15 * time.Second
	dhcpTries   = 3
)

type Lease interface {
	Configure() error
	Boot() (*url.URL, error)
	Link() netlink.Link
}

func lease4(iface netlink.Link) (Lease, error) {
	client, err := dhcp4client.New(iface,
		dhcp4client.WithTimeout(dhcpTimeout),
		dhcp4client.WithRetry(dhcpTries))
	if err != nil {
		return nil, err
	}

	log.Printf("Attempting to get DHCPv4 lease on %s", iface.Attrs().Name)
	p, err := client.Request()
	if err != nil {
		return nil, err
	}

	packet := dhclient.NewPacket4(iface, p)
	if _, err := packet.Boot(); err != nil {
		return nil, fmt.Errorf("valid DHCPv4 lease without PXE info: %v", err)
	}

	log.Printf("Got DHCPv4 lease on %s", iface.Attrs().Name)
	return packet, nil
}

func lease6(iface netlink.Link) (Lease, error) {
	client, err := dhcp6client.New(iface,
		dhcp6client.WithTimeout(dhcpTimeout),
		dhcp6client.WithRetry(dhcpTries))
	if err != nil {
		return nil, err
	}

	log.Printf("Attempting to get DHCPv6 lease on %s", iface.Attrs().Name)
	iana, p, err := client.RapidSolicit()
	if err != nil {
		return nil, err
	}

	packet := dhclient.NewPacket6(iface, p, iana)
	if _, err := packet.Boot(); err != nil {
		return nil, fmt.Errorf("valid DHCPv6 lease without PXE info: %v", err)
	}

	log.Printf("Got DHCPv6 lease on %s", iface.Attrs().Name)
	return packet, nil
}

// Netboot boots all interfaces matched by the regex in ifaceNames.
func Netboot(ctx context.Context, ifaceNames string) error {
	ifs, err := netlink.LinkList()
	if err != nil {
		return err
	}

	ifregex := regexp.MustCompilePOSIX(ifaceNames)

	// Yeah, this is a hack, until we can cancel all leases in progress.
	leases := make(chan Lease, 3*len(ifs))
	defer close(leases)

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, iface := range ifs {
		if !ifregex.MatchString(iface.Attrs().Name) {
			continue
		}

		wg.Add(1)
		go func(iface netlink.Link) {
			defer wg.Done()

			debug("Bringing up interface %s...", iface.Attrs().Name)
			if _, err := dhclient.IfUp(iface.Attrs().Name); err != nil {
				log.Printf("Could not bring up interface %s: %v", iface.Attrs().Name, err)
				return
			}

			wg.Add(1)
			go func(iface netlink.Link) {
				defer wg.Done()
				lease, err := lease4(iface)
				if err == nil {
					leases <- lease
				}
			}(iface)

			wg.Add(1)
			go func(iface netlink.Link) {
				defer wg.Done()
				lease, err := lease6(iface)
				if err == nil {
					leases <- lease
				}
			}(iface)
		}(iface)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case lease := <-leases:
			if err := Boot(lease); err != nil {
				log.Printf("Failed to boot lease %v: %v", lease, err)
				continue
			} else {
				return nil
			}
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

func Boot(lease Lease) error {
	if err := lease.Configure(); err != nil {
		return err
	}

	uri, err := lease.Boot()
	if err != nil {
		return err
	}
	log.Printf("Boot URI: %s", uri)

	// IP only makes sense for v4 anyway.
	var ip net.IP
	if p4, ok := lease.(*dhclient.Packet4); ok {
		ip = p4.Lease().IP
	}
	img, err := getBootImage(uri, iface.Attrs().HardwareAddr, ip)
	if err != nil {
		return err
	}
	log.Printf("Got configuration: %s", img)

	if *dryRun {
		label.ExecutionInfo(log.New(os.Stderr, "", log.LstdFlags))
		return nil
	} else if err := img.Execute(); err != nil {
		return fmt.Errorf("kexec of %v failed: %v", img, err)
	}

	// Kexec should either return an error or not return.
	panic("unreachable")
}

func main() {
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	ctx, cancel := context.WithTimeout(context.Background(), dhcpTries*dhcpTimeout)
	defer cancel()
	if err := Netboot(ctx, "eth0"); err != nil {
		log.Fatal(err)
	}
}
