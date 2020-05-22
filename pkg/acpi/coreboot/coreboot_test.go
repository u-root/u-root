// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreboot

import (
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
)

func TestCB(t *testing.T) {
	var rawAPIC = []byte{
		0x00, 0x50, 0x49, 0x43, 0x24, 0x00, 0x00, 0x00, 0x01, 0x8f, 0x50, 0x54, 0x4c, 0x54, 0x44, 0x20,
		0x09, 0x20, 0x41, 0x50, 0x49, 0x43, 0x20, 0x20, 0x00, 0x00, 0x04, 0x06, 0x20, 0x4c, 0x54, 0x50,
		0x00, 0x00, 0x00, 0x00,
	}

	tab, err := acpi.NewRaw(rawAPIC)
	if err != nil {
		t.Fatalf("NewRaw: got %v, want nil", err)
	}
	// This tab has a bad signature
	if _, err := NewCorebooter(tab[0]); err == nil {
		t.Errorf("got nil, want err")
	}
	// Test good signature, bad table.
	rawAPIC[0] = 'A'
	if tab, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Fatalf("NewRaw: got %v, want nil", err)
	}
	if _, err := NewCorebooter(tab[0]); err == nil {
		t.Errorf("got nil, want err")
	}
	// Create a good table.
	rawAPIC = append(rawAPIC, make([]byte, minMADT)...)
	rawAPIC[4] += 8
	if tab, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Fatalf("NewRaw: got %v, want nil", err)
	}
	if _, err := NewCorebooter(tab[0]); err != nil {
		t.Fatalf("NewRaw: got %v, want nil", err)
	}
}
