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
	"testing"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/scuzz"
	"github.com/u-root/u-root/pkg/storage"
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

	device, err := storage.Device("/dev/sda")
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
