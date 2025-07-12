// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package mount

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/u-root/pkg/core/cp"
)

func TestIntegration(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	// qemu likes to lock files.
	// In practice we've seen issues with multiple instantiations of
	// qemu getting into lock wars. To avoid this, copy data files to
	// a temp directory.
	tmp := testtmp.TempDir(t)

	// We do not use CopyTree as it (1) recreates the full path in the tmp directory,
	// and (2) we want to only copy what we want to copy.
	for _, f := range []string{"1MB.ext4_vfat", "12Kzeros", "gptdisk", "gptdisk2"} {
		s := filepath.Join("./testdata", f)
		d := filepath.Join(tmp, f)
		if err := cp.Copy(s, d); err != nil {
			t.Fatalf("Copying %q to %q: got %v, want nil", s, d, err)
		}
	}

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/mount"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			// CONFIG_ATA_PIIX is required for this option to work.
			qemu.ArbitraryArgs("-hda", filepath.Join(tmp, "1MB.ext4_vfat")),
			qemu.ArbitraryArgs("-hdb", filepath.Join(tmp, "12Kzeros")),
			qemu.ArbitraryArgs("-hdc", filepath.Join(tmp, "gptdisk")),
			qemu.ArbitraryArgs("-drive", "file="+filepath.Join(tmp, "gptdisk2")+",if=none,id=NVME1"),
			// use-intel-id uses the vendor=0x8086 and device=0x5845 ids for NVME
			qemu.ArbitraryArgs("-device", "nvme,drive=NVME1,serial=nvme-1,use-intel-id"),
		),
	)
}
