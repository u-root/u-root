// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package securelaunch

import (
	"os"
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

// VM setup:
//
//  /dev/sda is ./testdata/mbrdisk
//	  /dev/sda1 is ext4
//	  /dev/sda2 is vfat
//	  /dev/sda3 is fat32
//	  /dev/sda4 is xfs
//
//  /dev/sdb is ./testdata/12Kzeros
//	  /dev/sdb1 exists, but is not formatted.
//
//  /dev/sdc and /dev/nvme0n1 are ./testdata/gptdisk
//    /dev/sdc1 and /dev/nvme0n1p1 exist (EFI system partition), but is not formatted
//    /dev/sdc2 and /dev/nvme0n1p2 exist (Linux), but is not formatted
//
//  /dev/sdd is ./testdata/gptdisk_label
//    /dev/sdd1 is ext4 with no GPT partition label
//    /dev/sdd2 is ext4 with GPT partition label "TEST_LABEL"
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

func TestVM(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				// CONFIG_ATA_PIIX is required for this option to work.
				qemu.ArbitraryArgs{"-hda", "testdata/mbrdisk"},
				qemu.ArbitraryArgs{"-hdb", "testdata/12Kzeros"},
				qemu.ArbitraryArgs{"-hdc", "testdata/gptdisk"},
				qemu.ArbitraryArgs{"-hdd", "testdata/gptdisk_label"},
				qemu.ArbitraryArgs{"-drive", "file=testdata/gptdisk2,if=none,id=NVME1"},
				// use-intel-id uses the vendor=0x8086 and device=0x5845 ids for NVME
				qemu.ArbitraryArgs{"-device", "nvme,drive=NVME1,serial=nvme-1,use-intel-id"},

				// With NVMe devices enabled, kernel crashes when not using q35 machine model.
				qemu.ArbitraryArgs{"-machine", "q35"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/securelaunch"}, o)
}

func TestMountDevice(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Skipping since we are not root")
	}

	if err := GetBlkInfo(); err != nil {
		t.Fatalf("GetBlkInfo() = %v, not nil", err)
	}

	if len(StorageBlkDevices) == 0 {
		t.Fatal("len(StorageBlockDevices) = 0, not > 0")
	}

	mounted := false
	matchExpr := regexp.MustCompile(`[hsv]d[a-z]\d+`)
	for _, device := range StorageBlkDevices {
		if matchExpr.MatchString(device.Name) {
			mountPath, err := MountDevice(device, mount.MS_RDONLY)
			if err != nil || mountPath == "" {
				continue
			}

			if err := UnmountAll(); err != nil {
				continue
			}

			mountPath, err = MountDevice(device, 0)
			if err != nil || mountPath == "" {
				continue
			}

			mounted = true
			break
		}
	}
	if !mounted {
		t.Skip("Skipping since no suitable block device was found to mount")
	}
}
