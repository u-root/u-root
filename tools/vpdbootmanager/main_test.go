// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/u-root/u-root/pkg/boot/systembooter"
)

func TestInvalidCommand(t *testing.T) {
	if err := cli([]string{"unknown"}); err.Error() != "Unrecognized action" {
		t.Errorf(`err.Error() = %q, want "Unrecognized action"`, err.Error())
	}
}

func TestNoEntryType(t *testing.T) {
	if err := cli([]string{"add", "localboot"}); err.Error() != "you need to provide method" {
		t.Errorf(`err.Error() = %q, want "you need to provide method"`, err.Error())
	}
}

func TestNoAction(t *testing.T) {
	if err := cli([]string{}); err.Error() != "you need to provide action" {
		t.Errorf(`err.Error() = %q, want "you need to provide action"`, err.Error())
	}
}

func TestAddNetbootEntryFull(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(path.Join(dir, "rw"), 0o700); err != nil {
		t.Errorf(`os.MkdirAll(path.Join(%q, "rw"), 0o700) = %v, want nil`, dir, err)
	}

	args := []string{
		"add",
		"netboot",
		"dhcpv6",
		"aa:bb:cc:dd:ee:ff",
		"-vpd-dir",
		dir,
	}
	if err := cli(args); err != nil {
		t.Errorf(`cli(%v) = %v, want nil`, args, err)
	}
	file, err := os.ReadFile(path.Join(dir, "rw", "Boot0001"))
	if err != nil {
		t.Errorf(`os.ReadFile(path.Join(%q, "rw", "Boot0001")) = %v, want nil`, file, err)
	}
	var out systembooter.NetBooter
	if err := json.Unmarshal([]byte(file), &out); err != nil {
		t.Errorf(`json.Unmarshal([]byte(%q), %v) = %v, want nil`, file, &out, err)
	}
	if out.Method != "dhcpv6" || out.MAC != "aa:bb:cc:dd:ee:ff" {
		t.Errorf(`out.Method, out.Mac = %q, %q, want "dhcpv6", "aa:bb:cc:dd:ee:ff"`, out.Method, out.MAC)
	}
}

func TestAddLocalbootEntryFull(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(path.Join(dir, "rw"), 0o700); err != nil {
		t.Errorf(`os.MkdirAll(path.Join(%q, "rw"), 0o700) = %v, want nil`, dir, err)
	}

	args := []string{
		"add",
		"localboot",
		"grub",
		"-vpd-dir",
		dir,
	}
	if err := cli(args); err != nil {
		t.Errorf(`cli(%v) = %v, want nil`, args, err)
	}
	file, err := os.ReadFile(path.Join(dir, "rw", "Boot0001"))
	if err != nil {
		t.Errorf(`os.ReadFile(path.Join(%q, "rw", "Boot0001")) = %v, want nil`, file, err)
	}
	var out systembooter.NetBooter
	if err := json.Unmarshal([]byte(file), &out); err != nil {
		t.Errorf(`json.Unmarshal([]byte(%q), %v) = %v, want nil`, file, &out, err)
	}
	if out.Method != "grub" {
		t.Errorf(`out.Method = %q, want "grub"`, out.Method)
	}
}
