// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/mount"
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
	guest.SkipIfNotInVM(t)

	prefix := getDevicePrefix()

	d := t.TempDir()

	sda1 := filepath.Join(d, prefix+"a1")
	deva1 := fmt.Sprintf("/dev/%sa1", prefix)
	if mp, err := mount.TryMount(deva1, sda1, "", mount.ReadOnly, func() error { return os.MkdirAll(sda1, 0o666) }); err != nil {
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
	if mp, err := mount.TryMount(deva2, sda2, "", mount.ReadOnly, func() error { return os.MkdirAll(sda2, 0o666) }); err != nil {
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
	if _, err := mount.TryMount(devb1, sdb1, "", mount.ReadOnly, func() error { return os.MkdirAll(sdb1, 0o666) }); !strings.Contains(err.Error(), "no suitable filesystem") {
		t.Errorf("TryMount(%s) = %v, want an error containing 'no suitable filesystem'", devb1, err)
	}

	sdz1 := filepath.Join(d, prefix+"z1")
	devz1 := fmt.Sprintf("/dev/%sz1", prefix)
	if _, err := mount.TryMount(devz1, sdz1, "", mount.ReadOnly, func() error { return os.MkdirAll(sdz1, 0o666) }); !os.IsNotExist(err) {
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

func (m *fakeMounter) Mount(path string, flags uintptr, opts ...func() error) (*mount.MountPoint, error) {
	m.count++
	return &mount.MountPoint{
		Path:   path,
		Device: filepath.Join("/dev", m.name),
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

func TestIsTmpRamfs(t *testing.T) {
	guest.SkipIfNotInVM(t)

	testRoot := t.TempDir()

	// This test assumes the old mount.Mount behavior of automagically
	// creating the mountpoint directory.
	// This is a good chance to test using the option to create
	// the destination directory ...
	// Is a tmpfs.
	tmpfsMount := filepath.Join(testRoot, "tmpfs")
	mp1, err := mount.Mount("somedevice", tmpfsMount, "tmpfs", "", 0, func() error { return os.MkdirAll(tmpfsMount, 0o666) })
	if err != nil {
		t.Fatalf("mount.Mount(somedevice, %s, tmpfs, 0) returned with error: %v", tmpfsMount, err)
	}
	defer mp1.Unmount(0)

	r, err := mount.IsTmpRamfs(tmpfsMount)
	if err != nil {
		t.Errorf("mount.IsTmpRamfs(%s) returned error: %v, want nil", tmpfsMount, err)
	}

	if !r {
		t.Errorf("mount.IsTmpRamfs(%s) = false, want true", tmpfsMount)
	}

	// Not a tmpfs.
	nottmpfsMount := filepath.Join(testRoot, "nottmpfs")
	mp2, err := mount.Mount("none", nottmpfsMount, "proc", "", 0, func() error { return os.MkdirAll(nottmpfsMount, 0o666) })
	if err != nil {
		t.Fatalf("mount.Mount(somedevice, %s, \"\", \"\", 0)", err)
	}
	defer mp2.Unmount(0)

	r, err = mount.IsTmpRamfs(nottmpfsMount)
	if err != nil {
		t.Errorf("mount.IsTmpRamfs(%s) returned error: %v, want nil", nottmpfsMount, err)
	}
	if r {
		t.Errorf("mount.IsTmpRamfs(%s) = true, want false", nottmpfsMount)
	}
}

func TestOpt(t *testing.T) {
	var (
		j      int
		errVal = fmt.Errorf("got it")
	)

	m, err := mount.Mount("", "", "", "", 0, func() error {
		j++
		return errVal
	})
	if !errors.Is(err, errVal) {
		t.Fatalf(`mount.Mount("", "", "", "", 0, func): %v != %v`, err, errVal)
	}
	if m != nil {
		t.Errorf(`mount.Mount("", "", "", "", 0, func): %v != nil`, m)
	}
	if j != 1 {
		t.Fatalf(`mount.Mount("", "", "", "", 0, func): %d != 1`, j)
	}
}
