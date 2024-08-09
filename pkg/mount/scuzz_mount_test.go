// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package mount_test

import (
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/mount/scuzz"
)

func TestIdentify(t *testing.T) {
	guest.SkipIfNotInVM(t)

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

	device, err := block.Device("/dev/sda")
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
