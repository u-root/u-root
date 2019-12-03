// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"net"
	"net/url"
	"os"

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

// Message returns the wrapped DHCPv6 packet.
func (p *Packet6) Message() *dhcpv6.Message {
	return p.p
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
	if err := netlink.AddrReplace(p.iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %v", dst, p.iface, err)
		}
	}

	if ips := p.DNS(); ips != nil {
		if err := WriteDNSSettings(ips, nil, ""); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet6) String() string {
	return fmt.Sprintf("IPv6 DHCP Lease IP %s", p.Lease().IPv6Addr)
}

// Lease returns lease information assigned.
func (p *Packet6) Lease() *dhcpv6.OptIAAddress {
	// TODO(chrisko): Reform dhcpv6 option handling to be like dhcpv4.
	ianaOpt := p.p.GetOneOption(dhcpv6.OptionIANA)
	iana, ok := ianaOpt.(*dhcpv6.OptIANA)
	if !ok {
		return nil
	}

	iaAddrOpt := iana.Options.GetOne(dhcpv6.OptionIAAddr)
	iaAddr, ok := iaAddrOpt.(*dhcpv6.OptIAAddress)
	if !ok {
		return nil
	}
	return iaAddr
}

// DNS returns DNS servers assigned.
func (p *Packet6) DNS() []net.IP {
	// TODO: Would the IANA contain this, or the packet?
	dnsOpt := p.p.GetOneOption(dhcpv6.OptionDNSRecursiveNameServer)
	dns, ok := dnsOpt.(*dhcpv6.OptDNSRecursiveNameServer)
	if !ok {
		return nil
	}
	return dns.NameServers
}

// Boot returns the boot file URL and parameters assigned.
//
// TODO: RFC 5970 is helpfully avoidant of where these options are used. Are
// they added to the packet? Are they added to an IANA?  It *seems* like it's
// in the packet.
func (p *Packet6) Boot() (*url.URL, error) {
	uriOpt := p.p.GetOneOption(dhcpv6.OptionBootfileURL)
	uri, ok := uriOpt.(dhcpv6.OptBootFileURL)
	if !ok {
		return nil, fmt.Errorf("packet does not contain boot file URL")
	}
	// Srsly, a []byte?
	return url.Parse(string(uri.ToBytes()))
}

// ISCSIBoot returns the target address and volume name to boot from if
// they were part of the DHCP message.
//
// Parses the DHCPv6 Boot File for iSCSI target and volume as specified by RFC
// 4173 and RFC 5970.
func (p *Packet6) ISCSIBoot() (*net.TCPAddr, string, error) {
	uriOpt := p.p.GetOneOption(dhcpv6.OptionBootfileURL)
	uri, ok := uriOpt.(dhcpv6.OptBootFileURL)
	if !ok {
		return nil, "", fmt.Errorf("packet does not contain boot file URL")
	}
	return parseISCSIURI(string(uri.ToBytes()))
}
