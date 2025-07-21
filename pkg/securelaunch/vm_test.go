// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package securelaunch

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/core/cp"
	"github.com/u-root/u-root/pkg/mount"
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
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	d := t.TempDir()
	mbrdisk := filepath.Join(d, "mbrdisk")
	if err := cp.Default.Copy("testdata/mbrdisk", mbrdisk); err != nil {
		t.Fatalf("copying testdata/mbrdisk to %q:got %v, want nil", mbrdisk, err)
	}
	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/securelaunch"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),

			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", mbrdisk),
			qemu.ArbitraryArgs("-hdb", "testdata/12Kzeros"),
			qemu.ArbitraryArgs("-hdc", "testdata/gptdisk"),
			qemu.ArbitraryArgs("-hdd", "testdata/gptdisk_label"),
			qemu.ArbitraryArgs("-drive", "file=testdata/gptdisk2,if=none,id=NVME1"),
			// use-intel-id uses the vendor=0x8086 and device=0x5845 ids for NVME
			qemu.ArbitraryArgs("-device", "nvme,drive=NVME1,serial=nvme-1,use-intel-id"),

			// With NVMe devices enabled, kernel crashes when not using q35 machine model.
			qemu.ArbitraryArgs("-machine", "q35"),
		),
	)
}

func mountMountDevice(t *testing.T) error {
	t.Helper()
	if err := GetBlkInfo(); err != nil {
		return fmt.Errorf("GetBlkInfo() = %w, not nil", err)
	}

	if len(StorageBlkDevices) == 0 {
		return fmt.Errorf("len(StorageBlockDevices) = 0, not > 0")
	}

	matchExpr := regexp.MustCompile(`[hsv]d[a-z]\d+`)
	t.Logf("Searching for storage in %v with regexp %s", StorageBlkDevices, `[hsv]d[a-z]\d+`)
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

			return nil
		}
	}
	return fmt.Errorf("Skipping since no suitable block device was found to mount:%w", os.ErrNotExist)
}

func TestWriteFile(t *testing.T) {
	guest.SkipIfNotInVM(t)

	Debug = t.Logf
	if err := mountMountDevice(t); err != nil {
		t.Skipf("no mountable device for test:%v", err)
	}

	tempFile := "sda1:" + "/testfile"
	dataStr := "Hello World!"

	if err := WriteFile([]byte(dataStr), tempFile); err != nil {
		t.Fatalf(`WriteFile(dataStr.bytes, tempFile) = %v, not nil`, err)
	}
}

func TestReadFile(t *testing.T) {
	guest.SkipIfNotInVM(t)

	Debug = t.Logf

	if err := mountMountDevice(t); err != nil {
		t.Skipf("no mountable device for test:%v", err)
	}

	tempFile := "sda1:" + "/testfile"
	dataStr := "Hello World!"

	if err := WriteFile([]byte(dataStr), tempFile); err != nil {
		t.Fatalf(`WriteFile(dataStr.bytes, tempFile) = %v, not nil`, err)
	}

	readBytes, err := ReadFile(tempFile)
	if err != nil {
		t.Fatalf(`ReadFile(tempFile) = %v, not nil`, err)
	}

	if !bytes.Equal(readBytes, []byte(dataStr)) {
		t.Fatalf(`ReadFile(tempFile) returned %q, not %q`, readBytes, []byte(dataStr))
	}
}

func TestGetFileBytes(t *testing.T) {
	guest.SkipIfNotInVM(t)

	Debug = t.Logf

	if err := mountMountDevice(t); err != nil {
		t.Skipf("no mountable device for test:%v", err)
	}

	file := "sda1:" + "/file.out"
	fileStr := "Hello, World!"
	fileBytes := []byte(fileStr)

	if err := WriteFile(fileBytes, file); err != nil {
		t.Fatalf(`WriteFile(str, file) = %v, not nil`, err)
	}

	readBytes, err := GetFileBytes(file)
	if err != nil {
		t.Fatalf(`GetFileBytes(file) = %v, not nil`, err)
	}

	if !bytes.Equal(fileBytes, readBytes) {
		t.Fatalf(`GetFileBytes(file) got '%v', want '%v'`, readBytes, fileBytes)
	}
}
