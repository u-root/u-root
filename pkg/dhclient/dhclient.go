// Package dhclient provides a unified interface for interfacing with both
// DHCPv4 and DHCPv6 clients.
package dhclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/dhcp6opts"
	"github.com/u-root/dhcp4"
	"github.com/vishvananda/netlink"
)

const linkUpAttempt = 30 * time.Second

// IfUp ensures the given network interface is up and returns the link object.
func IfUp(ifname string) (netlink.Link, error) {
	start := time.Now()
	for time.Since(start) < linkUpAttempt {
		// Note that it may seem odd to keep trying the LinkByName
		// operation, but consider that a hotplug device such as USB
		// ethernet can just vanish.
		iface, err := netlink.LinkByName(ifname)
		if err != nil {
			return nil, fmt.Errorf("cannot get interface %q by name: %v", ifname, err)
		}

		if iface.Attrs().OperState == netlink.OperUp {
			return iface, nil
		}

		if err := netlink.LinkSetUp(iface); err != nil {
			return nil, fmt.Errorf("interface %q: %v can't make it up: %v", ifname, iface, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("link %q still down after %d seconds", ifname, linkUpAttempt)
}

// Configure4 adds IP addresses, routes, and DNS servers to the system.
func Configure4(iface netlink.Link, packet *dhcp4.Packet) error {
	p := NewPacket4(packet)

	l := p.Lease()
	if l == nil {
		return fmt.Errorf("no lease returned")
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet: l,
	}
	if err := netlink.AddrReplace(iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %v", dst, iface, err)
		}
	}

	if gw := p.Gateway(); gw != nil {
		r := &netlink.Route{
			LinkIndex: iface.Attrs().Index,
			Gw:        gw,
		}

		if err := netlink.RouteReplace(r); err != nil {
			return fmt.Errorf("%s: add %s: %v", iface.Attrs().Name, r, err)
		}
	}

	if ips := p.DNS(); ips != nil {
		if err := WriteDNSSettings(ips); err != nil {
			return err
		}
	}
	return nil
}

// Configure6 adds IPv6 addresses, routes, and DNS servers to the system.
func Configure6(iface netlink.Link, packet *dhcp6.Packet, iana *dhcp6opts.IANA) error {
	p := NewPacket6(packet, iana)

	l := p.Lease()
	if l == nil {
		return fmt.Errorf("no lease returned")
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   l.IP,
			Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::")),
		},
		PreferedLft: int(l.PreferredLifetime.Seconds()),
		ValidLft:    int(l.ValidLifetime.Seconds()),
	}
	if err := netlink.AddrReplace(iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %v", dst, iface, err)
		}
	}

	if ips := p.DNS(); ips != nil {
		if err := WriteDNSSettings(ips); err != nil {
			return err
		}
	}
	return nil
}

// WriteDNSSettings writes the given IPs as nameservers to resolv.conf.
func WriteDNSSettings(ips []net.IP) error {
	rc := &bytes.Buffer{}
	for _, ip := range ips {
		rc.WriteString(fmt.Sprintf("nameserver %s\n", ip))
	}
	return ioutil.WriteFile("/etc/resolv.conf", rc.Bytes(), 0644)
}
