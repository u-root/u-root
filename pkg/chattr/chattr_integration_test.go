// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package chattr

// Tests issue raw ioctl calls, they have to be run as root.

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestVM(t *testing.T) {
	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/chattr"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
		),
	)
}
