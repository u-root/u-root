// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package block

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/pci"
)

func TestDebug(t *testing.T) {
	Debug("foo")
}

func TestDeviceError(t *testing.T) {
	// Success paths are covered in vm_test.go
	_, err := Device("noexist")
	if err == nil {
		t.Errorf(`Device("noexist") = _,%v, expect an error`, err)
	}
}

func TestBlockDevString(t *testing.T) {
	for _, tt := range []struct {
		name     string
		blockdev *BlockDev
		want     []string
	}{
		{
			name: "complete",
			blockdev: &BlockDev{
				Name:   "devname",
				FSType: "sometype",
				FsUUID: "xxxx",
			},
			want: []string{"devname", "sometype", "xxxx"},
		},
		{
			name: "without FSType",
			blockdev: &BlockDev{
				Name:   "devname",
				FsUUID: "xxxx",
			},
			want: []string{"devname", "xxxx"},
		},
		{
			name: "without FSType and FsUUID",
			blockdev: &BlockDev{
				Name: "devname",
			},
			want: []string{"devname"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.blockdev.String()
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("String() = %s, does not contain %s", got, w)
				}
			}
		})
	}
}

func TestBlockDevDevicePath(t *testing.T) {
	name := "devname"
	b := &BlockDev{Name: name}
	want := filepath.Join("/dev", name)
	got := b.DevicePath()
	if got != want {
		t.Errorf("DevicePath() = %s, want %s", got, want)
	}
}

func TestBlockDevDeviceName(t *testing.T) {
	want := "devname"
	b := &BlockDev{Name: want}
	got := b.DevName()
	if got != want {
		t.Errorf("DevName() = %s, want %s", got, want)
	}
}

func TestBlockDevGPTTableError(t *testing.T) {
	// Success paths are covered in vm_test.go
	tests := []struct {
		name string
		dev  *BlockDev
	}{
		{
			name: "Empty BlockDev",
			dev:  &BlockDev{},
		},
		{
			name: "Not exist",
			dev:  &BlockDev{Name: "noexist"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dev.GPTTable()
			if err == nil {
				t.Error("Expect an error")
			}
		})
	}
}

