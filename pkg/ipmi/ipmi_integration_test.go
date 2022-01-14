// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegrationIPMI(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				// This integration test requires kernel built with the following options set:
				// CONFIG_IPMI=y
				// CONFIG_IPMI_DEVICE_INTERFACE=y
				// CONFIG_IPMI_WATCHDOG=y
				// CONFIG_IPMI_SI=y
				qemu.ArbitraryArgs{"-device", "ipmi-bmc-sim,id=bmc0"},
				qemu.ArbitraryArgs{"-device", "isa-ipmi-kcs,bmc=bmc0,irq=5"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/ipmi"}, o)
}
