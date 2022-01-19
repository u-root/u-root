// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"

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
				qemu.ArbitraryArgs{"-device", "pci-ipmi-kcs,bmc=bmc0"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/ipmi"}, o)
}

func TestWatchdogRunningQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if _, err := i.WatchdogRunning(); err != nil {
		t.Errorf(`i.WatchdogRunning() = nil, not %q`, err)
	}
}
func TestShutoffWatchdogQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if err := i.ShutoffWatchdog(); err != nil {
		t.Errorf(`i.ShutoffWatchdog() = nil, not %q`, err)
	}
}
func TestGetDeviceIDQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if _, err := i.GetDeviceID(); err != nil {
		t.Errorf(`i.GetDeviceID() = nil, not %q`, err)
	}
}
func TestEnableSELQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if _, err := i.EnableSEL(); err != nil {
		t.Errorf(`i.EnableSEL() = nil, not %q`, err)
	}
}

func TestGetSELInfoQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if _, err := i.GetSELInfo(); err != nil {
		t.Errorf(`i.GetSELInfo() = nil, not %q`, err)
	}
}
func TestGetLanConfigQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	t.Skip("Not supported command")
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if _, err := i.GetLanConfig(1, 1); err != nil {
		t.Errorf(`i.GetLanConfig(1, 1) = nil, not %q`, err)
	}
}
func TestRawCmdQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	data := []byte{0x6, 0x1}
	if _, err := i.RawCmd(data); err != nil {
		t.Errorf(`i.RawCmd(data) = nil, not %q`, err)
	}
}
func TestSetSystemFWVersionQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	t.Skip("Not supported command")
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	if err := i.SetSystemFWVersion("TestTest"); err == nil {
		t.Errorf(`i.SetSystemFWVersion("TestTest") = nil, not %q`, err)
	}
}
func TestLogSystemEventQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0):= i,nil, not nil, %q", err)
	}
	defer i.Close()

	e := &Event{}
	if err := i.LogSystemEvent(e); err != nil {
		t.Errorf(`i.LogSystemEvent(e) = nil, not %q`, err)
	}
}
