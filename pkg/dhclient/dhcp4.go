// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"net/url"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4opts"
)

// Packet4 implements convenience functions for DHCPv4 packets.
type Packet4 struct {
	P *dhcp4.Packet
}

// NewPacket4 wraps a DHCPv4 packet with some convenience methods.
func NewPacket4(p *dhcp4.Packet) *Packet4 {
	return &Packet4{
		P: p,
	}
}

// Lease returns the IPNet assigned.
func (p *Packet4) Lease() *net.IPNet {
	netmask, err := dhcp4opts.GetSubnetMask(p.P.Options)
	if err != nil {
		// If they did not offer a subnet mask, we choose the most
		// restrictive option.
		netmask = []byte{255, 255, 255, 255}
	}

	return &net.IPNet{
		IP:   p.P.YIAddr,
		Mask: net.IPMask(netmask),
	}
}

// Gateway returns the gateway IP assigned.
//
// OptionRouter is used as opposed to GIAddr, which seems unused by most DHCP
// servers?
func (p *Packet4) Gateway() net.IP {
	gw, err := dhcp4opts.GetRouters(p.P.Options)
	if err != nil {
		return nil
	}
	return gw[0]
}

// DNS returns DNS IPs assigned.
func (p *Packet4) DNS() []net.IP {
	ips, err := dhcp4opts.GetDomainNameServers(p.P.Options)
	if err != nil {
		return nil
	}
	return []net.IP(ips)
}

// Boot returns the boot file assigned.
func (p *Packet4) Boot() (url.URL, error) {
	// TODO: This is not 100% right -- if a certain option is set, this
	// stuff is encoded in options instead of in the packet's BootFile and
	// ServerName fields.

	// While the default is tftp, servers may specify HTTP or FTP URIs.
	u, err := url.Parse(p.P.BootFile)
	if err != nil {
		return url.URL{}, err
	}

	if len(u.Scheme) == 0 {
		// Defaults to tftp is not specified.
		u.Scheme = "tftp"
		u.Path = p.P.BootFile
		if len(p.P.ServerName) == 0 {
			server, err := dhcp4opts.GetServerIdentifier(p.P.Options)
			if err != nil {
				return url.URL{}, err
			}
			u.Host = net.IP(server).String()
		} else {
			u.Host = p.P.ServerName
		}
	}
	return *u, nil
}
