// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/vishvananda/netlink"
)

// Packet4 implements convenience functions for DHCPv4 packets.
type Packet4 struct {
	iface netlink.Link
	P     *dhcpv4.DHCPv4
}

var _ Lease = &Packet4{}

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
	l := p.Lease()
	if l == nil {
		return fmt.Errorf("packet has no IP lease")
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet: l,
	}
	if err := netlink.AddrReplace(p.iface, dst); err != nil {
		return fmt.Errorf("add/replace %s to %v: %v", dst, p.iface, err)
	}

	// RFC 3442 notes that if classless static routes are available, they
	// have priority. You have to ignore the Route Option.
	if routes := p.P.ClasslessStaticRoute(); routes != nil {
		for _, route := range routes {
			r := &netlink.Route{
				LinkIndex: p.iface.Attrs().Index,
				Dst:       route.Dest,
				Gw:        route.Router,
			}
			// If no gateway is specified, the destination must be link-local.
			if r.Gw == nil || r.Gw.Equal(net.IPv4zero) {
				r.Scope = netlink.SCOPE_LINK
			}

			if err := netlink.RouteReplace(r); err != nil {
				return fmt.Errorf("%s: add %s: %v", p.iface.Attrs().Name, r, err)
			}
		}
	} else if gw := p.P.Router(); gw != nil && len(gw) > 0 {
		r := &netlink.Route{
			LinkIndex: p.iface.Attrs().Index,
			Gw:        gw[0],
		}

		if err := netlink.RouteReplace(r); err != nil {
			return fmt.Errorf("%s: add %s: %v", p.iface.Attrs().Name, r, err)
		}
	}

	if ips := p.P.DNS(); ips != nil {
		if err := WriteDNSSettings(ips); err != nil {
			return err
		}
	}
	return nil
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

var (
	ErrNoBootFile       = errors.New("no boot file name present in DHCP message")
	ErrNoServerHostName = errors.New("no server host name present in DHCP message")
)

// Boot returns the boot file assigned.
func (p *Packet4) Boot() (*url.URL, error) {
	// Look for dhcp option presence first, then legacy BootFileName in header.
	bootFileName := p.P.BootFileNameOption()
	bootFileName = strings.TrimRight(bootFileName, "\x00")
	if bootFileName == "" {
		if len(p.P.BootFileName) == 0 {
			return nil, ErrNoBootFile
		}
		bootFileName = p.P.BootFileName
	}

	// While the default is tftp, servers may specify HTTP or FTP URIs.
	u, err := url.Parse(bootFileName)
	if err != nil {
		return nil, err
	}

	if len(u.Scheme) == 0 {
		// Defaults to tftp is not specified.
		u.Scheme = "tftp"
		u.Path = bootFileName
		if len(p.P.ServerHostName) == 0 {
			server := p.P.ServerIdentifier()
			if server != nil {
				u.Host = server.String()
			} else if !p.P.ServerIPAddr.Equal(net.IPv4zero) {
				u.Host = p.P.ServerIPAddr.String()
			} else {
				return nil, ErrNoServerHostName
			}
		} else {
			u.Host = p.P.ServerHostName
		}
	}
	return u, nil
}
