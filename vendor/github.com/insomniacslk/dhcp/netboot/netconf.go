package netboot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/jsimonetti/rtnetlink"
	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/mdlayher/netlink"
)

// AddrConf holds a single IP address configuration for a NIC
type AddrConf struct {
	IPNet             net.IPNet
	PreferredLifetime time.Duration
	ValidLifetime     time.Duration
}

// NetConf holds multiple IP configuration for a NIC, and DNS configuration
type NetConf struct {
	Addresses     []AddrConf
	DNSServers    []net.IP
	DNSSearchList []string
	Routers       []net.IP
}

// GetNetConfFromPacketv6 extracts network configuration information from a DHCPv6
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv6(d *dhcpv6.Message) (*NetConf, error) {
	iana := d.Options.OneIANA()
	if iana == nil {
		return nil, errors.New("no option IA NA found")
	}
	netconf := NetConf{}

	// get IP configuration
	iaaddrs := make([]*dhcpv6.OptIAAddress, 0)
	for _, o := range iana.Options {
		if o.Code() == dhcpv6.OptionIAAddr {
			iaaddrs = append(iaaddrs, o.(*dhcpv6.OptIAAddress))
		}
	}
	netmask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::"))
	for _, iaaddr := range iaaddrs {
		netconf.Addresses = append(netconf.Addresses, AddrConf{
			IPNet: net.IPNet{
				IP:   iaaddr.IPv6Addr,
				Mask: netmask,
			},
			PreferredLifetime: iaaddr.PreferredLifetime,
			ValidLifetime:     iaaddr.ValidLifetime,
		})
	}
	// get DNS configuration
	dns := d.Options.DNS()
	if len(dns) == 0 {
		return nil, errors.New("no option DNS Recursive Name Servers found")
	}
	netconf.DNSServers = dns

	domains := d.Options.DomainSearchList()
	if domains != nil {
		netconf.DNSSearchList = domains.Labels
	}

	return &netconf, nil
}

// GetNetConfFromPacketv4 extracts network configuration information from a DHCPv4
// Reply packet and returns a populated NetConf structure
func GetNetConfFromPacketv4(d *dhcpv4.DHCPv4) (*NetConf, error) {
	// extract the address from the DHCPv4 address
	ipAddr := d.YourIPAddr
	if ipAddr == nil || ipAddr.Equal(net.IPv4zero) {
		return nil, errors.New("ip address is null (0.0.0.0)")
	}
	netconf := NetConf{}

	// get the subnet mask from OptionSubnetMask. If the netmask is not defined
	// in the packet, an error is returned
	netmask := d.SubnetMask()
	if netmask == nil {
		return nil, errors.New("no netmask option in response packet")
	}
	ones, _ := netmask.Size()
	if ones == 0 {
		return nil, errors.New("netmask extracted from OptSubnetMask options is null")
	}

	// netconf struct requires a valid lifetime to be specified. ValidLifetime is a dhcpv6
	// concept, the closest mapping in dhcpv4 world is "IP Address Lease Time". If the lease
	// time option is nil, we set it to 0
	leaseTime := d.IPAddressLeaseTime(0)

	netconf.Addresses = append(netconf.Addresses, AddrConf{
		IPNet: net.IPNet{
			IP:   ipAddr,
			Mask: netmask,
		},
		PreferredLifetime: 0,
		ValidLifetime:     leaseTime,
	})

	// get DNS configuration
	dnsServers := d.DNS()
	if len(dnsServers) == 0 {
		return nil, errors.New("no dns servers options in response packet")
	}
	netconf.DNSServers = dnsServers

	// get domain search list
	dnsSearchList := d.DomainSearch()
	if dnsSearchList != nil {
		if len(dnsSearchList.Labels) == 0 {
			return nil, errors.New("dns search list is empty")
		}
		netconf.DNSSearchList = dnsSearchList.Labels
	}

	// get default gateway
	routersList := d.Router()
	if len(routersList) == 0 {
		return nil, errors.New("no routers specified in the corresponding option")
	}
	netconf.Routers = routersList
	return &netconf, nil
}