func TestBlockDevPhysicalBlockSizeError(t *testing.T) {
	// Success paths are covered in vm_test.go
	dev := &BlockDev{Name: "noexist"}
	_, err := dev.PhysicalBlockSize()
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestBlockDevBlockSizeError(t *testing.T) {
	// Success paths are covered in vm_test.go
	dev := &BlockDev{Name: "noexist"}
	_, err := dev.BlockSize()
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestBlockDevKernelBlockSizeError(t *testing.T) {
	// Success paths are covered in vm_test.go
	dev := &BlockDev{Name: "noexist"}
	_, err := dev.KernelBlockSize()
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestBlockDevSizeError(t *testing.T) {
	// Success paths are covered in vm_test.go
	tests := []struct {
		name string
		dev  *BlockDev
	}{
		{
			name: "Empty BlockDev",
			dev:  &BlockDev{},
		},
		{
			name: "Not exist",
			dev:  &BlockDev{Name: "noexist"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dev.Size()
			if err == nil {
				t.Error("Expect an error")
			}
		})
	}
}

func TestBlockDevReadPartitionTableError(t *testing.T) {
	// Success paths are covered in vm_test.go
	dev := &BlockDev{Name: "noexist"}
	err := dev.ReadPartitionTable()
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestBlockDevPCIInfoError(t *testing.T) {
	// Success paths are covered in vm_test.go
	dev := &BlockDev{Name: "noexist"}
	_, err := dev.PCIInfo()
	if err == nil {
		t.Error("Expect an error")
	}
}

func TestBlockDevicesFilterName(t *testing.T) {
	devs := BlockDevices{
		&BlockDev{Name: "devA", FsUUID: "1234-abcd"},
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
		&BlockDev{Name: "devC", FsUUID: "1a2b-3c4d"},
	}

	devs = devs.FilterName("devB")

	want := BlockDevices{
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func TestBlockDevicesFilterNames(t *testing.T) {
	devs := BlockDevices{
		&BlockDev{Name: "devA", FsUUID: "1234-abcd"},
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
		&BlockDev{Name: "devC", FsUUID: "1a2b-3c4d"},
	}

	devs = devs.FilterNames("devA", "devB")

	want := BlockDevices{
		&BlockDev{Name: "devA", FsUUID: "1234-abcd"},
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func TestBlockDevicesFilterFSUUID(t *testing.T) {
	devs := BlockDevices{
		&BlockDev{Name: "devA", FsUUID: "1234-abcd"},
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
		&BlockDev{Name: "devC", FsUUID: "1a2b-3c4d"},
	}

	devs = devs.FilterFSUUID("abcd-1234")

	want := BlockDevices{
		&BlockDev{Name: "devB", FsUUID: "abcd-1234"},
	}
	if !reflect.DeepEqual(devs, want) {
		t.Fatalf("Filtered block devices: \n\t%v \nwant: \n\t%v", devs, want)
	}
}

func TestGetMountpointByDevice(t *testing.T) {
	LinuxMountsPath = "testdata/mounts"

	t.Run("Not exist", func(t *testing.T) {
		if _, err := GetMountpointByDevice("/dev/mapper/sys-oldxxxxxx"); err == nil {
			t.Errorf(`GetMountpointByDevice("/dev/mapper/sys-oldxxxxxx") = _, %v, expect an error`, err)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		mountpoint, err := GetMountpointByDevice("/dev/mapper/sys-old")
		if err != nil {
			t.Errorf(`GetMountpointByDevice("/dev/mapper/sys-old") = _, %v, unexpected error`, err)
		}
		if *mountpoint != "/media/usb" {
			t.Errorf(`*mountpoint = %q, want "/media/usb"`, *mountpoint)
		}
	})
}

func TestParsePCIList(t *testing.T) {
	for _, tt := range []struct {
		name        string
		blockString string
		want        pci.Devices
		err         error
	}{
		{
			name:        "one device",
			blockString: "0x8086:0x1234",
			want:        pci.Devices{&pci.PCI{Vendor: 0x8086, Device: 0x1234}},
			err:         nil,
		},
		{
			name:        "two devices",
			blockString: "0x8086:0x1234,0x1234:0xabcd",
			want: pci.Devices{
				&pci.PCI{Vendor: 0x8086, Device: 0x1234},
				&pci.PCI{Vendor: 0x1234, Device: 0xabcd},
			},
			err: nil,
		},
		{
			name:        "no 0x",
			blockString: "8086:1234,1234:abcd",
			want: pci.Devices{
				&pci.PCI{Vendor: 0x8086, Device: 0x1234},
				&pci.PCI{Vendor: 0x1234, Device: 0xabcd},
			},
			err: nil,
		},
		{
			name:        "capitals",
			blockString: "0x8086:0x1234,0x1234:0xABCD",
			want: pci.Devices{
				&pci.PCI{Vendor: 0x8086, Device: 0x1234},
				&pci.PCI{Vendor: 0x1234, Device: 0xabcd},
			},
			err: nil,
		},
		{
			name:        "not hex vendor",
			blockString: "0xghij:0x1234",
			want:        nil,
			err:         strconv.ErrSyntax,
		},
		{
			name:        "not hex vendor",
			blockString: "0x1234:0xghij",
			want:        nil,
			err:         strconv.ErrSyntax,
		},
		{
			name:        "bad format",
			blockString: "0xghij,0x1234",
			want:        nil,
			err:         ErrListFormat,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			devices, err := parsePCIList(tt.blockString)
			if !errors.Is(err, tt.err) {
				t.Errorf(`errors.Is(%v, %v) = false, want true`, err, tt.err)
			}
			if !reflect.DeepEqual(devices, tt.want) {
				// Need to do this because stringer does not print device and vendor
				s := "got:\n"
				for _, d := range devices {
					s = fmt.Sprintf("%s{Vendor: %v, Device %v}\n", s, d.Vendor, d.Device)
				}
				s = fmt.Sprintf("%swant:\n", s)
				for _, d := range tt.want {
					s = fmt.Sprintf("%s{Vendor: %v, Device %v}\n", s, d.Vendor, d.Device)
				}
				t.Errorf("reflect.DeepEqual(%v, %v) = false, want true", devices, tt.want)
			}
		})
	}
}

func TestComposePartName(t *testing.T) {
	for _, tt := range []struct {
		name    string
		devName string
		partNo  int
		want    string
	}{
		{
			name:    "parent device name ends with a letter #1",
			devName: "sda",
			partNo:  1,
			want:    "sda1",
		},
		{
			name:    "parent device name ends with a letter #2",
			devName: "sdb",
			partNo:  1,
			want:    "sdb1",
		},
		{
			name:    "parent device name ends with a letter #3",
			devName: "sda",
			partNo:  2,
			want:    "sda2",
		},
		{
			name:    "parent device name ends with a letter #4",
			devName: "sdb",
			partNo:  2,
			want:    "sdb2",
		},
		{
			name:    "parent device name ends with a letter, more than 9 partitions",
			devName: "sda",
			partNo:  11,
			want:    "sda11",
		},
		{
			name:    "parent device name ends with a number #1",
			devName: "nvme0n1",
			partNo:  1,
			want:    "nvme0n1p1",
		},
		{
			name:    "parent device name ends with a number #2",
			devName: "nvme0n1",
			partNo:  2,
			want:    "nvme0n1p2",
		},
		{
			name:    "parent device name ends with a number, more than 9 devices",
			devName: "nvme0n10",
			partNo:  1,
			want:    "nvme0n10p1",
		},
		{
			name:    "parent device name ends with a number, more than 9 partitions",
			devName: "nvme0n1",
			partNo:  10,
			want:    "nvme0n1p10",
		},
		{
			name:    "parent device name ends with a number, more than 9 devices ans partitions",
			devName: "nvme0n10",
			partNo:  10,
			want:    "nvme0n10p10",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := ComposePartName(tt.devName, tt.partNo)
			if got != tt.want {
				t.Errorf("ComposePartName(%q, %d) = %q, want %q", tt.devName, tt.partNo, got, tt.want)
			}
		})
	}
}
