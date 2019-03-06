// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"net/url"
	"strings"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/dhcp6opts"
)

// Packet6 implements Packet for IPv6 DHCP.
type Packet6 struct {
	p    *dhcp6.Packet
	iana *dhcp6opts.IANA
}

// NewPacket6 wraps a DHCPv6 packet with some convenience methods.
func NewPacket6(p *dhcp6.Packet, ianaLease *dhcp6opts.IANA) *Packet6 {
	return &Packet6{
		p:    p,
		iana: ianaLease,
	}
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
func (p *Packet6) Boot() (*url.URL, string, error) {
	uri, err := dhcp6opts.GetBootFileURL(p.p.Options)
	if err != nil {
		return nil, "", err
	}

	// Having this value is optional.
	bfp, err := dhcp6opts.GetBootFileParam(p.p.Options)
	if err != dhcp6.ErrOptionNotPresent {
		return nil, "", err
	}

	var cmdline string
	if bfp != nil {
		cmdline = strings.Join(bfp, " ")
	}
	return (*url.URL)(uri), cmdline, nil
}
