// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/scuzz"
	"github.com/u-root/u-root/pkg/testutil"
)

// Assumptions:
//
//   /dev/sda is ./testdata/1MB.ext4_vfat
//	/dev/sda1 is ext4
//	/dev/sda2 is vfat
//
//   /dev/sdb is ./testdata/12Kzeros
//	/dev/sdb1 exists, but is not formatted.

func TestIdentify(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	disk, err := scuzz.NewSGDisk("/dev/sda")
	if err != nil {
		t.Fatal(err)
	}
	defer disk.Close()

	info, err := disk.Identify()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Identify(/dev/sda): %v", info)

	device, err := block.Device("/dev/sda")
	if err != nil {
		t.Fatal(err)
	}
	size, err := device.Size()
	if err != nil {
		t.Fatal(err)
	}

	if info.NumberSectors != size/512 {
		t.Errorf("Identify(/dev/sda).NumberSectors = %d, want %d", info.NumberSectors, size/512)
	}
}

func TestBlockDevices(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	want := block.BlockDevices{
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		// TODO: Why does the vfat one not have a UUID? That must be a bug?
		&block.BlockDev{Name: "sda2"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = %v want %v", devs, want)
	}

	sizes := map[string]uint64{
		"sda":  2048 * 512,
		"sda1": 1024 * 512,
		"sda2": 1023 * 512,
		"sdb":  24 * 512,
		"sdb1": 23 * 512,
	}
	for _, dev := range devs {
		size, err := dev.Size()
		if err != nil {
			t.Errorf("Size(%s) error: %v", dev, err)
		}
		if size != sizes[dev.Name] {
			t.Errorf("Size(%s) = %v, want %v", dev, size, sizes[dev.Name])
		}
	}

	blkSizes := map[string]int{
		"sda":  4096,
		"sda1": 4096,
		"sda2": 512,
		"sdb":  4096,
		"sdb1": 512,
	}
	for _, dev := range devs {
		blkSize, err := dev.BlockSize()
		if err != nil {
			t.Errorf("BlockSize(%s) error: %v", dev, err)
		}
		if blkSize != blkSizes[dev.Name] {
			t.Errorf("BlockSize(%s) = %v, want %v", dev, blkSize, blkSizes[dev.Name])
		}
	}

	wantRR := map[string]error{
		"sda":  nil,
		"sda1": syscall.EINVAL,
		"sda2": syscall.EINVAL,
		"sdb":  nil,
		"sdb1": syscall.EINVAL,
	}
	for _, dev := range devs {
		if err := dev.ReadPartitionTable(); err != wantRR[dev.Name] {
			t.Errorf("ReadPartitionTable(%s) = %v, want %v", dev, err, wantRR[dev.Name])
		}
	}
}

func TestTryMount(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	d, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	sda1 := filepath.Join(d, "sda1")
	if mp, err := mount.TryMount("/dev/sda1", sda1, mount.ReadOnly); err != nil {
		t.Errorf("TryMount(/dev/sda1) = %v, want nil", err)
	} else {
		want := &mount.MountPoint{
			Path:   sda1,
			Device: "/dev/sda1",
			FSType: "ext4",
			Flags:  mount.ReadOnly,
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(/dev/sda1) = %v, want %v", mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}

	sda2 := filepath.Join(d, "sda2")
	if mp, err := mount.TryMount("/dev/sda2", sda2, mount.ReadOnly); err != nil {
		t.Errorf("TryMount(/dev/sda2) = %v, want nil", err)
	} else {
		want := &mount.MountPoint{
			Path:   sda2,
			Device: "/dev/sda2",
			FSType: "vfat",
			Flags:  mount.ReadOnly,
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(/dev/sda2) = %v, want %v", mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}

	sdb1 := filepath.Join(d, "sdb1")
	if _, err := mount.TryMount("/dev/sdb1", sdb1, mount.ReadOnly); !strings.Contains(err.Error(), "no suitable filesystem") {
		t.Errorf("TryMount(/dev/sdb1) = %v, want an error containing 'no suitable filesystem'", err)
	}

	sdc1 := filepath.Join(d, "sdc1")
	if _, err := mount.TryMount("/dev/sdc1", sdc1, mount.ReadOnly); !os.IsNotExist(err) {
		t.Errorf("TryMount(/dev/sdc1) = %v, want an error equivalent to Does Not Exist", err)
	}
}
