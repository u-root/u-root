// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/mount"
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
//   /dev/sdc and /dev/nvme0n1 are ./testdata/gptdisk
//      /dev/sdc1 and /dev/nvme0n1p1 exist (EFI system partition), but is not formatted
//      /dev/sdc2 and /dev/nvme0n1p2 exist (Linux), but is not formatted
//
//   ARM tests will load drives as virtio-blk devices (/dev/vd*)

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

// fakeMount lets us count how many times it has been mounted.
type fakeMounter struct {
	name  string
	count int
}

func (m *fakeMounter) DevName() string {
	return m.name
}

func (m *fakeMounter) Mount(path string, flags uintptr) (*mount.MountPoint, error) {
	m.count++
	return &mount.MountPoint{
		Path:   path,
		Device: m.name,
		Flags:  flags,
	}, nil
}

func TestMountPool(t *testing.T) {
	sda1 := &fakeMounter{name: "sda1"}
	sda2 := &fakeMounter{name: "sda2"}
	mp := &mount.Pool{}
	m1, err := mp.Mount(sda1, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = mp.Mount(sda2, 2)
	if err != nil {
		t.Fatal(err)
	}
	m3, err := mp.Mount(sda1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if m1 != m3 {
		t.Fatalf("mp.Mount(sda1) = %v; expected to reuse mount %v", m3, m1)
	}
	if sda1.count != 1 {
		t.Fatalf("expected sda1 mounted 1 times; but mounted %d times", sda1.count)
	}
}

func TestUnmount(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	d, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	mp, err := mount.Mount("", d, "tmpfs", "", 0)
	if err != nil {
		t.Fatalf("Mount(tmpfs) = %v, want nil", err)
	}

	want := &mount.MountPoint{
		Path:   d,
		FSType: "tmpfs",
	}
	if !reflect.DeepEqual(mp, want) {
		t.Fatalf("Mount(tmpfs) = %v, want %v", mp, want)
	}

	if err := mount.Unmount("", false, false); err == nil {
		t.Errorf("Unmount('') got %v, want error", err)
	}
	if err := mount.Unmount(d, true, true); err == nil {
		t.Errorf("Unmount(%s) got %v, want error", d, err)
	}

	// Make the mount busy.
	f, err := os.Create(filepath.Join(d, "busy.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if err := mount.Unmount(d, false, false); !errors.Is(err, syscall.EBUSY) {
		t.Errorf("Unmount(%s) = %v, want EBUSY", d, err)
	}

	// Make unbusy.
	f.Close()

	if err := mount.Unmount(d, true, false); err != nil {
		t.Errorf("Unmount(%s) = %v, want nil", d, err)
	}
}
