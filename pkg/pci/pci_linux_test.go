// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"log"
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
	r, err := NewBusReader()
	if err != nil {
		log.Fatalf("%v", err)
	}

	d, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}
	d.SetVendorDeviceName()
	if err := d.ReadConfig(); err != nil {
		log.Fatalf("ReadConfig: got %v, want nil", err)
	}
}
