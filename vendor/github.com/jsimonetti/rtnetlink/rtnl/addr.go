package rtnl

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/jsimonetti/rtnetlink"
)

// AddrAdd associates an IP-address with an interface.
func (c *Conn) AddrAdd(ifc *net.Interface, addr *net.IPNet) error {
	af, err := addrFamily(addr.IP)
	if err != nil {
		return err
	}
	scope := addrScope(addr.IP)
	prefixlen, _ := addr.Mask.Size()
	tx := &rtnetlink.AddressMessage{
		Family:       uint8(af),
		PrefixLength: uint8(prefixlen),
		Scope:        uint8(scope),
		Index:        uint32(ifc.Index),
		Attributes: rtnetlink.AddressAttributes{
			Address: addr.IP,
			Local:   addr.IP,
		},
	}
	if ip4 := addr.IP.To4(); ip4 != nil {
		tx.Attributes.Broadcast = broadcastAddr(addr)
		tx.Attributes.Address = ip4
		tx.Attributes.Local = ip4
	}
	return c.Conn.Address.New(tx)
}

// AddrDel revokes an IP-address from an interface.
func (c *Conn) AddrDel(ifc *net.Interface, addr *net.IPNet) error {
	af, err := addrFamily(addr.IP)
	if err != nil {
		return err
	}
	prefixlen, _ := addr.Mask.Size()
	rx, err := c.Addrs(ifc, af)
	if err != nil {
		return err
	}
	for _, v := range rx {
		plen, _ := v.Mask.Size()
		if plen == prefixlen && v.IP.Equal(addr.IP) {
			tx := &rtnetlink.AddressMessage{
				Family:       uint8(af),
				PrefixLength: uint8(prefixlen),
				Index:        uint32(ifc.Index),
				Attributes: rtnetlink.AddressAttributes{
					Address: addr.IP,
				},
			}
			if ip4 := addr.IP.To4(); ip4 != nil {
				tx.Attributes.Broadcast = broadcastAddr(addr)
				tx.Attributes.Address = ip4
			}
			return c.Conn.Address.Delete(tx)
		}
	}
	return &net.AddrError{Err: "address not found", Addr: addr.IP.String()}
}

// Addrs returns IP addresses matching the interface and address family.
// To retrieve all addresses configured for the system, run:
//
//	conn.Addrs(nil, 0)
//
func (c *Conn) Addrs(ifc *net.Interface, family int) (out []*net.IPNet, err error) {
	rx, err := c.Conn.Address.List()
	if err != nil {
		return nil, err
	}
	match := func(v *rtnetlink.AddressMessage, ifc *net.Interface, family int) bool {
		if ifc != nil && v.Index != uint32(ifc.Index) {
			return false
		}
		if family != 0 && v.Family != uint8(family) {
			return false
		}
		return true
	}
	for _, m := range rx {
		if match(&m, ifc, family) {
			bitlen := 8 * len(m.Attributes.Address)
			a := &net.IPNet{
				IP:   m.Attributes.Address,
				Mask: net.CIDRMask(int(m.PrefixLength), bitlen),
			}
			out = append(out, a)
		}
	}
	return
}

func addrFamily(ip net.IP) (int, error) {
	if ip.To4() != nil {
		return unix.AF_INET, nil
	}
	if len(ip) == net.IPv6len {
		return unix.AF_INET6, nil
	}
	return 0, &net.AddrError{Err: "invalid IP address", Addr: ip.String()}
}

func addrScope(ip net.IP) int {
	if ip.IsGlobalUnicast() {
		return unix.RT_SCOPE_UNIVERSE
	}
	return unix.RT_SCOPE_LINK
}

func broadcastAddr(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP.To4()
	if ip == nil {
		return nil
	}
	mask := net.IP(ipnet.Mask).To4()
	n := len(ip)
	if n != len(mask) {
		return nil
	}
	out := make(net.IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] | ^mask[i]
	}
	return out
}
