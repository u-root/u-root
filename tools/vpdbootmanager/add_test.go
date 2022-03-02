// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/u-root/u-root/pkg/boot/systembooter"
)

func TestParseNetboot(t *testing.T) {
	b, _, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{})
	if err != nil {
		t.Errorf(`parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{}) = _, _, %v, want nil`, err)
	}

	if b.Type != "netboot" || b.Method != "dhcpv4" || b.MAC != "aa:bb:cc:dd:ee:ff" || b.OverrideURL != nil || b.Retries != nil {
		t.Errorf(`b.Type, b.Method, b.MAC, b.OverrideURL, b.Retries = %q, %q, %q, %v, %v, want "netboot, "dpcpv4", "aa:bb:cc:dd:ee:ff", nil, nil`,
			b.Type, b.Method, b.MAC, b.OverrideURL, b.Retries)
	}
}

func TestParseNetbootWithFlags(t *testing.T) {
	flags := []string{
		"-override-url",
		"http://url",
		"-retries",
		"1",
		"-vpd-dir",
		"test",
	}

	b, vpdDir, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", flags)
	if err != nil {
		t.Errorf(`parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", %v) = _, _, %v, want nil`, flags, err)
	}

	if *b.OverrideURL != "http://url" || *b.Retries != 1 || vpdDir != "test" {
		t.Errorf(`*b.OverrideURL, *b.Retries, vpdDir = %q, %d, %q, want "http://url", 1, "test"`, *b.OverrideURL, *b.Retries, vpdDir)
	}
}

func TestParseLocalboot(t *testing.T) {
	b, _, err := parseLocalbootFlags("grub", []string{})
	if err != nil {
		t.Errorf(`parseLocalbootFlags("grub", []string{}) = _, _, %v, want nil`, err)
	}
	if b.Method != "grub" {
		t.Errorf(`b.Method = %q, want "grub"`, b.Method)
	}

	flags := []string{
		"device",
		"path",
	}
	b, _, err = parseLocalbootFlags("path", flags)
	if err != nil {
		t.Errorf(`parseLocalbootFlags("grub", %v) = _, _, %v, want nil`, flags, err)
	}
	if b.Method != "path" || b.DeviceGUID != "device" || b.Kernel != "path" {
		t.Errorf(`b.Method, b.DeviceGUID, b.Kernel = %q, %q, %q`, b.Method, b.DeviceGUID, b.Kernel)
	}
}

func TestParseLocalbootWithFlags(t *testing.T) {
	flags := []string{
		"-kernel-args",
		"kernel-argument-test",
		"-ramfs",
		"ramfs-test",
		"-vpd-dir",
		"test",
	}
	b, vpdDir, err := parseLocalbootFlags("grub", flags)
	if err != nil {
		t.Errorf(`parseLocalbootFlags("grub", %v) = _, _, %v, want nil`, flags, err)
	}

	if b.Method != "grub" || b.KernelArgs != "kernel-argument-test" || b.Initramfs != "ramfs-test" || vpdDir != "test" {
		t.Errorf(`b.Method, b.KernelArgs, b.Initramfs, vpdDir = %q, %q, %q, %q, want "grub", "kernel-argument-test", "ramfs-test", "test"`,
			b.Method, b.KernelArgs, b.Initramfs, vpdDir)
	}

	flags = []string{
		"device",
		"path",
		"-kernel-args",
		"kernel-argument-test",
		"-ramfs",
		"ramfs-test",
		"-vpd-dir",
		"test",
	}
	b, vpdDir, err = parseLocalbootFlags("path", flags)
	if err != nil {
		t.Errorf(`parseLocalbootFlags("path", %v) = _, _, %v, want nil`, flags, err)
	}
	if b.Method != "path" || b.DeviceGUID != "device" || b.Kernel != "path" || b.KernelArgs != "kernel-argument-test" ||
		b.Initramfs != "ramfs-test" || vpdDir != "test" {
		t.Errorf(`b.Method, b.DeviceGUID, b.Kernel, b.KernelArgs, b.Initramfs, vpdDir = %q, %q, %q, %q, %q, %q, want "path", "device", "path", "kernel-argument-test", "ramfs-test", "test"`,
			b.Method, b.DeviceGUID, b.Kernel, b.KernelArgs, b.Initramfs, vpdDir)
	}
}

