// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package block

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/pci"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/vmtest"
)

// VM setup:
//
//   /dev/sda is ./testdata/1MB.ext4_vfat
//	/dev/sda1 is ext4
//	/dev/sda2 is vfat
//
//   /dev/sdb is ./testdata/12Kzeros
//	/dev/sdb1 exists, but is not formatted.
//
//   /dev/sdc and /dev/nvme0n1 are ./testdata/gptdisk
//      /dev/sdc1 and /dev/nvme0n1p1 exist (EFI system partition), but is not formatted
//      /dev/sdc2 and /dev/nvme0n1p2 exist (Linux), but is not formatted
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestVM(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				// CONFIG_ATA_PIIX is required for this option to work.
				qemu.ArbitraryArgs{"-hda", "../testdata/1MB.ext4_vfat"},
				qemu.ArbitraryArgs{"-hdb", "../testdata/12Kzeros"},
				qemu.ArbitraryArgs{"-hdc", "../testdata/gptdisk"},
				qemu.ArbitraryArgs{"-drive", "file=../testdata/gptdisk2,if=none,id=NVME1"},
				// use-intel-id uses the vendor=0x8086 and device=0x5845 ids for NVME
				qemu.ArbitraryArgs{"-device", "nvme,drive=NVME1,serial=nvme-1,use-intel-id"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/mount"}, o)
}

func TestGPT(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devpath := fmt.Sprintf("/dev/%sc", getDevicePrefix())

	sdcDisk, err := Device(devpath)
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

func TestBlockDevices(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	want := BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: "nvme0n1p1"},
		&BlockDev{Name: "nvme0n1p2"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	sizes := map[string]uint64{
		"nvme0n1":     50 * 4096,
		"nvme0n1p1":   167 * 512,
		"nvme0n1p2":   166 * 512,
		prefix + "a":  2048 * 512,
		prefix + "a1": 1024 * 512,
		prefix + "a2": 1023 * 512,
		prefix + "b":  24 * 512,
		prefix + "b1": 23 * 512,
		prefix + "c":  50 * 4096,
		prefix + "c1": 167 * 512,
		prefix + "c2": 166 * 512,
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
		"nvme0n1":     nil,
		"nvme0n1p1":   syscall.EINVAL,
		"nvme0n1p2":   syscall.EINVAL,
		prefix + "a":  nil,
		prefix + "a1": syscall.EINVAL,
		prefix + "a2": syscall.EINVAL,
		prefix + "b":  nil,
		prefix + "b1": syscall.EINVAL,
		prefix + "c":  nil,
		prefix + "c1": syscall.EINVAL,
		prefix + "c2": syscall.EINVAL,
	}
	for _, dev := range devs {
		if err := dev.ReadPartitionTable(); err != wantRR[dev.Name] {
			t.Errorf("ReadPartitionTable(%s) = %v, want %v", dev, err, wantRR[dev.Name])
		}
	}
}

func TestFilterHavingPartitions(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	want := BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: "nvme0n1p1"},
		&BlockDev{Name: "nvme0n1p2"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	Debug = log.Printf
	devs = devs.FilterHavingPartitions([]int{1, 2})

	want = BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "c"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func TestFilterPartID(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devname := fmt.Sprintf("%sc2", getDevicePrefix())

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	for _, tt := range []struct {
		guid string
		want BlockDevices
	}{
		{
			guid: "C9865081-266C-4A23-A948-C03DAB506198",
			want: BlockDevices{
				&BlockDev{Name: "nvme0n1p2"},
				&BlockDev{Name: devname},
			},
		},
		{
			guid: "c9865081-266c-4a23-a948-c03dab506198",
			want: BlockDevices{
				&BlockDev{Name: "nvme0n1p2"},
				&BlockDev{Name: devname},
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

	prefix := getDevicePrefix()

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	for _, tt := range []struct {
		guid string
		want BlockDevices
	}{
		{
			// EFI system partition.
			guid: "C12A7328-F81F-11D2-BA4B-00A0C93EC93B",
			want: BlockDevices{
				&BlockDev{Name: "nvme0n1p1"},
				&BlockDev{Name: prefix + "c1"},
			},
		},
		{
			// EFI system partition. mixed case.
			guid: "c12a7328-f81F-11D2-BA4B-00A0C93ec93B",
			want: BlockDevices{
				&BlockDev{Name: "nvme0n1p1"},
				&BlockDev{Name: prefix + "c1"},
			},
		},
		{
			// This is some random Linux GUID.
			guid: "0FC63DAF-8483-4772-8E79-3D69D8477DE4",
			want: BlockDevices{
				&BlockDev{Name: "nvme0n1p2"},
				&BlockDev{Name: prefix + "c2"},
			},
		},
	} {
		parts := devs.FilterPartType(tt.guid)
		if !reflect.DeepEqual(parts, tt.want) {
			t.Errorf("FilterPartType(%s) = %v, want %v", tt.guid, parts, tt.want)
		}
	}
}

func TestFilterBlockPCIString(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	// Check that NVME devices are present.
	want := BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: "nvme0n1p1"},
		&BlockDev{Name: "nvme0n1p2"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	Debug = log.Printf
	devs, err = devs.FilterBlockPCIString("0x8086:0x5845")
	if err != nil {
		t.Fatal(err)
	}

	want = BlockDevices{
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func TestFilterBlockPCI(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	devs, err := GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	// Check that NVME devices are present.
	want := BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: "nvme0n1p1"},
		&BlockDev{Name: "nvme0n1p2"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	p := &pci.PCI{Vendor: 0x8086, Device: 0x5845}
	pl := pci.Devices{p}

	Debug = log.Printf
	devs = devs.FilterBlockPCI(pl)

	want = BlockDevices{
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func getDevicePrefix() string {
	if _, err := os.Stat("/dev/sdc"); err != nil {
		return "vd"
	}
	return "sd"
}
