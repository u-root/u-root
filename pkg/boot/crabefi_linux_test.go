// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"testing"

	"github.com/u-root/uio/uio"
)

func TestCrabEFIImage(t *testing.T) {
	ci := &CrabEFIImage{Payload: uio.NewLazyFile("/missing/crabefi.elf"), BootRank: 42}
	if got := ci.Label(); got != "UEFI (CrabEFI)" {
		t.Errorf("Label() = %q; want UEFI (CrabEFI)", got)
	}
	ci.Name = "Firmware"
	if got := ci.Label(); got != ci.Name {
		t.Errorf("Label() = %q; want %q", got, ci.Name)
	}
	if got := ci.Rank(); got != ci.BootRank {
		t.Errorf("Rank() = %d; want %d", got, ci.BootRank)
	}
	if got, want := ci.String(), "CrabEFI(payload=/missing/crabefi.elf)"; got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}
	ci.Edit(func(string) string {
		t.Error("Edit() called the Linux command-line editor")
		return ""
	})
}

func TestCrabEFINilPayload(t *testing.T) {
	ci := &CrabEFIImage{}
	for _, opts := range [][]LoadOption{nil, {WithDryRun(true), WithLogger(nil)}} {
		if err := ci.Load(opts...); err == nil {
			t.Error("Load(nil payload) = nil; want error, even in dry-run mode")
		}
	}
}
