// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"os"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
)

const (
	//Number of configs in the fitimage.itb
	fbcCnt = 2
)

func TestLoadConfig(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	kn, rn, err := i.LoadConfig(false, true)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("kernel name: %s", kn)
	t.Logf("ramdisk name: %s", rn)
	if kn != "kernel@0" {
		t.Fatalf("Expected kernel %s, got %s", "kernel@0", kn)
	}
	if rn != "ramdisk@0" {
		t.Fatalf("Expected ramdisk %s, got %s", "ramdisk@0", rn)
	}

}

func TestLoadConfigMiss(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	i.ConfigOverride = "MagicNonExistentConfig"
	kn, rn, err := i.LoadConfig(false, true)

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

	defer func(old func(i *boot.LinuxImage, verbose bool) error) { loadImage = old }(loadImage)

	loadImage = func(i *boot.LinuxImage, verbose bool) error {
		t.Log("mock load image")
		return nil
	}

	if err = i.Load(true); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfig(t *testing.T) {
	f, err := os.Open("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := ParseConfig(f, true)

	if err != nil {
		t.Fatal(err)
	}

	if len(imgs) != fbcCnt {
		t.Fatalf("Expected 2 images from ParseConfig, got %x", len(imgs))
	}

	cs := [fbcCnt]string{"conf@1", "conf_bz@1"}
	for c, i := range imgs {
		t.Logf("config name: %s", i.ConfigOverride)
		t.Logf("kernel name: %s", i.Kernel)
		t.Logf("ramdisk name: %s", i.InitRAMFS)
		if i.ConfigOverride != cs[c] {
			t.Fatalf("Expected config %s, got %s", cs[c], i.ConfigOverride)
		}
		if i.Kernel != "kernel@0" {
			t.Fatalf("Expected kernel %s, got %s", "kernel@0", i.Kernel)
		}
		if i.InitRAMFS != "ramdisk@0" {
			t.Fatalf("Expected ramdisk %s, got %s", "ramdisk@0", i.InitRAMFS)
		}
	}
}

func TestLabel(t *testing.T) {
	n, kn, rn := "conf_bz@1", "kernel@0", "ramdisk@0"
	img := &Image{name: n, Kernel: kn, InitRAMFS: rn}
	l := img.Label()
	if !strings.Contains(l, n) {
		t.Fatalf("Expected Image label to contain name %s, got %s", n, l)
	}
}
