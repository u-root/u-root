// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qnetwork"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/u-root/pkg/testutil"
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

	script := fmt.Sprintf(`
		dhclient -ipv6=false -v
		ip a
		wget http://192.168.0.2:%d/foobar
		cat ./foobar
		sleep 5
	`, port)

	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			// Build dhclient as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands(
				"github.com/u-root/u-root/cmds/core/dhclient",
			),
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/core/sleep",
				"github.com/u-root/u-root/cmds/core/wget",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qnetwork.HostNetwork("192.168.0.0/24"),
			qnetwork.ServeHTTP(s, ln),
			qemu.VirtioRandom(),
		),
	)
	t.Logf("Command: %v", vm.CmdlineQuoted())
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
	script := `
		dhclient -v -retry 2 -timeout 10
		echo "DHCP timed out"
		sleep 5
	`

	net := qnetwork.NewInterVM()
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			// Build dhclient as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/dhclient"),
			uimage.WithBusyboxCommands("github.com/u-root/u-root/cmds/core/sleep"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			// An empty network so DHCP has something to send packets to.
			net.NewVM(),
			qemu.VirtioRandom(),
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
	serverScript := `
		ip link set eth0 up
		pxeserver -6 -your-ip6=fec0::3 -4=false
	`
	// QEMU doesn't support DHCPv6 for getting IP configuration, so we have
	// to supply our own server.
	//
	// We don't currently have a radvd server we can use, so we also cannot
	// try to download a file using the DHCP configuration.
	net := qnetwork.NewInterVM()
	serverVM := scriptvm.Start(t, "dhcp6_server", serverScript,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/ip",
				"github.com/u-root/u-root/cmds/exp/pxeserver",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
			qemu.VirtioRandom(),
		),
	)

	clientScript := `
		dhclient -ipv4=false -vv
		ip a
	`
	clientVM := scriptvm.Start(t, "dhcp6_client", clientScript,
		scriptvm.WithUimage(
			// Build dhclient as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/dhclient"),
			uimage.WithBusyboxCommands("github.com/u-root/u-root/cmds/core/ip"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			net.NewVM(),
			qemu.VirtioRandom(),
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
