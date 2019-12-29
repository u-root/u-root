// +build integration

package rtnl

import (
	"errors"
	"net"
	"syscall"
	"testing"

	"github.com/mdlayher/netlink"
)

var errNoLoopback = errors.New("no loopback interface")

func loopbackInterface(c *Conn) (*net.Interface, error) {
	links, err := c.Links()
	if err != nil {
		return nil, err
	}
	for _, ifc := range links {
		if ifc.Flags&net.FlagLoopback != 0 {
			return ifc, nil
		}
	}
	return nil, errNoLoopback
}

func interfaceHasAddr(c *Conn, ifc *net.Interface, addr *net.IPNet) (bool, error) {
	addrs, err := c.Addrs(ifc, 0)
	if err != nil {
		return false, err
	}
	for _, a := range addrs {
		if ipnetEqual(a, addr) {
			return true, nil
		}
	}
	return false, nil
}

func ipnetEqual(a, b *net.IPNet) bool {
	na, _ := a.Mask.Size()
	nb, _ := b.Mask.Size()
	return na == nb && a.IP.Equal(b.IP)
}

func TestLiveAddrs(t *testing.T) {
	c, err := Dial(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	addrs, err := c.Addrs(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Skip("no network addresses")
	}
	for i, a := range addrs {
		t.Logf("* entry %d: %#v", i, a)
		if a == nil {
			t.Error("nil address result")
			continue
		}
		if a.IP.IsUnspecified() {
			t.Error("address unspecified")
		}
		if a.Mask == nil {
			t.Error("mask unspecified")
		} else {
			ones, bits := a.Mask.Size()
			if bits != 8*net.IPv4len && bits != 8*net.IPv6len {
				t.Error("bad mask length")
			}
			if ones == 0 || ones > bits {
				t.Error("bad prefix length")
			}
		}
	}
}

func TestLiveAddrAddDel(t *testing.T) {
	c, err := Dial(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	lo, err := loopbackInterface(c)
	if err != nil {
		t.Skip(err)
	}
	if lo == nil {
		t.Fatal("no error but nil result")
	}

	testip := MustParseAddr("127.0.0.99/32")

	c.AddrDel(lo, testip) // ok if fails

	if err := c.AddrAdd(lo, testip); err != nil {
		// requires specific privilege - skip the test if can't do
		if v, ok := err.(*netlink.OpError); ok && v.Err == syscall.EPERM {
			t.Skip("AddrAdd: ", err)
		}
		t.Fatal("AddrAdd:", err)
	}
	// if the above suceeded, everything below should do, too
	ok, err := interfaceHasAddr(c, lo, testip)
	if err != nil {
		t.Error(err)
	} else if !ok {
		t.Error("address reported as added but can't confirm")
	}
	if err := c.AddrDel(lo, testip); err != nil {
		t.Error("AddrDel: ", err)
	}
}
