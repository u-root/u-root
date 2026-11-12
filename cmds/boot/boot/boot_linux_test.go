// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo || tinygo.enable

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/menu"
)

func TestCrabEFIMenuOption(t *testing.T) {
	payload := filepath.Join(t.TempDir(), "crabefi.elf")
	if err := os.WriteFile(payload, []byte("payload"), 0600); err != nil {
		t.Fatal(err)
	}
	linux := &boot.LinuxImage{Name: "Linux", Cmdline: "quiet"}
	images, err := appendCrabEFI([]boot.OSImage{linux}, payload)
	if err != nil {
		t.Fatal(err)
	}
	entries := menu.OSImages(false, images...)
	if len(entries) != 2 || entries[0].Label() != "Linux" || entries[1].Label() != "UEFI (CrabEFI)" {
		t.Fatalf("boot menu = %v; want Linux followed by UEFI (CrabEFI)", entries)
	}
	entries[1].Edit(func(s string) string { t.Error("firmware entry exposed a Linux command line"); return s })
	if linux.Cmdline != "quiet" {
		t.Error("adding firmware entry changed Linux options")
	}
	images, err = appendCrabEFI(images, "")
	if err != nil || len(images) != 2 {
		t.Fatalf("appendCrabEFI(disabled) = %v, %v", images, err)
	}
	images, err = appendCrabEFI(nil, payload)
	if err != nil || len(images) != 1 {
		t.Fatalf("appendCrabEFI(no Linux images) = %v, %v; want one firmware image", images, err)
	}
	ci, ok := images[0].(*boot.CrabEFIImage)
	if !ok {
		t.Fatalf("appendCrabEFI() returned %T; want *boot.CrabEFIImage", images[0])
	}
	data := make([]byte, len("payload"))
	if _, err := ci.Payload.ReadAt(data, 0); err != nil || string(data) != "payload" {
		t.Errorf("CrabEFI payload ReadAt() = %q, %v; want payload, nil", data, err)
	}
}

func TestCrabEFIDisabled(t *testing.T) {
	images, err := appendCrabEFI(nil, "")
	if err != nil || len(images) != 0 {
		t.Errorf("appendCrabEFI(disabled, no images) = %v, %v; want no images, nil", images, err)
	}
}

func TestCrabEFIInvalidPath(t *testing.T) {
	for _, path := range []string{filepath.Join(t.TempDir(), "missing"), t.TempDir()} {
		if _, err := appendCrabEFI(nil, path); err == nil {
			t.Errorf("appendCrabEFI(%q) succeeded; want error", path)
		}
	}
}
