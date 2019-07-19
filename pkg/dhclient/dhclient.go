// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dhclient provides a unified interface for interfacing with both
// DHCPv4 and DHCPv6 clients.
package dhclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/nclient4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/nclient6"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const linkUpAttempt = 30 * time.Second

// isIpv6LinkReady returns true iff the interface has a link-local address
// which is not tentative.
func isIpv6LinkReady(l netlink.Link) (bool, error) {
	addrs, err := netlink.AddrList(l, netlink.FAMILY_V6)
	if err != nil {
		return false, err
	}
	for _, addr := range addrs {
		if addr.IP.IsLinkLocalUnicast() && (addr.Flags&unix.IFA_F_TENTATIVE == 0) {
			if addr.Flags&unix.IFA_F_DADFAILED != 0 {
				log.Printf("DADFAILED for %v, continuing anyhow", addr.IP)
			}
			return true, nil
		}
	}
	return false, nil
}

// IfUp ensures the given network interface is up and returns the link object.
func IfUp(ifname string) (netlink.Link, error) {
	start := time.Now()
	for time.Since(start) < linkUpAttempt {
		// Note that it may seem odd to keep trying the LinkByName
		// operation, but consider that a hotplug device such as USB
		// ethernet can just vanish.
		iface, err := netlink.LinkByName(ifname)
		if err != nil {
			return nil, fmt.Errorf("cannot get interface %q by name: %v", ifname, err)
		}

		if iface.Attrs().Flags&net.FlagUp == net.FlagUp {
			return iface, nil
		}

		if err := netlink.LinkSetUp(iface); err != nil {
			return nil, fmt.Errorf("interface %q: %v can't make it up: %v", ifname, iface, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("link %q still down after %d seconds", ifname, linkUpAttempt)
}

// Configure4 adds IP addresses, routes, and DNS servers to the system.
func Configure4(iface netlink.Link, packet *dhcpv4.DHCPv4) error {
	p := NewPacket4(iface, packet)
	return p.Configure()
}

// Configure6 adds IPv6 addresses, routes, and DNS servers to the system.
func Configure6(iface netlink.Link, packet *dhcpv6.Message) error {
	p := NewPacket6(iface, packet)

	l := p.Lease()
	if l == nil {
		return fmt.Errorf("no lease returned")
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   l.IPv6Addr,
			Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::")),
		},
		PreferedLft: int(l.PreferredLifetime),
		ValidLft:    int(l.ValidLifetime),
		// Optimistic DAD (Duplicate Address Detection) means we can
		// use the address before DAD is complete. The DHCP server's
		// job was to give us a unique IP so there is little risk of a
		// collision.
		Flags: unix.IFA_F_OPTIMISTIC,
	}
	if err := netlink.AddrReplace(iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %v", dst, iface, err)
		}
	}

	if ips := p.DNS(); ips != nil {
		if err := WriteDNSSettings(ips, nil, ""); err != nil {
			return err
		}
	}
	return nil
}

// WriteDNSSettings writes the given nameservers, search list, and domain to resolv.conf.
func WriteDNSSettings(ns []net.IP, sl []string, domain string) error {
	rc := &bytes.Buffer{}
	if domain != "" {
		rc.WriteString(fmt.Sprintf("domain %s\n", domain))
	}
	if ns != nil {
		for _, ip := range ns {
			rc.WriteString(fmt.Sprintf("nameserver %s\n", ip))
		}
	}
	if sl != nil {
		rc.WriteString("search ")
		rc.WriteString(strings.Join(sl, " "))
		rc.WriteString("\n")
	}
	return ioutil.WriteFile("/etc/resolv.conf", rc.Bytes(), 0644)
}

// Lease is a network configuration obtained by DHCP.
type Lease interface {
	fmt.Stringer

	// Configure configures the associated interface with the network
	// configuration.
	Configure() error

	// Boot is a URL to obtain booting information from that was part of
	// the network config.
	Boot() (*url.URL, error)

	// Link is the interface the configuration is for.
	Link() netlink.Link
}

// LogLevel is the amount of information to log.
type LogLevel uint8

// LogLevel are the levels.
const (
	LogInfo    LogLevel = 0
	LogSummary LogLevel = 1
	LogDebug   LogLevel = 2
)

// Config is a DHCP client configuration.
type Config struct {
	// Timeout is the timeout for one DHCP request attempt.
	Timeout time.Duration

	// Retries is how many times to retry DHCP attempts.
	Retries int

	// LogLevel determines the amount of information printed for each
	// attempt. The highest log level should print each entire packet sent
	// and received.
	LogLevel LogLevel
}

func lease4(ctx context.Context, iface netlink.Link, c Config) (Lease, error) {
	mods := []nclient4.ClientOpt{
		nclient4.WithTimeout(c.Timeout),
		nclient4.WithRetry(c.Retries),
	}
	switch c.LogLevel {
	case LogSummary:
		mods = append(mods, nclient4.WithSummaryLogger())
	case LogDebug:
		mods = append(mods, nclient4.WithDebugLogger())
	}
	client, err := nclient4.New(iface.Attrs().Name, mods...)
	if err != nil {
		return nil, err
	}

	log.Printf("Attempting to get DHCPv4 lease on %s", iface.Attrs().Name)
	_, p, err := client.Request(ctx, dhcpv4.WithNetboot,
		dhcpv4.WithOption(dhcpv4.OptClassIdentifier("PXE UROOT")),
		dhcpv4.WithRequestedOptions(dhcpv4.OptionSubnetMask))
	if err != nil {
		return nil, err
	}

	packet := NewPacket4(iface, p)
	log.Printf("Got DHCPv4 lease on %s: %v", iface.Attrs().Name, p.Summary())
	return packet, nil
}

func lease6(ctx context.Context, iface netlink.Link, c Config) (Lease, error) {
	// For ipv6, we cannot bind to the port until Duplicate Address
	// Detection (DAD) is complete which is indicated by the link being no
	// longer marked as "tentative". This usually takes about a second.
	for {
		if ready, err := isIpv6LinkReady(iface); err != nil {
			return nil, err
		} else if ready {
			break
		}
		select {
		case <-time.After(100 * time.Millisecond):
			continue
		case <-ctx.Done():
			return nil, errors.New("timeout after waiting for a non-tentative IPv6 address")
		}
	}

	client, err := nclient6.New(iface.Attrs().Name,
		nclient6.WithTimeout(c.Timeout),
		nclient6.WithRetry(c.Retries))
	if err != nil {
		return nil, err
	}

	log.Printf("Attempting to get DHCPv6 lease on %s", iface.Attrs().Name)
	p, err := client.RapidSolicit(ctx, dhcpv6.WithNetboot)
	if err != nil {
		return nil, err
	}

	packet := NewPacket6(iface, p)
	log.Printf("Got DHCPv6 lease on %s: %v", iface.Attrs().Name, p.Summary())
	return packet, nil
}

// Result is the result of a particular DHCP attempt.
type Result struct {
	// Interface is the network interface the attempt was sent on.
	Interface netlink.Link

	// Lease is the DHCP configuration returned.
	//
	// If Lease is set, Err is nil.
	Lease Lease

	// Err is an error that occured during the DHCP attempt.
	Err error
}

// SendRequests coordinates soliciting DHCP configuration on all ifs.
//
// ipv4 and ipv6 determine whether to send DHCPv4 and DHCPv6 requests,
// respectively.
//
// The *Result channel will be closed when all requests have completed.
func SendRequests(ctx context.Context, ifs []netlink.Link, ipv4, ipv6 bool, c Config) chan *Result {
	// Yeah, this is a hack, until we can cancel all leases in progress.
	r := make(chan *Result, 3*len(ifs))

	var wg sync.WaitGroup
	for _, iface := range ifs {
		wg.Add(1)
		go func(iface netlink.Link) {
			defer wg.Done()

			log.Printf("Bringing up interface %s...", iface.Attrs().Name)
			if _, err := IfUp(iface.Attrs().Name); err != nil {
				log.Printf("Could not bring up interface %s: %v", iface.Attrs().Name, err)
				return
			}

			if ipv4 {
				wg.Add(1)
				go func(iface netlink.Link) {
					defer wg.Done()
					lease, err := lease4(ctx, iface, c)
					r <- &Result{iface, lease, err}
				}(iface)
			}

			if ipv6 {
				wg.Add(1)
				go func(iface netlink.Link) {
					defer wg.Done()
					lease, err := lease6(ctx, iface, c)
					r <- &Result{iface, lease, err}
				}(iface)
			}
		}(iface)
	}

	go func() {
		wg.Wait()
		close(r)
	}()
	return r
}
