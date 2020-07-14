// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/mount/scuzz"
	"github.com/u-root/u-root/pkg/pci"
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
//
//   /dev/sdc is ./testdata/gptdisk
//      /dev/sdc1 exists (EFI system partition), but is not formatted
//      /dev/sdc2 exists (Linux), but is not formatted

func TestGPT(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	sdcDisk, err := block.Device("/dev/sdc")
	if err != nil {
		t.Fatal(err)
	}

	parts, err := sdcDisk.GPTTable()
	if err != nil {
		t.Fatal(err)
	}
	wantParts := []gpt.Partition{
		{FirstLBA: 34, LastLBA: 200},
		{FirstLBA: 201, LastLBA: 366},
	}
	for i, p := range parts.Partitions {
		if !p.IsEmpty() {
			want := wantParts[i]
			if p.FirstLBA != want.FirstLBA {
				t.Errorf("partition %d: got FirstLBA %d want %d", i, p.FirstLBA, want.FirstLBA)
			}
			if p.LastLBA != want.LastLBA {
				t.Errorf("partition %d: got LastLBA %d want %d", i, p.LastLBA, want.LastLBA)
			}

			t.Logf("partition: %v", p)
		}
	}
}

func TestFilterPartID(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	for _, tt := range []struct {
		guid string
		want block.BlockDevices
	}{
		{
			guid: "C9865081-266C-4A23-A948-C03DAB506198",
			want: block.BlockDevices{
				&block.BlockDev{Name: "sdc2"},
			},
		},
		{
			guid: "c9865081-266c-4a23-a948-c03dab506198",
			want: block.BlockDevices{
				&block.BlockDev{Name: "sdc2"},
			},
		},
		{
			guid: "",
			want: nil,
		},
	} {
		parts := devs.FilterPartID(tt.guid)
		if !reflect.DeepEqual(parts, tt.want) {
			t.Errorf("FilterPartID(%s) = %v, want %v", tt.guid, parts, tt.want)
		}
	}
}

func TestFilterPartType(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	for _, tt := range []struct {
		guid string
		want block.BlockDevices
	}{
		{
			// EFI system partition.
			guid: "C12A7328-F81F-11D2-BA4B-00A0C93EC93B",
			want: block.BlockDevices{
				&block.BlockDev{Name: "sdc1"},
			},
		},
		{
			// EFI system partition. mixed case.
			guid: "c12a7328-f81F-11D2-BA4B-00A0C93ec93B",
			want: block.BlockDevices{
				&block.BlockDev{Name: "sdc1"},
			},
		},
		{
			// This is some random Linux GUID.
			guid: "0FC63DAF-8483-4772-8E79-3D69D8477DE4",
			want: block.BlockDevices{
				&block.BlockDev{Name: "sdc2"},
			},
		},
	} {
		parts := devs.FilterPartType(tt.guid)
		if !reflect.DeepEqual(parts, tt.want) {
			t.Errorf("FilterPartType(%s) = %v, want %v", tt.guid, parts, tt.want)
		}
	}
}

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

