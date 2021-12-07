// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/boot/systembooter"
)

// Turn these off until testify dies.

func testParseNetboot(t *testing.T) {
	b, _, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{})
	require.NoError(t, err)
	require.Equal(t, "netboot", b.Type)
	require.Equal(t, "dhcpv4", b.Method)
	require.Equal(t, "aa:bb:cc:dd:ee:ff", b.MAC)
	require.Nil(t, b.OverrideURL)
	require.Nil(t, b.Retries)
}

func testParseNetbootWithFlags(t *testing.T) {
	b, vpdDir, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{
		"-override-url",
		"http://url",
		"-retries",
		"1",
		"-vpd-dir",
		"test",
	})
	require.NoError(t, err)
	require.Equal(t, "http://url", *b.OverrideURL)
	require.Equal(t, 1, *b.Retries)
	require.Equal(t, "test", vpdDir)
}

func testParseLocalboot(t *testing.T) {
	b, _, err := parseLocalbootFlags("grub", []string{})
	require.NoError(t, err)
	require.Equal(t, "grub", b.Method)

	b, _, err = parseLocalbootFlags("path", []string{
		"device",
		"path",
	})
	require.NoError(t, err)
	require.Equal(t, "path", b.Method)
	require.Equal(t, "device", b.DeviceGUID)
	require.Equal(t, "path", b.Kernel)
}

func testParseLocalbootWithFlags(t *testing.T) {
	b, vpdDir, err := parseLocalbootFlags("grub", []string{
		"-kernel-args",
		"kernel-argument-test",
		"-ramfs",
		"ramfs-test",
		"-vpd-dir",
		"test",
	})
	require.NoError(t, err)
	require.Equal(t, "grub", b.Method)
	require.Equal(t, "kernel-argument-test", b.KernelArgs)
	require.Equal(t, "ramfs-test", b.Initramfs)
	require.Equal(t, "test", vpdDir)

	b, vpdDir, err = parseLocalbootFlags("path", []string{
		"device",
		"path",
		"-kernel-args",
		"kernel-argument-test",
		"-ramfs",
		"ramfs-test",
		"-vpd-dir",
		"test",
	})
	require.NoError(t, err)
	require.Equal(t, "path", b.Method)
	require.Equal(t, "device", b.DeviceGUID)
	require.Equal(t, "path", b.Kernel)
	require.Equal(t, "kernel-argument-test", b.KernelArgs)
	require.Equal(t, "ramfs-test", b.Initramfs)
	require.Equal(t, "test", vpdDir)
}

func testFailGracefullyMissingArg(t *testing.T) {
	err := add("localboot", []string{})
	require.Equal(t, "you need to provide method", err.Error())

	err = add("localboot", []string{"path"})
	require.Equal(t, "you need to pass DeviceGUID and Kernel path", err.Error())

	err = add("localboot", []string{"path", "device"})
	require.Equal(t, "you need to pass DeviceGUID and Kernel path", err.Error())

	err = add("netboot", []string{})
	require.Equal(t, "you need to pass method and MAC address", err.Error())

	err = add("netboot", []string{"dhcpv6"})
	require.Equal(t, "you need to pass method and MAC address", err.Error())
}

func testFailGracefullyBadMACAddress(t *testing.T) {
	err := add("netboot", []string{"dhcpv6", "test"})
	require.Equal(t, "address test: invalid MAC address", err.Error())
}

func testFailGracefullyBadNetworkType(t *testing.T) {
	err := add("netboot", []string{"not-valid", "test"})
	require.Equal(t, "method needs to be either dhcpv4 or dhcpv6", err.Error())
}

func testFailGracefullyBadLocalbootType(t *testing.T) {
	err := add("localboot", []string{"not-valid"})
	require.Equal(t, "method needs to be grub or path", err.Error())
}

func testFailGracefullyUnknownEntryType(t *testing.T) {
	err := add("test", []string{})
	require.Equal(t, "unknown entry type", err.Error())
}

func testAddBootEntry(t *testing.T) {
	vpdDir := t.TempDir()
	os.MkdirAll(path.Join(vpdDir, "rw"), 0o700)
	defer os.RemoveAll(vpdDir)
	err := addBootEntry(&systembooter.LocalBooter{
		Method: "grub",
	}, vpdDir)
	require.NoError(t, err)
	file, err := os.ReadFile(path.Join(vpdDir, "rw", "Boot0001"))
	require.NoError(t, err)
	var out systembooter.LocalBooter
	err = json.Unmarshal([]byte(file), &out)
	require.NoError(t, err)
	require.Equal(t, "grub", out.Method)
}

func testAddBootEntryMultiple(t *testing.T) {
	vpdDir := t.TempDir()
	os.MkdirAll(path.Join(vpdDir, "rw"), 0o700)
	defer os.RemoveAll(vpdDir)
	for i := 1; i < 5; i++ {
		err := addBootEntry(&systembooter.LocalBooter{
			Method: "grub",
		}, vpdDir)
		require.NoError(t, err)
		file, err := os.ReadFile(path.Join(vpdDir, "rw", fmt.Sprintf("Boot%04d", i)))
		require.NoError(t, err)
		var out systembooter.LocalBooter
		err = json.Unmarshal([]byte(file), &out)
		require.NoError(t, err)
		require.Equal(t, "grub", out.Method)
	}
}
