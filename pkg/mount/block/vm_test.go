// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package block

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/pci"
)

// VM setup:
//
//  /dev/sda is ./testdata/mbrdisk
//	  /dev/sda1 is ext4
//	  /dev/sda2 is vfat
//	  /dev/sda3 is fat32
//	  /dev/sda4 is xfs
//
//  /dev/sdb is ../testdata/12Kzeros
//	  /dev/sdb1 exists, but is not formatted.
//
//  /dev/sdc and /dev/nvme0n1 are ../testdata/gptdisk
//    /dev/sdc1 and /dev/nvme0n1p1 exist (EFI system partition), but is not formatted
//    /dev/sdc2 and /dev/nvme0n1p2 exist (Linux), but is not formatted
//
//  /dev/sdd is ./testdata/gptdisk_label
//    /dev/sdd1 is ext4 with no GPT partition label
//    /dev/sdd2 is ext4 with GPT partition label "TEST_LABEL"
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestVM(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)
	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/mount/block"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", "testdata/mbrdisk"),
			qemu.ArbitraryArgs("-hdb", "../testdata/12Kzeros"),
			qemu.ArbitraryArgs("-hdc", "../testdata/gptdisk"),
			qemu.ArbitraryArgs("-hdd", "testdata/gptdisk_label"),
			qemu.ArbitraryArgs("-drive", "file=../testdata/gptdisk2,if=none,id=NVME1"),
			// use-intel-id uses the vendor=0x8086 and device=0x5845 ids for NVME
			qemu.ArbitraryArgs("-device", "nvme,drive=NVME1,serial=nvme-1,use-intel-id"),

			// With NVMe devices enabled, kernel crashes when not using q35 machine model.
			qemu.ArbitraryArgs("-machine", "q35"),
		),
	)
}

func TestBlockDevMount(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)
	devs = devs.FilterName("sdd1")
	want := BlockDevices{
		&BlockDev{Name: prefix + "d1", FsUUID: "02175989-d49f-4e8e-836e-99300af66fc1"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Test block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}

	mountPath := t.TempDir()

	// Only testing for the right calls to the mount pkg here.
	// Mounting itself is out of scope here and covered in pkg mount.

	dev := devs[0] // FSType unset
	mp, err := dev.Mount(mountPath, mount.ReadOnly)
	if err != nil {
		t.Errorf("%s.Mount() = _,%v \nunexpected error", dev.Name, err)
	}
	if err := mp.Unmount(0); err != nil {
		t.Fatal(err)
	}

	dev.FSType = "ext4" // FSType set
	mp, err = dev.Mount(mountPath, mount.ReadOnly)
	if err != nil {
		t.Errorf("%s.Mount() = _,%v \nunexpected error", dev.Name, err)
	}
	if err := mp.Unmount(0); err != nil {
		t.Fatal(err)
	}
}

func TestBlockDevGPTTable(t *testing.T) {
	guest.SkipIfNotInVM(t)

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

func TestBlockDevSizes(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	sizes := map[string]uint64{
		"nvme0n1":     50 * 4096,
		"nvme0n1p1":   167 * 512,
		"nvme0n1p2":   166 * 512,
		prefix + "a":  36864 * 512,
		prefix + "a1": 1024 * 512,
		prefix + "a2": 1023 * 512,
		prefix + "a3": 1024 * 512,
		prefix + "a4": 32768 * 512,
		prefix + "b":  24 * 512,
		prefix + "b1": 23 * 512,
		prefix + "c":  50 * 4096,
		prefix + "c1": 167 * 512,
		prefix + "c2": 166 * 512,
		prefix + "d":  512 * 4096,
		prefix + "d1": 1025 * 512,
		prefix + "d2": 990 * 512,
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

		kSize, err := dev.KernelBlockSize()
		if err != nil {
			t.Errorf("KernelBlockSize(%s) = %v", dev, err)
		}
		if pSize != wantBlkSize {
			t.Errorf("KernelBlockSize(%s) = %d, want %d", dev, kSize, wantBlkSize)
		}
	}
}

func TestBlockDevReadPartitionTable(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	wantRR := map[string]error{
		"nvme0n1":     nil,
		"nvme0n1p1":   syscall.EINVAL,
		"nvme0n1p2":   syscall.EINVAL,
		prefix + "a":  nil,
		prefix + "a1": syscall.EINVAL,
		prefix + "a2": syscall.EINVAL,
		prefix + "a3": syscall.EINVAL,
		prefix + "a4": syscall.EINVAL,
		prefix + "b":  nil,
		prefix + "b1": syscall.EINVAL,
		prefix + "c":  nil,
		prefix + "c1": syscall.EINVAL,
		prefix + "c2": syscall.EINVAL,
		prefix + "d":  nil,
		prefix + "d1": syscall.EINVAL,
		prefix + "d2": syscall.EINVAL,
	}
	for _, dev := range devs {
		if err := dev.ReadPartitionTable(); err != wantRR[dev.Name] {
			t.Errorf("ReadPartitionTable(%s) : \n\t%v \nwant: \n\t%v", dev, err, wantRR[dev.Name])
		}
	}
}

func TestBlockDevicesFilterHavingPartitions(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	devs = devs.FilterHavingPartitions([]int{1, 2})

	want := BlockDevices{
		&BlockDev{Name: "nvme0n1"},
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "d"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func TestBlockDevicesFilterPartID(t *testing.T) {
	guest.SkipIfNotInVM(t)

	devname := fmt.Sprintf("%sc2", getDevicePrefix())
	devs := testDevs(t)

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
			t.Errorf("FilterPartID(%s) : \n\t%v \nwant: \n\t%v", tt.guid, parts, tt.want)
		}
	}
}

func TestBlockDevicesFilterPartType(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

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
				&BlockDev{Name: prefix + "d1", FsUUID: "02175989-d49f-4e8e-836e-99300af66fc1"},
				&BlockDev{Name: prefix + "d2", FsUUID: "f3323a7f-a90a-4342-9508-d042afed287d"},
			},
		},
	} {
		parts := devs.FilterPartType(tt.guid)
		if !reflect.DeepEqual(parts, tt.want) {
			t.Errorf("FilterPartType(%s) : \n\t%v \nwant: \n\t%v", tt.guid, parts, tt.want)
		}
	}
}

