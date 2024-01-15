// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/network"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestDhclientQEMU4 uses QEMU's DHCP server to test dhclient.
func TestDhclientQEMU4(t *testing.T) {
	// Create the file to download
	dir := t.TempDir()
	want := "Hello, world!"
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
			network.IPv4HostNetwork(&net.IPNet{
				IP:   net.IP{192, 168, 0, 0},
				Mask: net.CIDRMask(24, 32),
			}),
			qemu.ServeHTTP(s, ln),
		),
	)
	if _, err := vm.Console.ExpectString("Configured eth0 with IPv4 DHCP Lease"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}
	// "cat ./foobar" should be outputting this.
	if _, err := vm.Console.ExpectString("Hello, world!"); err != nil {
		t.Errorf("%s: %v", testutil.NowLog(), err)
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

func TestDhclientTimesOut(t *testing.T) {
	testCmds := []string{
		"dhclient -v -retry 2 -timeout 10",
		"echo \"DHCP timed out\"",
		"sleep 5",
		"shutdown -h",
	}

	net := network.NewInterVM()
	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/dhclient",
			"github.com/u-root/u-root/cmds/core/sleep",
			"github.com/u-root/u-root/cmds/core/shutdown",
		)}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			// An empty network so DHCP has something to send packets to.
			net.NewVM(),
		),
	)

	if _, err := vm.Console.ExpectString("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if _, err := vm.Console.ExpectString("Could not configure eth0 for IPv"); err != nil {
		t.Error(err)
	}
	if _, err := vm.Console.ExpectString("DHCP timed out"); err != nil {
		t.Error(err)
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

func TestDhclient6(t *testing.T) {
	serverCmds := []string{
		"ip link set eth0 up",
		"pxeserver -6 -your-ip6=fec0::3 -4=false",
	}
	// QEMU doesn't support DHCPv6 for getting IP configuration, so we have
	// to supply our own server.
	//
	// We don't currently have a radvd server we can use, so we also cannot
	// try to download a file using the DHCP configuration.
	net := network.NewInterVM()
	serverVM := vmtest.StartVMAndRunCmds(t, serverCmds,
		vmtest.WithName("TestDhclient6_Server"),
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/exp/pxeserver",
		)}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
		),
	)

	testCmds := []string{
		"dhclient -ipv4=false -vv",
		"ip a",
		"shutdown -h",
	}
	clientVM := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithName("TestDhclient6_Client"),
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/ip",
			"github.com/u-root/u-root/cmds/core/dhclient",
			"github.com/u-root/u-root/cmds/core/shutdown",
		)}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
		),
	)

	if _, err := serverVM.Console.ExpectString("starting dhcpv6 server"); err != nil {
		t.Errorf("%s dhcpv6 server: %v", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("Configured eth0 with IPv6 DHCP Lease IP fec0::3"); err != nil {
		t.Errorf("%s configure: %v", testutil.NowLog(), err)
	}
	if _, err := clientVM.Console.ExpectString("inet6 fec0::3"); err != nil {
		t.Errorf("%s ip: %v", testutil.NowLog(), err)
	}

	if err := clientVM.Wait(); err != nil {
		t.Errorf("Client VM wait: %v", err)
	}
	if err := serverVM.Kill(); err != nil {
		t.Errorf("Server VM could not be killed: %v", err)
	}
	// Would return signal: killed.
	serverVM.Wait()
}
