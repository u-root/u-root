// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package gpio

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/gpio"}, &vmtest.Options{
		QEMUOpts: qemu.Options{
			// Make GPIOs nums 10 to 20 available through the
			// mockup driver.
			KernelArgs: "gpio-mockup.gpio_mockup_ranges=10,20",
		},
	})
}
