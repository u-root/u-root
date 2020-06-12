// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/boot/systembooter"
	"github.com/u-root/u-root/pkg/vpd"
)

func TestParseNetboot(t *testing.T) {
	b, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{})
	require.NoError(t, err)
	require.Equal(t, "netboot", b.Type)
	require.Equal(t, "dhcpv4", b.Method)
	require.Equal(t, "aa:bb:cc:dd:ee:ff", b.MAC)
	require.Nil(t, b.OverrideURL)
	require.Nil(t, b.Retries)
}

func TestParseNetbootWithFlags(t *testing.T) {
	b, err := parseNetbootFlags("dhcpv4", "aa:bb:cc:dd:ee:ff", []string{
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
	require.Equal(t, "test", vpd.VpdDir)
}

func TestParseLocalboot(t *testing.T) {
	b, err := parseLocalbootFlags("grub", []string{})
	require.NoError(t, err)
	require.Equal(t, "grub", b.Method)

	b, err = parseLocalbootFlags("path", []string{
		"device",
		"path",
	})
	require.NoError(t, err)
	require.Equal(t, "path", b.Method)
	require.Equal(t, "device", b.DeviceGUID)
	require.Equal(t, "path", b.Kernel)
}

func TestParseLocalbootWithFlags(t *testing.T) {
	b, err := parseLocalbootFlags("grub", []string{
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
	require.Equal(t, "test", vpd.VpdDir)

	b, err = parseLocalbootFlags("path", []string{
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
	require.Equal(t, "test", vpd.VpdDir)
}

func TestFailGracefullyMissingArg(t *testing.T) {
	err := add("localboot", []string{})
	require.Equal(t, "You need to provide method", err.Error())

	err = add("localboot", []string{"path"})
	require.Equal(t, "You need to pass DeviceGUID and Kernel path", err.Error())

	err = add("localboot", []string{"path", "device"})
	require.Equal(t, "You need to pass DeviceGUID and Kernel path", err.Error())

	err = add("netboot", []string{})
	require.Equal(t, "You need to pass method and MAC address", err.Error())

	err = add("netboot", []string{"dhcpv6"})
	require.Equal(t, "You need to pass method and MAC address", err.Error())
}

func TestFailGracefullyBadMACAddress(t *testing.T) {
	err := add("netboot", []string{"dhcpv6", "test"})
	require.Equal(t, "address test: invalid MAC address", err.Error())
}

func TestFailGracefullyBadNetworkType(t *testing.T) {
	err := add("netboot", []string{"not-valid", "test"})
	require.Equal(t, "Method needs to be either dhcpv4 or dhcpv6", err.Error())
}

func TestFailGracefullyBadLocalbootType(t *testing.T) {
	err := add("localboot", []string{"not-valid"})
	require.Equal(t, "Method needs to be grub or path", err.Error())
}

func TestFailGracefullyUnknownEntryType(t *testing.T) {
	err := add("test", []string{})
	require.Equal(t, "Unknown entry type", err.Error())
}

func TestAddBootEntry(t *testing.T) {
	dir, err := ioutil.TempDir("", "vpdbootmanager")
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(path.Join(dir, "rw"), 0700)
	defer os.RemoveAll(dir)
	vpd.VpdDir = dir
	err = addBootEntry(&systembooter.LocalBooter{
		Method: "grub",
	})
	require.NoError(t, err)
	file, err := ioutil.ReadFile(path.Join(dir, "rw", "Boot0001"))
	require.NoError(t, err)
	var out systembooter.LocalBooter
	err = json.Unmarshal([]byte(file), &out)
	require.NoError(t, err)
	require.Equal(t, "grub", out.Method)
}

func TestAddBootEntryMultiple(t *testing.T) {
	dir, err := ioutil.TempDir("", "vpdbootmanager")
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(path.Join(dir, "rw"), 0700)
	defer os.RemoveAll(dir)
	vpd.VpdDir = dir
	for i := 1; i < 5; i++ {
		err = addBootEntry(&systembooter.LocalBooter{
			Method: "grub",
		})
		require.NoError(t, err)
		file, err := ioutil.ReadFile(path.Join(dir, "rw", fmt.Sprintf("Boot%04d", i)))
		require.NoError(t, err)
		var out systembooter.LocalBooter
		err = json.Unmarshal([]byte(file), &out)
		require.NoError(t, err)
		require.Equal(t, "grub", out.Method)
	}
}
