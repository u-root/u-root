// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gpio

import (
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

// GPIO allocations so the tests don't conflict with each other:
//
//   TestExport - 9 and 10
//   TestReadValue - 11
//   TestGetPinID - 12

func TestReadValue(t *testing.T) {
	guest.SkipIfNotInVM(t)

	const gpioNum = 11

	if err := Export(gpioNum); err != nil {
		t.Fatal(err)
	}

	for _, want := range []Value{Low, High} {
		if err := SetOutputValue(gpioNum, want); err != nil {
			t.Fatal(err)
		}

		if val, err := ReadValue(gpioNum); err != nil {
			t.Fatal(err)
		} else if val != want {
			t.Errorf("ReadValue(%d) = %v, want %v", gpioNum, val, want)
		}
	}
}

func TestExport(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if err := Export(10); err != nil {
		t.Errorf("Could not export pin 10: %v", err)
	}

	// Only 10-20 are valid GPIOs in the mock chip.
	if err := Export(9); err == nil {
		t.Errorf("Export(pin 9) should have failed, got nil")
	}
}

func TestGetPinID(t *testing.T) {
	guest.SkipIfNotInVM(t)

	// Base is 10, so we expect 10+2.
	if pin, err := GetPinID("gpio-mockup-A", 2); err != nil {
		t.Errorf("GetPinID(gpio-mockup-A, 2) = %v, want nil", err)
	} else if pin != 12 {
		t.Errorf("GetPinID(gpio-mockup-A, 2) = %v, want 12", pin)
	}

	// There are only 10 GPIOs, so expect this to fail.
	if _, err := GetPinID("gpio-mockup-A", 12); err == nil {
		t.Errorf("GetPinID(gpio-mockup-A, 12) = nil, but wanted error")
	}
}
