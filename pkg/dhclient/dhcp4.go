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
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/vishvananda/netlink"
)

// DefaultScheme for boot file if there are none in the lease
var DefaultScheme = "tftp"

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

// Message returns the unwrapped DHCPv4 packet.
func (p *Packet4) Message() (*dhcpv4.DHCPv4, *dhcpv6.Message) {
	return p.P, nil
}

// Link is a netlink link
func (p *Packet4) Link() netlink.Link {
	return p.iface
}

// GatherDNSSettings gets the DNS related infromation from a dhcp packet
// including, nameservers, domain, and search options
func (p *Packet4) GatherDNSSettings() (ns []net.IP, sl []string, dom string) {
	if nameservers := p.P.DNS(); nameservers != nil {
		ns = nameservers
	}
	if searchList := p.P.DomainSearch(); searchList != nil {
		sl = searchList.Labels
	}
	if domain := p.P.DomainName(); domain != "" {
		dom = domain
	}
	return
}

// Configure4 adds IP addresses, routes, and DNS servers to the system.
func Configure4(iface netlink.Link, packet *dhcpv4.DHCPv4) error {
	p := NewPacket4(iface, packet)
	return p.Configure()
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
		return fmt.Errorf("add/replace %s to %v: %w", dst, p.iface, err)
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
				return fmt.Errorf("%s: add %s: %w", p.iface.Attrs().Name, r, err)
			}
		}
	} else if gw := p.P.Router(); len(gw) > 0 {
		r := &netlink.Route{
			LinkIndex: p.iface.Attrs().Index,
			Gw:        gw[0],
		}

		if err := netlink.RouteReplace(r); err != nil {
			return fmt.Errorf("%s: add %s: %w", p.iface.Attrs().Name, r, err)
		}
	}

	nameServers, searchList, domain := p.GatherDNSSettings()
	if err := WriteDNSSettings(nameServers, searchList, domain, ResolvConfPath); err != nil {
		return err
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
	// ErrNoBootFile represents that no pxe boot file was found.
	ErrNoBootFile = errors.New("no boot file name present in DHCP message")

	// ErrNoRootPath means no root path option was found in DHCP message.
	ErrNoRootPath = errors.New("no root path in DHCP message")

	// ErrNoServerHostName represents that no pxe boot server was found.
	ErrNoServerHostName = errors.New("no server host name present in DHCP message")
)

func (p *Packet4) bootfilename() string {
	// Look for dhcp option presence first, then legacy BootFileName in header.
	bootFileName := p.P.BootFileNameOption()
	bootFileName = strings.TrimRight(bootFileName, "\x00")
	if len(bootFileName) > 0 {
		return bootFileName
	}
	if len(p.P.BootFileName) >= 0 {
		return p.P.BootFileName
	}
	return ""
}

// Boot returns the boot file assigned.
func (p *Packet4) Boot() (*url.URL, error) {
	bootFileName := p.bootfilename()
	if len(bootFileName) == 0 {
		return nil, ErrNoBootFile
	}

	// While the default is tftp, servers may specify HTTP or FTP URIs.
	u, err := url.Parse(bootFileName)
	if err != nil {
		return nil, err
	}

	if len(u.Scheme) == 0 {
		// Use the DefaultScheme if not specified
		u.Scheme = DefaultScheme
		u.Path = bootFileName
		if len(p.P.ServerHostName) == 0 {
			server := p.P.ServerIdentifier()
			if !p.P.ServerIPAddr.Equal(net.IPv4zero) {
				u.Host = p.P.ServerIPAddr.String()
			} else if server != nil {
				u.Host = server.String()
			} else {
				return nil, ErrNoServerHostName
			}
		} else {
			u.Host = p.P.ServerHostName
		}
	}
	return u, nil
}

// ISCSIBoot returns the target address and volume name to boot from if
// they were part of the DHCP message.
//
// Parses the IPv4 DHCP Root Path for iSCSI target and volume as specified by
// RFC 4173.
func (p *Packet4) ISCSIBoot() (*net.TCPAddr, string, error) {
	rp := p.P.RootPath()
	if len(rp) > 0 {
		return ParseISCSIURI(rp)
	}
	bootfilename := p.bootfilename()
	if len(bootfilename) > 0 && strings.HasPrefix(bootfilename, "iscsi:") {
		return ParseISCSIURI(bootfilename)
	}
	return nil, "", ErrNoRootPath
}

// Response returns the DHCP response
func (p *Packet4) Response() interface{} {
	return p.P
}
