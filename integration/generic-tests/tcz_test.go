// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"regexp"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

func TestTczclient(t *testing.T) {
	// TODO: support arm
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	t.Skip("This test is flaky, and must be fixed")

	serverScript := `
		ip addr add 192.168.0.1/24 dev eth0
		ip link set eth0 up
		ip route add 255.255.255.255/32 dev eth0
		ip l
		ip a
		srvfiles -h 192.168.0.1 -d /
		echo The Server Completes
		shutdown -h
	`
	net := qnetwork.NewInterVM()
	serverVM := scriptvm.Start(t, "tcz_server", serverScript,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/ls",
				"github.com/u-root/u-root/cmds/core/shutdown",
				"github.com/u-root/u-root/cmds/exp/srvfiles",
				"github.com/u-root/u-root/cmds/exp/pxeserver",
			),
			uimage.WithFiles(
				"./testdata/tczserver:tcz",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
			qemu.VirtioRandom(),
		),
	)

	clientScript := `
		ip addr add 192.168.0.2/24 dev eth0
		ip link set eth0 up
		ip a
		tcz -d -h 192.168.0.1 -p 8080 libXcomposite libXdamage libXinerama libxshmfence
		tcz -d -h 192.168.0.1 -p 8080 libXdmcp
		echo HI THERE
		ls /TinyCorePackages/tcloop
	`

	var b wc
	clientVM := scriptvm.Start(t, "tcz_client", clientScript,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/ls",
				"github.com/u-root/u-root/cmds/core/shutdown",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/exp/tcz",
			),
			uimage.WithFiles(
				"./testdata/tczclient:tcz",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
			qemu.WithSerialOutput(&b),
			qemu.VirtioRandom(),
		),
	)

	// The directory list is the last thing we get. At that point,
	// b will have the output we care about and the VM will have shut
	// down. We can do the rest of the RE matching on b.String()
	// This is a bit of a hack but it frees us from worrying
	// about the order in which things appear.
	tczs := []string{"libXcomposite", "libXdamage", "libXinerama", "libxshmfence"}
	for _, s := range tczs {
		if _, err := clientVM.Console.ExpectString(s); err != nil {
			t.Errorf("got %v, want nil", err)
		}
		t.Logf("Matched %s", s)
	}

	if false {
		for _, s := range tczs {
			re, err := regexp.Compile(".*loop.*" + s)
			if err != nil {
				t.Errorf("Check loop device re %s: got %v, want nil", s, err)
				continue
			}
			if ok := re.MatchString(b.String()); !ok {
				t.Errorf("Check loop device %s: got no match, want match", s)
				continue
			}
		}
	}

	if err := clientVM.Wait(); err != nil {
		t.Errorf("Client Wait: %v", err)
	}

	serverVM.Kill()
	serverVM.Wait()
}
