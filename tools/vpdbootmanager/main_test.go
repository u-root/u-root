// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/boot/systembooter"
)

func TestInvalidCommand(t *testing.T) {
	err := cli([]string{"unknown"})
	require.Equal(t, "Unrecognized action", err.Error())
}

func TestNoEntryType(t *testing.T) {
	err := cli([]string{"add", "localboot"})
	require.Equal(t, "You need to provide method", err.Error())
}

func TestNoAction(t *testing.T) {
	err := cli([]string{})
	require.Equal(t, "You need to provide action", err.Error())
}

func TestAddNetbootEntryFull(t *testing.T) {
	dir, err := ioutil.TempDir("", "vpdbootmanager")
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(path.Join(dir, "rw"), 0700)
	defer os.RemoveAll(dir)
	err = cli([]string{
		"add",
		"netboot",
		"dhcpv6",
		"aa:bb:cc:dd:ee:ff",
		"-vpd-dir",
		dir,
	})
	require.NoError(t, err)
	file, err := ioutil.ReadFile(path.Join(dir, "rw", "Boot0001"))
	require.NoError(t, err)
	var out systembooter.NetBooter
	err = json.Unmarshal([]byte(file), &out)
	require.NoError(t, err)
	require.Equal(t, "dhcpv6", out.Method)
	require.Equal(t, "aa:bb:cc:dd:ee:ff", out.MAC)
}

func TestAddLocalbootEntryFull(t *testing.T) {
	dir, err := ioutil.TempDir("", "vpdbootmanager")
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(path.Join(dir, "rw"), 0700)
	defer os.RemoveAll(dir)
	err = cli([]string{
		"add",
		"localboot",
		"grub",
		"-vpd-dir",
		dir,
	})
	require.NoError(t, err)
	file, err := ioutil.ReadFile(path.Join(dir, "rw", "Boot0001"))
	require.NoError(t, err)
	var out systembooter.NetBooter
	err = json.Unmarshal([]byte(file), &out)
	require.NoError(t, err)
	require.Equal(t, "grub", out.Method)
}
