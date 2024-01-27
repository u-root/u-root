// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestPxeboot runs a server and client to test pxebooting a node.
func TestPxeboot4(t *testing.T) {
	serverScript := `
		ip addr add 192.168.0.1/24 dev eth0
		ip link set eth0 up
		ip route add 0.0.0.0/0 dev eth0
		ls -l /pxeroot
		pxeserver -tftp-dir=/pxeroot
	`
	net := qnetwork.NewInterVM()
	serverVM := vmtest.StartVMAndRunCmds(t, serverScript,
		vmtest.WithName("TestPxeboot_Server"),
		vmtest.WithMergedInitramfs(uroot.Opts{
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/ls",
				"github.com/u-root/u-root/cmds/exp/pxeserver",
			),
			ExtraFiles: []string{
				"./testdata/pxe:pxeroot",
			},
		}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
		),
	)

	clientScript := "pxeboot --no-exec -v"
	clientVM := vmtest.StartVMAndRunCmds(t, clientScript,
		vmtest.WithName("TestPxeboot_Client"),
		vmtest.WithBusyboxCommands("github.com/u-root/u-root/cmds/boot/pxeboot"),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
		),
	)

	if _, err := serverVM.Console.ExpectString("starting file server"); err != nil {
		t.Errorf("%s File server: %v", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("Got DHCPv4 lease on eth0:"); err != nil {
		t.Errorf("%s Lease %v:", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("Boot URI: tftp://192.168.0.1/pxelinux.0"); err != nil {
		t.Errorf("%s Boot: %v", testutil.NowLog(), err)
	}

	// Boot menu should show the label from the pxelinux file.
	if _, err := clientVM.Console.ExpectString("01. some-random-kernel"); err != nil {
		t.Errorf("%s Boot Menu: %v", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("Attempting to boot"); err != nil {
		t.Errorf("%s Boot Menu: %v", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("Kernel: tftp://192.168.0.1/kernel"); err != nil {
		t.Errorf("%s parsed kernel: %v", testutil.NowLog(), err)
	}

	if err := serverVM.Kill(); err != nil {
		t.Error(err)
	}
	serverVM.Wait()

	if err := clientVM.Wait(); err != nil {
		t.Errorf("Client VM Wait: %v", err)
	}
}
