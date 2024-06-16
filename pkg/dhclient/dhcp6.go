// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// Packet6 implements Packet for IPv6 DHCP.
type Packet6 struct {
	p     *dhcpv6.Message
	iface netlink.Link
}

// NewPacket6 wraps a DHCPv6 packet with some convenience methods.
func NewPacket6(iface netlink.Link, p *dhcpv6.Message) *Packet6 {
	return &Packet6{
		p:     p,
		iface: iface,
	}
}

// Message returns the unwrapped DHCPv6 packet.
func (p *Packet6) Message() (*dhcpv4.DHCPv4, *dhcpv6.Message) {
	return nil, p.p
}

// Link returns the interface this packet was received for.
func (p *Packet6) Link() netlink.Link {
	return p.iface
}

// Configure6 adds IPv6 addresses, routes, and DNS servers to the system.
func Configure6(iface netlink.Link, packet *dhcpv6.Message) error {
	p := NewPacket6(iface, packet)
	return p.Configure()
}

// Configure configures interface using this packet.
func (p *Packet6) Configure() error {
	l := p.Lease()
	if l == nil {
		return fmt.Errorf("no lease returned")
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet: &net.IPNet{
			IP: l.IPv6Addr,

			// This mask tells Linux which addresses we know to be
			// "on-link" (i.e., reachable on this interface without
			// having to talk to a router).
			//
			// Since DHCPv6 does not give us that information, we
			// have to assume that no addresses are on-link. To do
			// that, we use /128. (See also RFC 5942 Section 5,
			// "Observed Incorrect Implementation Behavior".)
			Mask: net.CIDRMask(128, 128),
		},
		PreferedLft: int(l.PreferredLifetime.Seconds()),
		ValidLft:    int(l.ValidLifetime.Seconds()),
		// Optimistic DAD (Duplicate Address Detection) means we can
		// use the address before DAD is complete. The DHCP server's
		// job was to give us a unique IP so there is little risk of a
		// collision.
		Flags: unix.IFA_F_OPTIMISTIC,
	}

	if err := netlink.AddrReplace(p.iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %w", dst, p.iface, err)
		}
	}

	if ips := p.DNS(); ips != nil {
		if err := WriteDNSSettings(ips, nil, "", ResolvConfPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet6) String() string {
	if p.Lease() != nil {
		return fmt.Sprintf("IPv6 DHCP Lease IP %s", p.Lease().IPv6Addr)
	}
	return "IPv6 DHCP Lease came with no IP"
}

// Lease returns lease information assigned.
func (p *Packet6) Lease() *dhcpv6.OptIAAddress {
	iana := p.p.Options.OneIANA()
	if iana == nil {
		return nil
	}
	return iana.Options.OneAddress()
}

// DNS returns DNS servers assigned.
func (p *Packet6) DNS() []net.IP {
	return p.p.Options.DNS()
}

// Boot returns the boot file URL and parameters assigned.
func (p *Packet6) Boot() (*url.URL, error) {
	uri := p.p.Options.BootFileURL()
	if len(uri) == 0 {
		return nil, fmt.Errorf("packet does not contain boot file URL")
	}
	return url.Parse(uri)
}

// ISCSIBoot returns the target address and volume name to boot from if
// they were part of the DHCP message.
//
// Parses the DHCPv6 Boot File for iSCSI target and volume as specified by RFC
// 4173 and RFC 5970.
func (p *Packet6) ISCSIBoot() (*net.TCPAddr, string, error) {
	uri := p.p.Options.BootFileURL()
	if len(uri) == 0 {
		return nil, "", fmt.Errorf("packet does not contain boot file URL")
	}
	return ParseISCSIURI(uri)
}