func TestFailGracefullyMissingArg(t *testing.T) {
	err := add("localboot", []string{})
	if err.Error() != "you need to provide method" {
		t.Error("error message should be: you need to provide method")
	}

	err = add("localboot", []string{"path"})
	if err.Error() != "you need to pass DeviceGUID and Kernel path" {
		t.Error("error message should be: you need to pass DeviceGUID and Kernel path")
	}

	err = add("localboot", []string{"path", "device"})
	if err.Error() != "you need to pass DeviceGUID and Kernel path" {
		t.Error("error message should be: you need to pass DeviceGUID and Kernel path")
	}

	err = add("netboot", []string{})
	if err.Error() != "you need to pass method and MAC address" {
		t.Error("error message should be: you need to pass method and MAC address")
	}

	err = add("netboot", []string{"dhcpv6"})
	if err.Error() != "you need to pass method and MAC address" {
		t.Error("error message should be: you need to pass method and MAC address")
	}
}

func TestFailGracefullyBadMACAddress(t *testing.T) {
	err := add("netboot", []string{"dhcpv6", "test"})
	if err.Error() != "address test: invalid MAC address" {
		t.Errorf(`err.Error() = %q, want "error message should be: address test: invalid MAC address"`, err.Error())
	}
}

func TestFailGracefullyBadNetworkType(t *testing.T) {
	err := add("netboot", []string{"not-valid", "test"})
	if err.Error() != "method needs to be either dhcpv4 or dhcpv6" {
		t.Errorf(`err.Error() = %q, want "error message should be: method needs to be either dhcpv4 or dhcpv6"`, err.Error())
	}
}

func TestFailGracefullyBadLocalbootType(t *testing.T) {
	err := add("localboot", []string{"not-valid"})
	if err.Error() != "method needs to be grub or path" {
		t.Errorf(`err.Error() = %q, want "error message: method needs to be grub or path"`, err.Error())
	}
}

func TestFailGracefullyUnknownEntryType(t *testing.T) {
	err := add("test", []string{})
	if err.Error() != "unknown entry type" {
		t.Errorf(`err.Error() = %q, want "unknown entry type"`, err.Error())
	}
}

func TestAddBootEntry(t *testing.T) {
	vpdDir := t.TempDir()
	if err := os.MkdirAll(path.Join(vpdDir, "rw"), 0o700); err != nil {
		t.Errorf(`os.MkdirAll(path.Join(%q, "rw"), 0o700) = %v, want nil`, vpdDir, err)
	}

	if err := addBootEntry(&systembooter.LocalBooter{
		Method: "grub",
	}, vpdDir); err != nil {
		t.Errorf(`addBootEntry(&systembooter.LocalBooter{"grub"}, %q) = %v, want nil`, vpdDir, err)
	}

	file, err := os.ReadFile(path.Join(vpdDir, "rw", "Boot0001"))
	if err != nil {
		t.Errorf(`os.ReadFile(path.Join(%q, "rw", "Boot0001") = %v, want nil`, vpdDir, err)
	}
	var out systembooter.LocalBooter
	if err = json.Unmarshal([]byte(file), &out); err != nil {
		t.Errorf(`json.Unmarshal([]byte(%v), %v) = %v, want nil`, file, &out, err)
	}

	if out.Method != "grub" {
		t.Errorf(`out.Method = %q, want grub`, out.Method)
	}
}

func TestAddBootEntryMultiple(t *testing.T) {
	vpdDir := t.TempDir()
	err := os.MkdirAll(path.Join(vpdDir, "rw"), 0o700)
	if err != nil {
		t.Errorf(`os.MkdirAll(path.Join(%q, "rw"), 0o700) = %v, want nil`, vpdDir, err)
	}

	for i := 1; i < 5; i++ {
		if err := addBootEntry(&systembooter.LocalBooter{
			Method: "grub",
		}, vpdDir); err != nil {
			t.Errorf(`addBootEntry(&systembooter.LocalBooter{Method: "grub"}, %q) = %v, want nil`, vpdDir, err)
		}
		file, err := os.ReadFile(path.Join(vpdDir, "rw", fmt.Sprintf("Boot%04d", i)))
		if err != nil {
			t.Errorf(`os.ReadFile(path.Join(%q, "rw", fmt.Sprintf("Boot%04d", i))) = %v, want nil`, vpdDir, i, err)
		}
		var out systembooter.LocalBooter
		if err := json.Unmarshal([]byte(file), &out); err != nil {
			t.Errorf(`json.Unmarshal([]byte(%q), %v) = %v, want nil`, file, &out, err)
		}
		if out.Method != "grub" {
			t.Errorf(`out.Method = %q, want grub`, out.Method)
		}
	}
}
