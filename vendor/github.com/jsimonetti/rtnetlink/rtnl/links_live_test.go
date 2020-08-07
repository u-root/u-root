// +build integration

package rtnl

import (
	"bytes"
	"net"
	"testing"
	"strconv"
)

const (
	hardwareAddrLen = 6
)

var (
	hardwareAddrZero = net.HardwareAddr{0, 0, 0, 0, 0, 0}
)

func hardwareAddrEqual(a, b net.HardwareAddr) bool {
	return bytes.Equal(a, b)
}

func hardwareAddrIsUnspecified(hw net.HardwareAddr) bool {
	return len(hw) != hardwareAddrLen || hardwareAddrEqual(hw, hardwareAddrZero)
}

// TestLinks tests the Live function returns sane results
func TestLiveLinks(t *testing.T) {
	c, err := Dial(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	links, err := c.Links()
	if err != nil {
		t.Fatal(err)
	}
	if len(links) == 0 {
		t.Skip("no network interfaces")
	}
	ieth := 0
	ilo := 0
	for i, ifc := range links {
		t.Logf("* entry %d: %#v", i, ifc)
		if ifc.Index == 0 {
			t.Error("zero ifc.Index")
		}
		if ifc.MTU == 0 {
			t.Error("zero ifc.MTU")
		}
		if len(ifc.Name) == 0 {
			t.Error("zero-length ifc.Name")
		}
		if !hardwareAddrIsUnspecified(ifc.HardwareAddr) {
			ieth = ifc.Index
		}
		if ifc.Flags&net.FlagLoopback != 0 {
			ilo = ifc.Index
		}
	}
	if ieth == 0 {
		t.Skip("no interfaces with non-zero link-level address")
	}
	if ilo == 0 {
		t.Skip("no loopback interfaces")
	}

	t.Run("LinkByIndex", func(t *testing.T) {
		for i, ifindex := range []int{ilo, ieth} {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				ifc, err := c.LinkByIndex(ifindex)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("* %#v", ifc)
				if ifc.Index == 0 {
					t.Error("zero ifc.Index")
				}
				if ifc.Index != ifindex {
					t.Errorf("returned wronk interface (%d != %d)", ifc.Index, ifindex)
				}
				if ifc.MTU == 0 {
					t.Error("zero ifc.MTU")
				}
				if len(ifc.Name) == 0 {
					t.Error("zero-length ifc.Name")
				}
				if ifindex == ieth && hardwareAddrIsUnspecified(ifc.HardwareAddr) {
					t.Error("zero ifc.HardwareAddr, expected non-zero")
				}
				if ifindex == ilo && ifc.Flags&net.FlagLoopback == 0  {
					t.Error("no FlagLoopback in ifc.Flags, expected to be set")
				}
			})
		}
	})
}
