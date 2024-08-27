// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dhclient allows for getting both DHCPv4 and DHCPv6 leases on
// multiple network interfaces in parallel.
package dhclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

const (
	// ResolvConfPath is the path to the resolv.conf file.
	ResolvConfPath = "/etc/resolv.conf"
)

// isIPv6LinkReady returns true if the interface has a link-local address
// which is not tentative.
func isIPv6LinkReady(l netlink.Link) (bool, error) {
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

// isIPv6RouteReady returns true if serverAddr is reachable.
func isIPv6RouteReady(l netlink.Link, serverAddr net.IP) (bool, error) {
	if serverAddr.IsMulticast() {
		return true, nil
	}

	routes, err := netlink.RouteList(l, netlink.FAMILY_V6)
	if err != nil {
		return false, err
	}
	for _, route := range routes {
		if route.LinkIndex != l.Attrs().Index {
			continue
		}
		// Default route.
		if route.Dst == nil {
			return true, nil
		}
		if route.Dst.Contains(serverAddr) {
			return true, nil
		}
	}
	return false, nil
}

// IfUp ensures the given network interface is up and returns the link object.
func IfUp(ifname string, linkUpTimeout time.Duration) (netlink.Link, error) {
	start := time.Now()
	for time.Since(start) < linkUpTimeout {
		// Note that it may seem odd to keep trying the LinkByName
		// operation, but consider that a hotplug device such as USB
		// ethernet can just vanish.
		iface, err := netlink.LinkByName(ifname)
		if err != nil {
			return nil, fmt.Errorf("cannot get interface %q by name: %w", ifname, err)
		}

		// Check if link is actually operational.
		// https://www.kernel.org/doc/Documentation/networking/operstates.txt states that we should check
		// for OperUp and OperUnknown.
		if o := iface.Attrs().OperState; o == netlink.OperUp || o == netlink.OperUnknown {
			return iface, nil
		}

		if err := netlink.LinkSetUp(iface); err != nil {
			return nil, fmt.Errorf("interface %q: %v can't make it up: %w", ifname, iface, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("link %q still down after %v seconds", ifname, linkUpTimeout.Seconds())
}

// WriteDNSSettings writes the given nameservers, search list, and domain to the resolv.conf at the specified path.
func WriteDNSSettings(ns []net.IP, sl []string, domain, resolvConfPath string) error {
	rc := &bytes.Buffer{}
	if domain != "" {
		rc.WriteString(fmt.Sprintf("domain %s\n", domain))
	}
	for _, ip := range ns {
		rc.WriteString(fmt.Sprintf("nameserver %s\n", ip))
	}
	if sl != nil {
		rc.WriteString("search ")
		rc.WriteString(strings.Join(sl, " "))
		rc.WriteString("\n")
	}
	return os.WriteFile(resolvConfPath, rc.Bytes(), 0o644)
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

	// ISCSIBoot returns the target address and volume name to boot from if
	// they were part of the DHCP message.
	ISCSIBoot() (*net.TCPAddr, string, error)

	// Link is the interface the configuration is for.
	Link() netlink.Link

	// Return the full DHCP response, this is either a *dhcpv4.DHCPv4 or a
	// *dhcpv6.Message.
	Message() (*dhcpv4.DHCPv4, *dhcpv6.Message)
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

	// Modifiers4 allows modifications to the IPv4 DHCP request.
	Modifiers4 []dhcpv4.Modifier

	// Modifiers6 allows modifications to the IPv6 DHCP request.
	Modifiers6 []dhcpv6.Modifier

	// V6ServerAddr can be a unicast or broadcast destination for DHCPv6
	// messages.
	//
	// If not set, it will default to nclient6's default (all servers &
	// relay agents).
	V6ServerAddr *net.UDPAddr

	// V6ClientPort is the port that is used to send and receive DHCPv6
	// messages.
	//
	// If not set, it will default to dhcpv6's default (546).
	V6ClientPort *int

	// V4ServerAddr can be a unicast or broadcast destination for IPv4 DHCP
	// messages.
	//
	// If not set, it will default to nclient4's default (DHCP broadcast
	// address).
	V4ServerAddr *net.UDPAddr

	// If true, add Client Identifier (61) option to the IPv4 request.
	V4ClientIdentifier bool
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
	if c.V4ServerAddr != nil {
		mods = append(mods, nclient4.WithServerAddr(c.V4ServerAddr))
	}
	client, err := nclient4.New(iface.Attrs().Name, mods...)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Prepend modifiers with default options, so they can be overriden.
	reqmods := append(
		[]dhcpv4.Modifier{
			dhcpv4.WithOption(dhcpv4.OptClassIdentifier("PXE UROOT")),
			dhcpv4.WithRequestedOptions(dhcpv4.OptionSubnetMask),
			dhcpv4.WithNetboot,
		},
		c.Modifiers4...)

	if c.V4ClientIdentifier {
		// Client Id is hardware type + mac per RFC 2132 9.14.
		ident := []byte{0x01} // Type ethernet
		ident = append(ident, iface.Attrs().HardwareAddr...)
		reqmods = append(reqmods, dhcpv4.WithOption(dhcpv4.OptClientIdentifier(ident)))
	}

	log.Printf("Attempting to get DHCPv4 lease on %s", iface.Attrs().Name)
	lease, err := client.Request(ctx, reqmods...)
	if err != nil {
		return nil, err
	}

	packet := NewPacket4(iface, lease.ACK)
	log.Printf("Got DHCPv4 lease on %s: %v", iface.Attrs().Name, lease.ACK.Summary())
	return packet, nil
}

func lease6(ctx context.Context, iface netlink.Link, c Config, linkUpTimeout time.Duration) (Lease, error) {
	clientPort := dhcpv6.DefaultClientPort
	if c.V6ClientPort != nil {
		clientPort = *c.V6ClientPort
	}

	// For ipv6, we cannot bind to the port until Duplicate Address
	// Detection (DAD) is complete which is indicated by the link being no
	// longer marked as "tentative". This usually takes about a second.

	// If the link is never going to be ready, don't wait forever.
	// (The user may not have configured a ctx with a timeout.)
	//
	// Hardcode the timeout to 30s for now.
	linkTimeout := time.After(linkUpTimeout)
	for {
		if ready, err := isIPv6LinkReady(iface); err != nil {
			return nil, err
		} else if ready {
			break
		}
		select {
		case <-time.After(100 * time.Millisecond):
			continue
		case <-linkTimeout:
			return nil, errors.New("timeout after waiting for a non-tentative IPv6 address")
		case <-ctx.Done():
			return nil, errors.New("timeout after waiting for a non-tentative IPv6 address")
		}
	}

	// If user specified a non-multicast address, make sure it's routable before we start.
	if c.V6ServerAddr != nil {
		for {
			if ready, err := isIPv6RouteReady(iface, c.V6ServerAddr.IP); err != nil {
				return nil, err
			} else if ready {
				break
			}
			select {
			case <-time.After(100 * time.Millisecond):
				continue
			case <-linkTimeout:
				return nil, errors.New("timeout after waiting for a route")
			case <-ctx.Done():
				return nil, errors.New("timeout after waiting for a route")
			}
		}
	}

	mods := []nclient6.ClientOpt{
		nclient6.WithTimeout(c.Timeout),
		nclient6.WithRetry(c.Retries),
	}
	switch c.LogLevel {
	case LogSummary:
		mods = append(mods, nclient6.WithSummaryLogger())
	case LogDebug:
		mods = append(mods, nclient6.WithDebugLogger())
	}
	if c.V6ServerAddr != nil {
		mods = append(mods, nclient6.WithBroadcastAddr(c.V6ServerAddr))
	}
	conn, err := nclient6.NewIPv6UDPConn(iface.Attrs().Name, clientPort)
	if err != nil {
		return nil, err
	}
	i, err := net.InterfaceByName(iface.Attrs().Name)
	if err != nil {
		return nil, err
	}
	client, err := nclient6.NewWithConn(conn, i.HardwareAddr, mods...)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Prepend modifiers with default options, so they can be overriden.
	reqmods := append(
		[]dhcpv6.Modifier{
			dhcpv6.WithNetboot,
		},
		c.Modifiers6...)

	log.Printf("Attempting to get DHCPv6 lease on %s", iface.Attrs().Name)
	p, err := client.RapidSolicit(ctx, reqmods...)
	if err != nil {
		return nil, err
	}

	packet := NewPacket6(iface, p)
	log.Printf("Got DHCPv6 lease on %s: %v", iface.Attrs().Name, p.Summary())
	return packet, nil
}

// NetworkProtocol is either IPv4 or IPv6.
type NetworkProtocol int

// Possible network protocols; either IPv4, IPv6, or both.
const (
	NetIPv4 NetworkProtocol = 1
	NetIPv6 NetworkProtocol = 2
	NetBoth NetworkProtocol = 3
)

func (n NetworkProtocol) String() string {
	switch n {
	case NetIPv4:
		return "IPv4"
	case NetIPv6:
		return "IPv6"
	case NetBoth:
		return "IPv4+IPv6"
	}
	return fmt.Sprintf("unknown network protocol (%#x)", n)
}

// Result is the result of a particular DHCP attempt.
type Result struct {
	// Protocol is the IP protocol that we tried to configure.
	Protocol NetworkProtocol

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
func SendRequests(ctx context.Context, ifs []netlink.Link, ipv4, ipv6 bool, c Config, linkUpTimeout time.Duration) chan *Result {
	// Yeah, this is a hack, until we can cancel all leases in progress.
	r := make(chan *Result, 3*len(ifs))

	var wg sync.WaitGroup
	for _, iface := range ifs {
		wg.Add(1)
		go func(iface netlink.Link) {
			defer wg.Done()

			log.Printf("Bringing up interface %s...", iface.Attrs().Name)
			if _, err := IfUp(iface.Attrs().Name, linkUpTimeout); err != nil {
				log.Printf("Could not bring up interface %s: %v", iface.Attrs().Name, err)
				return
			}

			if ipv4 {
				wg.Add(1)
				go func(iface netlink.Link) {
					defer wg.Done()
					lease, err := lease4(ctx, iface, c)
					r <- &Result{NetIPv4, iface, lease, err}
				}(iface)
			}

			if ipv6 {
				wg.Add(1)
				go func(iface netlink.Link) {
					defer wg.Done()
					lease, err := lease6(ctx, iface, c, linkUpTimeout)
					r <- &Result{NetIPv6, iface, lease, err}
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
