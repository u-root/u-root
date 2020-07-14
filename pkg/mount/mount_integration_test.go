// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package mount

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				// CONFIG_ATA_PIIX is required for this option to work.
				qemu.ArbitraryArgs{"-hda", "testdata/1MB.ext4_vfat"},
				qemu.ArbitraryArgs{"-hdb", "testdata/12Kzeros"},
				qemu.ArbitraryArgs{"-hdc", "testdata/gptdisk"},
				qemu.ArbitraryArgs{"-drive", "file=testdata/gptdisk2,if=none,id=NVME1"},
				qemu.ArbitraryArgs{"-device", "nvme,drive=NVME1,serial=nvme-1"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/mount"}, o)
}
