// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
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
	// Filter by vendor and dev id.
	ven, dev := d[0].Vendor, d[0].Device
	d, err = n.Read(func(p *PCI) bool {
		if p.Vendor == ven && p.Device == dev {
			return false
		}
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	// That should filter just one thing.
	if len(n.(*bus).Devices)-1 != len(d) {
		t.Fatalf("Got %d devices, wanted %d", len(d), len(n.(*bus).Devices)-1)
	}

}
