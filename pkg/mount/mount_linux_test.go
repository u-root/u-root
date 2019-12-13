// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestTryMount(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	d, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	sda1 := filepath.Join(d, "sda1")
	if mp, err := TryMount("/dev/sda1", sda1, 0); err != nil {
		t.Errorf("TryMount(/dev/sda1) = %v, want nil", err)
	} else {
		want := &MountPoint{
			Path:   sda1,
			Device: "/dev/sda1",
			FSType: "ext4",
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(/dev/sda1) = %v, want %v", mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}

	sda2 := filepath.Join(d, "sda2")
	if mp, err := TryMount("/dev/sda2", sda2, 0); err != nil {
		t.Errorf("TryMount(/dev/sda2) = %v, want nil", err)
	} else {
		want := &MountPoint{
			Path:   sda2,
			Device: "/dev/sda2",
			FSType: "vfat",
		}
		if !reflect.DeepEqual(mp, want) {
			t.Errorf("TryMount(/dev/sda2) = %v, want %v", mp, want)
		}

		if err := mp.Unmount(0); err != nil {
			t.Errorf("Unmount(%q) = %v, want nil", sda1, err)
		}
	}
}
