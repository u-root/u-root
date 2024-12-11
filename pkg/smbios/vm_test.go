// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package smbios

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/guest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestIntegration(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest("github.com/u-root/u-root/pkg/smbios"),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute*2),
			qemu.ArbitraryArgs("-smbios", "type=2,manufacturer=u-root"),
		),
	)
}

func TestFromSysfs(t *testing.T) {
	guest.SkipIfNotInVM(t)

	info, err := FromSysfs()
	if err != nil || info == nil {
		t.Errorf("FromSysfs() = %q, '%v', want nil", info, err)
	}
}

func TestGetBIOSInfo(t *testing.T) {
	guest.SkipIfNotInVM(t)

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
	guest.SkipIfNotInVM(t)

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
	guest.SkipIfNotInVM(t)

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
	guest.SkipIfNotInVM(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	processorinfo, err := info.GetProcessorInfo()
	if err != nil || processorinfo == nil {
		t.Errorf("GetProcessorInfo() = %q, '%v'", processorinfo, err)
	}
}

func TestGetMemoryDevices(t *testing.T) {
	guest.SkipIfNotInVM(t)

	info, err := FromSysfs()
	if err != nil {
		t.Errorf("FromSysfs as a requirement failed.")
	}

	memorydevices, err := info.GetMemoryDevices()
	if err != nil || memorydevices == nil {
		t.Errorf("GetMemoryDevices() = %q, '%v'", memorydevices, err)
	}
}
