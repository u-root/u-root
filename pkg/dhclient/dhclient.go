// Package dhclient provides a unified interface for interfacing with both
// DHCPv4 and DHCPv6 clients.
package dhclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"time"

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

// Lease is a DHCP lease.
type Lease struct {
	IPNet             *net.IPNet
	PreferredLifetime time.Duration
	ValidLifetime     time.Duration
}

// Packet is a DHCP packet.
type Packet interface {
	IPs() []net.IP

	// Leases are the leases returned by the DHCP server if specified,
	// otherwise nil.
	Leases() []Lease

	// Gateway is the gateway server, if specified, otherwise nil.
	Gateway() net.IP

	// DNS are DNS server addresses, if specified, otherwise nil.
	DNS() []net.IP

	// Boot returns the boot URI and boot parameters if specified,
	// otherwise an error.
	Boot() (url.URL, string, error)
}

// Client is a DHCP client.
type Client interface {
	// Solicit solicits a new DHCP lease.
	Solicit() (Packet, error)

	// Renew renews an existing DHCP lease.
	Renew(p Packet) (Packet, error)
}

// HandlePacket adds IP addresses, routes, and DNS servers to the system.
func HandlePacket(iface netlink.Link, packet Packet) error {
	l := packet.Leases()

	// We currently only know how to handle one lease.
	if len(l) > 1 {
		log.Printf("interface %s: only handling one lease.", iface.Attrs().Name)
	}

	// Add the address to the iface.
	dst := &netlink.Addr{
		IPNet:       l[0].IPNet,
		PreferedLft: int(l[0].PreferredLifetime.Seconds()),
		ValidLft:    int(l[0].ValidLifetime.Seconds()),
	}
	if err := netlink.AddrReplace(iface, dst); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("add/replace %s to %v: %v", dst, iface, err)
		}
	}

	if gw := packet.Gateway(); gw != nil {
		r := &netlink.Route{
			LinkIndex: iface.Attrs().Index,
			Gw:        gw,
		}

		if err := netlink.RouteReplace(r); err != nil {
			return fmt.Errorf("%s: add %s: %v", iface.Attrs().Name, r, err)
		}
	}

	if ips := packet.DNS(); ips != nil {
		rc := &bytes.Buffer{}
		for _, ip := range ips {
			rc.WriteString(fmt.Sprintf("nameserver %s\n", ip))
		}
		if err := ioutil.WriteFile("/etc/resolv.conf", rc.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}
