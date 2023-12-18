// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uroot"
	"golang.org/x/exp/slices"
)

// TestDhclientQEMU4 uses QEMU's DHCP server to test dhclient.
func TestDhclientQEMU4(t *testing.T) {
	// TODO: support arm
	if arch := qemu.GuestArch(); !slices.Contains([]qemu.Arch{qemu.ArchAMD64, qemu.ArchArm64}, arch) {
		t.Skipf("test not supported on %s", arch)
	}

	// Create the file to download
	dir := t.TempDir()

	want := "conteeent"
	foobarFile := filepath.Join(dir, "foobar")
	if err := os.WriteFile(foobarFile, []byte(want), 0o644); err != nil {
		t.Fatal(err)
	}

	// Serve HTTP on the host on a random port.
	http.Handle("/", http.FileServer(http.Dir(dir)))
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	s := &http.Server{}
	port := ln.Addr().(*net.TCPAddr).Port

	testCmds := []string{
		"dhclient -ipv6=false -v",
		"ip a",
		// Download a file to make sure dhclient configures kernel networking correctly.
		fmt.Sprintf("wget http://192.168.0.2:%d/foobar", port),
		"cat ./foobar",
		"sleep 5",
		"shutdown -h",
	}
	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/cat",
			"github.com/u-root/u-root/cmds/core/dhclient",
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/sleep",
			"github.com/u-root/u-root/cmds/core/shutdown",
			"github.com/u-root/u-root/cmds/core/wget",
		)}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs(
				"-device", "e1000,netdev=host0",
				"-netdev", "user,id=host0,net=192.168.0.0/24,dhcpstart=192.168.0.10,ipv6=off",
			),
			qemu.WithTask(func(ctx context.Context, n *qemu.Notifications) error {
				if err := s.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
					return err
				}
				return nil
			}),
			qemu.WithTask(func(ctx context.Context, n *qemu.Notifications) error {
				// Wait for VM exit.
				<-n.VMExited
				// Then close HTTP server.
				return s.Close()
			}),
		),
	)
	if _, err := vm.Console.ExpectString("Configured eth0 with IPv4 DHCP Lease"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
	if _, err := vm.Console.ExpectString("inet 192.168.0.10"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
	// "cat ./foobar" should be outputting this.
	if _, err := vm.Console.ExpectString(want); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

/*func TestDhclientTimesOut(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	network := qemu.NewNetwork()
	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestQEMUDHCPTimesOut",
		QEMUOpts: qemu.Options{
			Timeout: 50 * time.Second,
			Devices: []qemu.Device{
				// An empty new network is easier than
				// configuring QEMU not to expose any
				// networking. At the moment.
				network.NewVM(),
			},
		},
		TestCmds: []string{
			"dhclient -v -retry 2 -timeout 10",
			"echo \"DHCP timed out\"",
			"sleep 5",
			"shutdown -h",
		},
	})
	defer ccleanup()

	if err := dhcpClient.Expect("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if err := dhcpClient.Expect("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if err := dhcpClient.Expect("DHCP timed out"); err != nil {
		t.Error(err)
	}
}

func TestDhclient6(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	// QEMU doesn't support DHCPv6 for getting IP configuration, so we have
	// to supply our own server.
	//
	// We don't currently have a radvd server we can use, so we also cannot
	// try to download a file using the DHCP configuration.
	network := qemu.NewNetwork()
	dhcpServer, scleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestDhclient6_Server",
		TestCmds: []string{
			"ip link set eth0 up",
			"pxeserver -6 -your-ip6=fec0::3 -4=false",
		},
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "server"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
	})
	defer scleanup()

	dhcpClient, ccleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "TestDhclient6_Client",
		TestCmds: []string{
			"dhclient -ipv4=false -vv",
			"ip a",
			"shutdown -h",
		},
		QEMUOpts: qemu.Options{
			SerialOutput: vmtest.TestLineWriter(t, "client"),
			Timeout:      30 * time.Second,
			Devices: []qemu.Device{
				network.NewVM(),
			},
		},
	})
	defer ccleanup()

	if err := dhcpServer.Expect("starting dhcpv6 server"); err != nil {
		t.Errorf("%s dhcpv6 server: %v", testutil.NowLog(), err)
	}
	if err := dhcpClient.Expect("Configured eth0 with IPv6 DHCP Lease IP fec0::3"); err != nil {
		t.Errorf("%s configure: %v", testutil.NowLog(), err)
	}
	if err := dhcpClient.Expect("inet6 fec0::3"); err != nil {
		t.Errorf("%s ip: %v", testutil.NowLog(), err)
	}
}*/