func TestFilterBlockPCI(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	// Check that NVME devices are present.
	want := block.BlockDevices{
		&block.BlockDev{Name: "nvme0n1"},
		&block.BlockDev{Name: "nvme0n1p1"},
		&block.BlockDev{Name: "nvme0n1p2"},
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: "sda2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
		&block.BlockDev{Name: "sdc"},
		&block.BlockDev{Name: "sdc1"},
		&block.BlockDev{Name: "sdc2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	p := &pci.PCI{Vendor: "0x8086", Device: "0x5845"}
	pl := pci.Devices{p}

	block.Debug = log.Printf
	devs = devs.FilterBlockPCI(pl)

	want = block.BlockDevices{
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: "sda2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
		&block.BlockDev{Name: "sdc"},
		&block.BlockDev{Name: "sdc1"},
		&block.BlockDev{Name: "sdc2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func TestFilterBlockPCIString(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	// Check that NVME devices are present.
	want := block.BlockDevices{
		&block.BlockDev{Name: "nvme0n1"},
		&block.BlockDev{Name: "nvme0n1p1"},
		&block.BlockDev{Name: "nvme0n1p2"},
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: "sda2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
		&block.BlockDev{Name: "sdc"},
		&block.BlockDev{Name: "sdc1"},
		&block.BlockDev{Name: "sdc2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	block.Debug = log.Printf
	devs, err = devs.FilterBlockPCIString("0x8086:0x5845")
	if err != nil {
		t.Fatal(err)
	}

	want = block.BlockDevices{
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: "sda2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
		&block.BlockDev{Name: "sdc"},
		&block.BlockDev{Name: "sdc1"},
		&block.BlockDev{Name: "sdc2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
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
		&block.BlockDev{Name: "nvme0n1"},
		&block.BlockDev{Name: "nvme0n1p1"},
		&block.BlockDev{Name: "nvme0n1p2"},
		&block.BlockDev{Name: "sda"},
		&block.BlockDev{Name: "sda1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: "sda2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: "sdb"},
		&block.BlockDev{Name: "sdb1"},
		&block.BlockDev{Name: "sdc"},
		&block.BlockDev{Name: "sdc1"},
		&block.BlockDev{Name: "sdc2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	sizes := map[string]uint64{
		"sda":       2048 * 512,
		"sda1":      1024 * 512,
		"sda2":      1023 * 512,
		"sdb":       24 * 512,
		"sdb1":      23 * 512,
		"sdc":       50 * 4096,
		"sdc1":      167 * 512,
		"sdc2":      166 * 512,
		"nvme0n1":   50 * 4096,
		"nvme0n1p1": 167 * 512,
		"nvme0n1p2": 166 * 512,
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

	wantBlkSize := 512
	for _, dev := range devs {
		size, err := dev.BlockSize()
		if err != nil {
			t.Errorf("BlockSize(%s) = %v", dev, err)
		}
		if size != wantBlkSize {
			t.Errorf("BlockSize(%s) = %d, want %d", dev, size, wantBlkSize)
		}

		pSize, err := dev.PhysicalBlockSize()
		if err != nil {
			t.Errorf("PhysicalBlockSize(%s) = %v", dev, err)
		}
		if pSize != wantBlkSize {
			t.Errorf("PhysicalBlockSize(%s) = %d, want %d", dev, pSize, wantBlkSize)
		}
	}

	wantRR := map[string]error{
		"sda":       nil,
		"sda1":      syscall.EINVAL,
		"sda2":      syscall.EINVAL,
		"sdb":       nil,
		"sdb1":      syscall.EINVAL,
		"sdc":       nil,
		"sdc1":      syscall.EINVAL,
		"sdc2":      syscall.EINVAL,
		"nvme0n1":   nil,
		"nvme0n1p1": syscall.EINVAL,
		"nvme0n1p2": syscall.EINVAL,
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
	if mp, err := mount.TryMount("/dev/sda1", sda1, "", mount.ReadOnly); err != nil {
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
	if mp, err := mount.TryMount("/dev/sda2", sda2, "", mount.ReadOnly); err != nil {
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
	if _, err := mount.TryMount("/dev/sdb1", sdb1, "", mount.ReadOnly); !strings.Contains(err.Error(), "no suitable filesystem") {
		t.Errorf("TryMount(/dev/sdb1) = %v, want an error containing 'no suitable filesystem'", err)
	}

	sdz1 := filepath.Join(d, "sdz1")
	if _, err := mount.TryMount("/dev/sdz1", sdz1, "", mount.ReadOnly); !os.IsNotExist(err) {
		t.Errorf("TryMount(/dev/sdz1) = %v, want an error equivalent to Does Not Exist", err)
	}
}
