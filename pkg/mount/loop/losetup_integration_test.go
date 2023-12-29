// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package loop

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestIntegration(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	vmtest.RunGoTestsInVM(t, []string{"github.com/u-root/u-root/pkg/mount/loop"},
		vmtest.WithVMOpt(vmtest.WithQEMUFn(qemu.WithVMTimeout(time.Minute))),
	)
}
