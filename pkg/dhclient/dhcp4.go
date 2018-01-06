package dhclient

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/u-root/dhcp4"
	dhcp4client "github.com/u-root/dhcp4/client"
	dhcp4opts "github.com/u-root/dhcp4/opts"
	"github.com/vishvananda/netlink"
)

// Packet4 implements Packet for IPv4 DHCP.
type Packet4 struct {
	p *dhcp4.Packet
}

var _ Packet = &Packet4{}

func newPacket4(p *dhcp4.Packet) *Packet4 {
	return &Packet4{
		p: p,
	}
}

// IPs implements Packet.IPs.
func (p *Packet4) IPs() []net.IP {
	return []net.IP{p.p.YIAddr}
}

// Leases implements Packet.Leases.
//
// DHCPv4 only returns one lease at a time.
func (p *Packet4) Leases() []Lease {
	netmask, err := dhcp4opts.GetSubnetMask(p.p.Options)
	if err != nil {
		// If they did not offer a subnet mask, we choose the most
		// restrictive option.
		netmask = []byte{255, 255, 255, 255}
	}

	return []Lease{
		{
			IPNet: &net.IPNet{
				IP:   p.p.YIAddr,
				Mask: net.IPMask(netmask),
			},
		},
	}
}

// Gateway implements Packet.Gateway.
//
// OptionRouter is used as opposed to GIAddr, which seems unused by most DHCP
// servers?
func (p *Packet4) Gateway() net.IP {
	gw, err := dhcp4opts.GetRouters(p.p.Options)
	if err != nil {
		return nil
	}
	return gw[0]
}

// DNS implements Packet.DNS.
func (p *Packet4) DNS() []net.IP {
	ips, err := dhcp4opts.GetDomainNameServers(p.p.Options)
	if err != nil {
		return nil
	}
	return []net.IP(ips)
}

// Boot implements Packet.Boot.
func (p *Packet4) Boot() (url.URL, string, error) {
	// TODO: This is not 100% right -- if a certain option is set, this
	// stuff is encoded in options instead of in the packet's BootFile and
	// ServerName fields.

	// While the default is tftp, servers may specify HTTP or FTP URIs.
	u, err := url.Parse(p.p.BootFile)
	if err != nil {
		return url.URL{}, "", err
	}

	if len(u.Scheme) == 0 {
		// Defaults to tftp is not specified.
		u.Scheme = "tftp"
		u.Path = p.p.BootFile
		if len(p.p.ServerName) == 0 {
			server, err := dhcp4opts.GetServerIdentifier(p.p.Options)
			if err != nil {
				return url.URL{}, "", err
			}
			u.Host = net.IP(server).String()
		} else {
			u.Host = p.p.ServerName
		}
	}
	return *u, "", nil
}

// Client4 implements Client for DHCPv4.
type Client4 struct {
	iface  netlink.Link
	client *dhcp4client.Client
}

var _ Client = &Client4{}

// NewV4 implements a new DHCPv4 client.
func NewV4(iface netlink.Link, timeout time.Duration, retry int) (*Client4, error) {
	ifa, err := net.InterfaceByIndex(iface.Attrs().Index)
	if err != nil {
		return nil, err
	}

	client, err := dhcp4client.New(ifa /*, timeout, retry*/)
	if err != nil {
		return nil, err
	}

	return &Client4{
		iface:  iface,
		client: client,
	}, nil
}

// Solicit implements Client.Solicit.
func (c *Client4) Solicit() (Packet, error) {
	pkt, err := c.client.Request()
	if err != nil {
		return nil, err
	}
	return newPacket4(pkt), nil
}

// Renew implements Client.Renew.
func (c *Client4) Renew(p Packet) (Packet, error) {
	pkt, ok := p.(*Packet4)
	if !ok {
		return nil, fmt.Errorf("passed non-DHCPv4 packet to RenewPacket4")
	}
	pp, err := c.client.Renew(pkt.p)
	if err != nil {
		return nil, err
	}
	return newPacket4(pp), nil
}
