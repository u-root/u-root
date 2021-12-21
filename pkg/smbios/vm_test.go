// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"testing"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegration(t *testing.T) {
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/smbios"}, &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{
					"-smbios",
					"type=2,manufacturer=u-root",
				},
			},
		},
	})
}

func TestFromSysfs(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil || info == nil {
		t.Errorf("FromSysfs() = %q, '%v', want nil", info, err)
	}
}

func TestGetBIOSInfo(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	smbiosinfo, err := info.GetBIOSInfo()
	if err != nil || smbiosinfo == nil {
		t.Errorf("GetBiosInfo() = %q, '%v'", smbiosinfo, err)
	}
}

func TestGetSystemInfo(t *testing.T) {
	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	systeminfo, err := info.GetSystemInfo()
	if err != nil || systeminfo == nil {
		t.Errorf("GetSystemInfo() = %q, '%v'", systeminfo, err)
	}
}

func TestGetChassisInfo(t *testing.T) {

	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	chassisinfo, err := info.GetChassisInfo()
	if err != nil || chassisinfo == nil {
		t.Errorf("GetChassisInfo() = %q, '%v'", chassisinfo, err)
	}
}

func TestGetProcessorInfo(t *testing.T) {

	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	processorinfo, err := info.GetProcessorInfo()
	if err != nil || processorinfo == nil {
		t.Errorf("GetProcessorInfo() = %q, '%v'", processorinfo, err)
	}
}

func TestGetSystemSlots(t *testing.T) {

	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	systemslots, err := info.GetSystemSlots()
	if err != nil || systemslots == nil {
		t.Errorf("GetSystemSlots() = %q, '%v'", systemslots, err)
	}
}

func TestGetMemoryDevices(t *testing.T) {

	testutil.SkipIfNotRoot(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	memorydevices, err := info.GetMemoryDevices()
	if err != nil || memorydevices == nil {
		t.Errorf("GetMemoryDevices() = %q, '%v'", memorydevices, err)
	}
}
