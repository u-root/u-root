// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestNewBusReaderNoGlob(t *testing.T) {
	n, err := NewBusReader()
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewBusReader("*", "*")
	if err != nil {
		t.Fatal(err)
	}
	if len(n.(*bus).Devices) != len(g.(*bus).Devices) {
		t.Fatalf("Got %v, want %v", len(n.(*bus).Devices), len(g.(*bus).Devices))
	}

	for i := range n.(*bus).Devices {
		if n.(*bus).Devices[i] != g.(*bus).Devices[i] {
			t.Errorf("%d: got %q, want %q", i, n.(*bus).Devices[i], g.(*bus).Devices[i])
		}
	}
}

func TestBusReader(t *testing.T) {
	n, err := NewBusReader()
	if err != nil {
		t.Fatal(err)
	}
	if len(n.(*bus).Devices) == 0 {
		t.Fatal("got 0 devices, want at least 1")
	}

	// A single read should be okay.
	d, err := n.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(n.(*bus).Devices) != len(d) {
		t.Fatalf("Got %d devices, wanted %d", len(d), len(n.(*bus).Devices))
	}

	// Multiple reads should be ok
	d, err = n.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(n.(*bus).Devices) != len(d) {
		t.Fatalf("Got %d devices, wanted %d", len(d), len(n.(*bus).Devices))
	}

	// We are going to partition the set into devices which match and
	// devices which don't match ven:dev.
	ven, dev := d[0].Vendor, d[0].Device

	matches, err := n.Read(func(p *PCI) bool {
		return p.Vendor == ven && p.Device == dev
	})
	if err != nil {
		t.Fatal(err)
	}
	notMatches, err := n.Read(func(p *PCI) bool {
		return !(p.Vendor == ven && p.Device == dev)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check that the partitions add up.
	if len(matches)+len(notMatches) != len(n.(*bus).Devices) {
		t.Fatalf("Got %d+%d devices, wanted %d", len(matches), len(notMatches), len(n.(*bus).Devices))
	}
}

func TestBusReadConfig(t *testing.T) {
	var fullread bool
	if os.Getuid() == 0 {
		fullread = true
	}

	r, err := NewBusReader()
	if err != nil {
		log.Fatalf("%v", err)
	}

	d, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}
	o := &bytes.Buffer{}
	// First test is a low verbosity one that should only require 64 bytes.
	if err := d.Print(o, 0, 64); err != nil {
		log.Fatal(err)
	}
	// Second test is a bit more complex. If we are not root, it should
	// get an error. If we are root, it should be ok.
	err = d.Print(o, 0, 256)
	if fullread && err != nil {
		log.Fatalf("Doing a full config read as root: got %v, want nil", err)
	}
	if !fullread && err == nil {
		log.Fatalf("Doing a full config read as ! root: got nil, want %v", os.ErrPermission)
	}
}

func testBaseLimType(t *testing.T) {
	tests := []struct {
		bar    string
		r1, r2 string
	}{
		{bar: "0x0000000000001860 0x0000000000001867 0x0000000000040101", r1: "0x0000000000001860", r2: "0x0000000000001867"},
		{bar: "0x0000000000001867 0x0000000000040101"},
		{bar: "0x000000000001860 0x0?00000000001867 0x0000000000040101"},
		{bar: "0x000000000001860 0x0000000000001867 0x0?00000000040101"},
		{bar: "0x000000?000001860 0x0000000000001867 0x0000000000040101"},
	}
	for _, tt := range tests {
		b, l, a, err := BaseLimType(tt.bar)
		t.Logf("%v %v %v %v", b, l, a, err)
		// if r1 != tt.r1 || r2 != tt.r2 {
		// 	t.Errorf("BAR %s: got \n(%q,%q) want \n(%q,%q)", tt.bar, r1, r2, tt.r1, tt.r2)
		//}
	}
}
