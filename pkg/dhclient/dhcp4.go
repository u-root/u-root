// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"net"
	"net/url"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/vishvananda/netlink"
)

// Packet4 implements convenience functions for DHCPv4 packets.
type Packet4 struct {
	iface netlink.Link
	P     *dhcpv4.DHCPv4
}

// NewPacket4 wraps a DHCPv4 packet with some convenience methods.
func NewPacket4(iface netlink.Link, p *dhcpv4.DHCPv4) *Packet4 {
	return &Packet4{
		iface: iface,
		P:     p,
	}
}

func (p *Packet4) Link() netlink.Link {
	return p.iface
}

// Configure configures interface using this packet.
func (p *Packet4) Configure() error {
	return Configure4(p.iface, p.P)
}

func (p *Packet4) String() string {
	return fmt.Sprintf("IPv4 DHCP Lease IP %s", p.Lease())
}

// Lease returns the IPNet assigned.
func (p *Packet4) Lease() *net.IPNet {
	netmask := p.P.SubnetMask()
	if netmask == nil {
		// If they did not offer a subnet mask, we choose the most
		// restrictive option.
		netmask = []byte{255, 255, 255, 255}
	}

	return &net.IPNet{
		IP:   p.P.YourIPAddr,
		Mask: net.IPMask(netmask),
	}
}

// Gateway returns the gateway IP assigned.
//
// OptionRouter is used as opposed to GIAddr, which seems unused by most DHCP
// servers?
func (p *Packet4) Gateway() net.IP {
	gw := p.P.Router()
	if gw == nil {
		return nil
	}
	return gw[0]
}

// DNS returns DNS IPs assigned.
func (p *Packet4) DNS() []net.IP {
	return p.P.DNS()
}

// Boot returns the boot file assigned.
func (p *Packet4) Boot() (*url.URL, error) {
	// TODO: This is not 100% right -- if a certain option is set, this
	// stuff is encoded in options instead of in the packet's BootFile and
	// ServerName fields.

	// While the default is tftp, servers may specify HTTP or FTP URIs.
	u, err := url.Parse(p.P.BootFileName)
	if err != nil {
		return nil, err
	}

	if len(u.Scheme) == 0 {
		// Defaults to tftp is not specified.
		u.Scheme = "tftp"
		u.Path = p.P.BootFileName
		if len(p.P.ServerHostName) == 0 {
			server := p.P.ServerIdentifier()
			if server == nil {
				return nil, err
			}
			u.Host = server.String()
		} else {
			u.Host = p.P.ServerHostName
		}
	}
	return u, nil
}