// IfUp brings up an interface by name, and waits for it to come up until a timeout expires
func IfUp(ifname string, timeout time.Duration) (_ *net.Interface, err error) {
	start := time.Now()
	rt, err := rtnl.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rt.Close(); cerr != nil {
			err = cerr
		}
	}()
	for time.Since(start) < timeout {
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			return nil, err
		}
		// If the interface is up, return. According to kernel documentation OperState may
		// be either Up or Unknown:
		//   Interface is in RFC2863 operational state UP or UNKNOWN. This is for
		//   backward compatibility, routing daemons, dhcp clients can use this
		//   flag to determine whether they should use the interface.
		// Source: https://www.kernel.org/doc/Documentation/networking/operstates.txt
		operState, err := getOperState(iface.Index)
		if err != nil {
			return nil, err
		}
		if operState == rtnetlink.OperStateUp || operState == rtnetlink.OperStateUnknown {
			// XXX despite the OperUp state, upon the first attempt I
			// consistently get a "cannot assign requested address" error. Need
			// to investigate more.
			time.Sleep(time.Second)
			return iface, nil
		}
		// otherwise try to bring it up
		if err := rt.LinkUp(iface); err != nil {
			return nil, fmt.Errorf("interface %q: %v can't bring it up: %v", ifname, iface, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil, fmt.Errorf("timed out while waiting for %s to come up", ifname)

}

// ConfigureInterface configures a network interface with the configuration held by a
// NetConf structure
func ConfigureInterface(ifname string, netconf *NetConf) (err error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return err
	}
	rt, err := rtnl.Dial(nil)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := rt.Close(); err != nil {
			err = cerr
		}
	}()
	// configure interfaces
	for _, addr := range netconf.Addresses {
		if err := rt.AddrAdd(iface, &addr.IPNet); err != nil {
			return fmt.Errorf("cannot configure %s on %s: %v", ifname, addr.IPNet, err)
		}
	}
	// configure /etc/resolv.conf
	resolvconf := ""
	for _, ns := range netconf.DNSServers {
		resolvconf += fmt.Sprintf("nameserver %s\n", ns)
	}
	if len(netconf.DNSSearchList) > 0 {
		resolvconf += fmt.Sprintf("search %s\n", strings.Join(netconf.DNSSearchList, " "))
	}
	if err = ioutil.WriteFile("/etc/resolv.conf", []byte(resolvconf), 0644); err != nil {
		return fmt.Errorf("could not write resolv.conf file %v", err)
	}

	// FIXME wut? No IPv6 here?
	// add default route information for v4 space. only one default route is allowed
	// so ignore the others if there are multiple ones
	if len(netconf.Routers) > 0 {
		// if there is a default v4 route, remove it, as we want to add the one we just got during
		// the dhcp transaction. if the route is not present, which is the final state we want,
		// an error is returned so ignore it
		dst := net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.CIDRMask(0, 32),
		}
		// Remove a possible default route (dst 0.0.0.0) to the L2 domain (gw: 0.0.0.0), which is what
		// a client would want to add before initiating the DHCP transaction in order not to fail with
		// ENETUNREACH. If this default route has a specific metric assigned, it doesn't get removed.
		// The code doesn't remove any other default route (i.e. gw != 0.0.0.0).
		if err := rt.RouteDel(iface, net.IPNet{IP: net.IPv4zero}); err != nil {
			switch err := err.(type) {
			case *netlink.OpError:
				// ignore the error if it's -EEXIST or -ESRCH
				if !os.IsExist(err.Err) && err.Err != syscall.ESRCH {
					return fmt.Errorf("could not delete default route on interface %s: %v", ifname, err)
				}
			default:
				return fmt.Errorf("could not delete default route on interface %s: %v", ifname, err)
			}
		}

		src := netconf.Addresses[0].IPNet
		// TODO handle the remaining Routers if more than one
		if err := rt.RouteAddSrc(iface, dst, &src, netconf.Routers[0]); err != nil {
			return fmt.Errorf("could not add gateway %s for src %s dst %s to interface %s: %v", netconf.Routers[0], src, dst, ifname, err)
		}
	}

	return nil
}
