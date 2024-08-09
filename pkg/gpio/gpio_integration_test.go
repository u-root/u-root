// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package gpio

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestIntegration(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/gpio"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			// Make GPIOs nums 10 to 20 available through the
			// mockup driver.
			qemu.WithAppendKernel("gpio-mockup.gpio_mockup_ranges=10,20"),
		),
	)
}
