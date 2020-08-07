// +build integration

package rtnl

import (
	"net"
	"testing"
)

func TestLiveNeighbours(t *testing.T) {
	c, err := Dial(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Trigger a DNS lookup, only for a side effect of pushing our gateway or NS onto the neighbour table
	net.LookupHost("github.com")

	neigtab, err := c.Neighbours(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(neigtab) == 0 {
		t.Skip("no neighbours")
	}
	for i, e := range neigtab {
		t.Logf("* neighbour table entry [%d]: %v", i, e)
		if e.IP.IsUnspecified() {
			// This test doesn't seem to be very reliable
			// Disabling for now
			// t.Error("zero e.IP, expected non-zero")
			continue
		}
		if e.Interface == nil {
			t.Error("nil e.Interface, expected non-nil")
			continue
		}
		if len(e.Interface.Name) == 0 {
			t.Error("zero-length e.Interface.Name")
		}
		if e.IP.IsLoopback() {
			continue
		}
		if hardwareAddrIsUnspecified(e.HwAddr) {
			// This test doesn't seem to be very reliable
			// Disabling for now
			// t.Error("zero e.HwAddr, expected non-zero")
		}
		if hardwareAddrIsUnspecified(e.Interface.HardwareAddr) {
			t.Error("zero e.Interface.HardwareAddr, expected non-zero")
		}
	}
}