func TestBlockDevicesFilterPartLabel(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	label := "TEST_LABEL"
	want := BlockDevices{
		&BlockDev{Name: prefix + "d2", FsUUID: "f3323a7f-a90a-4342-9508-d042afed287d"},
	}

	parts := devs.FilterPartLabel(label)
	if !reflect.DeepEqual(parts, want) {
		t.Errorf("FilterPartLabel(%s) : \n\t%v \nwant: \n\t%v", label, parts, want)
	}
}

func TestBlockDevicesFilterBlockPCIString(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	Debug = log.Printf
	devs, err := devs.FilterBlockPCIString("0x8086:0x5845")
	if err != nil {
		t.Fatal(err)
	}

	want := BlockDevices{
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "a3", FsUUID: "a896-d7b8"},
		&BlockDev{Name: prefix + "a4", FsUUID: "dca5f234-726b-47e2-b16e-07d3dbde7d8c"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
		&BlockDev{Name: prefix + "d"},
		&BlockDev{Name: prefix + "d1", FsUUID: "02175989-d49f-4e8e-836e-99300af66fc1"},
		&BlockDev{Name: prefix + "d2", FsUUID: "f3323a7f-a90a-4342-9508-d042afed287d"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func TestBlockDevicesFilterBlockPCI(t *testing.T) {
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()
	devs := testDevs(t)

	p := &pci.PCI{Vendor: 0x8086, Device: 0x5845}
	pl := pci.Devices{p}

	Debug = log.Printf
	devs = devs.FilterBlockPCI(pl)

	want := BlockDevices{
		&BlockDev{Name: prefix + "a"},
		&BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&BlockDev{Name: prefix + "a3", FsUUID: "a896-d7b8"},
		&BlockDev{Name: prefix + "a4", FsUUID: "dca5f234-726b-47e2-b16e-07d3dbde7d8c"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
		&BlockDev{Name: prefix + "d"},
		&BlockDev{Name: prefix + "d1", FsUUID: "02175989-d49f-4e8e-836e-99300af66fc1"},
		&BlockDev{Name: prefix + "d2", FsUUID: "f3323a7f-a90a-4342-9508-d042afed287d"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func getDevicePrefix() string {
	if _, err := os.Stat("/dev/sdc"); err != nil {
		return "vd"
	}
	return "sd"
}

func testDevs(t *testing.T) BlockDevices {
	t.Helper()

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
		&BlockDev{Name: prefix + "a3", FsUUID: "a896-d7b8"},
		&BlockDev{Name: prefix + "a4", FsUUID: "dca5f234-726b-47e2-b16e-07d3dbde7d8c"},
		&BlockDev{Name: prefix + "b"},
		&BlockDev{Name: prefix + "b1"},
		&BlockDev{Name: prefix + "c"},
		&BlockDev{Name: prefix + "c1"},
		&BlockDev{Name: prefix + "c2"},
		&BlockDev{Name: prefix + "d"},
		&BlockDev{Name: prefix + "d1", FsUUID: "02175989-d49f-4e8e-836e-99300af66fc1"},
		&BlockDev{Name: prefix + "d2", FsUUID: "f3323a7f-a90a-4342-9508-d042afed287d"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Test block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}

	return devs
}
