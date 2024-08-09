// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package brctl

// Sometimes manual testing might be necessary or just more straight forward.
// To setup a local test environment similar to the integration test, run the following commands.
// Since the tests issue raw ioctl calls, they have to be run as root.
//
// ```
// ip link add eth10 type dummy
// ip link add eth10 type dummy
// brctl addbr br0
// brctl addbr br1
// brctl addif br0 eth0
// brctl addif br1 eth1
// ````

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestVM(t *testing.T) {
	// check if arch is amd64
	if runtime.GOARCH != "amd64" {
		t.Skip("skipping integration test")
	}

	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/brctl"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", BRCTL_TEST_IFACE_0)),
			qemu.ArbitraryArgs("-nic", fmt.Sprintf("user,id=%s", BRCTL_TEST_IFACE_1)),
		),
	)
}
