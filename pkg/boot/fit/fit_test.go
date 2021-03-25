// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"testing"

	"github.com/u-root/u-root/pkg/boot"
)

func TestLoadConfig(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	kn, rn, err := i.LoadConfig(true)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("kernel name: %s", kn)
	t.Logf("ramdisk name: %s", rn)
	if kn != "kernel@0" {
		t.Fatal(err)
	}
	if rn != "ramdisk@0" {
		t.Fatal(err)
	}
}

func TestLoadConfigMiss(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	i.ConfigOverride = "MagicNonExistentConfig"
	kn, rn, err := i.LoadConfig(true)

	if kn != "" {
		t.Fatalf("Kernel %s returned on expected config miss", kn)
	}

	if rn != "" {
		t.Fatalf("Initramfs %s returned on expected config miss", rn)
	}

	if err == nil {
		t.Fatal("Expected error message for miss on FIT config, got nil")
	}
}

func TestLoad(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	i.Kernel, i.InitRAMFS = "kernel@0", "ramdisk@0"

	// Save kexecReboot function and restore at the end
	defer func(old func() error) { kexecReboot = old }(kexecReboot)
	defer func(old func(i *boot.LinuxImage, verbose bool) error) { loadImage = old }(loadImage)
	kexecReboot = func() error {
		t.Log("mock reboot")
		return nil
	}

	loadImage = func(i *boot.LinuxImage, verbose bool) error {
		t.Log("mock load image")
		return nil
	}

	if err = i.Load(true); err != nil {
		t.Fatal(err)
	}
}
