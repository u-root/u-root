// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"net/url"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/dhcp6opts"
	"github.com/vishvananda/netlink"
)

// Packet6 implements Packet for IPv6 DHCP.
type Packet6 struct {
	p     *dhcp6.Packet
	iana  *dhcp6opts.IANA
	iface netlink.Link
}

// NewPacket6 wraps a DHCPv6 packet with some convenience methods.
func NewPacket6(iface netlink.Link, p *dhcp6.Packet, ianaLease *dhcp6opts.IANA) *Packet6 {
	return &Packet6{
		p:     p,
		iface: iface,
		iana:  ianaLease,
	}
}

func (p *Packet6) Link() netlink.Link {
	return p.iface
}

// Configure configures interface using this packet.
func (p *Packet6) Configure() error {
	return Configure6(p.iface, p.p, p.iana)
}

// Lease returns lease information assigned.
func (p *Packet6) Lease() *dhcp6opts.IAAddr {
	// TODO: Can a DHCPv6 server return multiple IAAddrs for one IANA?
	// There certainly doesn't seem to be a way to request multiple other
	// than requesting multiple IANAs.
	iaAddrs, err := dhcp6opts.GetIAAddr(p.iana.Options)
	if err != nil || len(iaAddrs) == 0 {
		return nil
	}

	return iaAddrs[0]
}

// DNS returns DNS servers assigned.
func (p *Packet6) DNS() []net.IP {
	// TODO: Would the IANA contain this, or the packet?
	ips, err := dhcp6opts.GetDNSServers(p.p.Options)
	if err != nil {
		return nil
	}
	return []net.IP(ips)
}

// Boot returns the boot file URL and parameters assigned.
//
// TODO: RFC 5970 is helpfully avoidant of where these options are used. Are
// they added to the packet? Are they added to an IANA?  It *seems* like it's
// in the packet.
func (p *Packet6) Boot() (*url.URL, error) {
	uri, err := dhcp6opts.GetBootFileURL(p.p.Options)
	if err != nil {
		return nil, err
	}
	return (*url.URL)(uri), nil
}
