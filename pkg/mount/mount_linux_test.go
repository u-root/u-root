// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount_test

import (
	"fmt"
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
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestGPT(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	devpath := fmt.Sprintf("/dev/%sc", getDevicePrefix())

	sdcDisk, err := block.Device(devpath)
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

	devname := fmt.Sprintf("%sc2", getDevicePrefix())

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
				&block.BlockDev{Name: devname},
			},
		},
		{
			guid: "c9865081-266c-4a23-a948-c03dab506198",
			want: block.BlockDevices{
				&block.BlockDev{Name: devname},
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
				&block.BlockDev{Name: prefix + "c1"},
			},
		},
		{
			// EFI system partition. mixed case.
			guid: "c12a7328-f81F-11D2-BA4B-00A0C93ec93B",
			want: block.BlockDevices{
				&block.BlockDev{Name: prefix + "c1"},
			},
		},
		{
			// This is some random Linux GUID.
			guid: "0FC63DAF-8483-4772-8E79-3D69D8477DE4",
			want: block.BlockDevices{
				&block.BlockDev{Name: prefix + "c2"},
			},
		},
	} {
		parts := devs.FilterPartType(tt.guid)
		if !reflect.DeepEqual(parts, tt.want) {
			t.Errorf("FilterPartType(%s) = %v, want %v", tt.guid, parts, tt.want)
		}
	}
}

func TestFilterBlockPCI(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

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
		&block.BlockDev{Name: prefix + "a"},
		&block.BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: prefix + "b"},
		&block.BlockDev{Name: prefix + "b1"},
		&block.BlockDev{Name: prefix + "c"},
		&block.BlockDev{Name: prefix + "c1"},
		&block.BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}

	p := &pci.PCI{Vendor: "0x8086", Device: "0x5845"}
	pl := pci.Devices{p}

	block.Debug = log.Printf
	devs = devs.FilterBlockPCI(pl)

	want = block.BlockDevices{
		&block.BlockDev{Name: prefix + "a"},
		&block.BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: prefix + "b"},
		&block.BlockDev{Name: prefix + "b1"},
		&block.BlockDev{Name: prefix + "c"},
		&block.BlockDev{Name: prefix + "c1"},
		&block.BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func TestFilterBlockPCIString(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

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
		&block.BlockDev{Name: prefix + "a"},
		&block.BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: prefix + "b"},
		&block.BlockDev{Name: prefix + "b1"},
		&block.BlockDev{Name: prefix + "c"},
		&block.BlockDev{Name: prefix + "c1"},
		&block.BlockDev{Name: prefix + "c2"},
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
		&block.BlockDev{Name: prefix + "a"},
		&block.BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: prefix + "b"},
		&block.BlockDev{Name: prefix + "b1"},
		&block.BlockDev{Name: prefix + "c"},
		&block.BlockDev{Name: prefix + "c1"},
		&block.BlockDev{Name: prefix + "c2"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("BlockDevices() = \n\t%v want\n\t%v", devs, want)
	}
}

func TestBlockDevices(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	devs, err := block.GetBlockDevices()
	if err != nil {
		t.Fatal(err)
	}
	devs = devs.FilterZeroSize()

	want := block.BlockDevices{
		&block.BlockDev{Name: "nvme0n1"},
		&block.BlockDev{Name: "nvme0n1p1"},
		&block.BlockDev{Name: "nvme0n1p2"},
		&block.BlockDev{Name: prefix + "a"},
		&block.BlockDev{Name: prefix + "a1", FsUUID: "2183ead8-a510-4b3d-9777-19c7090f66d9"},
		&block.BlockDev{Name: prefix + "a2", FsUUID: "ace5-5144"},
		&block.BlockDev{Name: prefix + "b"},
		&block.BlockDev{Name: prefix + "b1"},
		&block.BlockDev{Name: prefix + "c"},
		&block.BlockDev{Name: prefix + "c1"},
		&block.BlockDev{Name: prefix + "c2"},
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

func TestTryMount(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	prefix := getDevicePrefix()

	d, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	sda1 := filepath.Join(d, prefix+"a1")
	deva1 := fmt.Sprintf("/dev/%sa1", prefix)
	if mp, err := mount.TryMount(deva1, sda1, "", mount.ReadOnly); err != nil {
		t.Errorf("TryMount(%s) = %v, want nil", deva1, err)
	} else {
		want := &mount.MountPoint{
			Path:   sda1,
			Device: deva1,
			FSType: "ext4",
			Flags:  mount.ReadOnly,
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(%s) = %v, want %v", deva1, mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}

	sda2 := filepath.Join(d, prefix+"a2")
	deva2 := fmt.Sprintf("/dev/%sa2", prefix)
	if mp, err := mount.TryMount(deva2, sda2, "", mount.ReadOnly); err != nil {
		t.Errorf("TryMount(%s) = %v, want nil", deva2, err)
	} else {
		want := &mount.MountPoint{
			Path:   sda2,
			Device: deva2,
			FSType: "vfat",
			Flags:  mount.ReadOnly,
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(%s) = %v, want %v", deva2, mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}

	sdb1 := filepath.Join(d, prefix+"b1")
	devb1 := fmt.Sprintf("/dev/%sb1", prefix)
	if _, err := mount.TryMount(devb1, sdb1, "", mount.ReadOnly); !strings.Contains(err.Error(), "no suitable filesystem") {
		t.Errorf("TryMount(%s) = %v, want an error containing 'no suitable filesystem'", devb1, err)
	}

	sdz1 := filepath.Join(d, prefix+"z1")
	devz1 := fmt.Sprintf("/dev/%sz1", prefix)
	if _, err := mount.TryMount(devz1, sdz1, "", mount.ReadOnly); !os.IsNotExist(err) {
		t.Errorf("TryMount(%s) = %v, want an error equivalent to Does Not Exist", devz1, err)
	}
}

func getDevicePrefix() string {
	if _, err := os.Stat("/dev/sdc"); err != nil {
		return "vd"
	}
	return "sd"
}
