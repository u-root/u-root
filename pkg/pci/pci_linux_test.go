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
