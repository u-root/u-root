// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"strings"
	"testing"
)

func TestCmdline(t *testing.T) {

	exampleCmdLine := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`test-flag test2-flag=8 uroot.uinitflags="a=3 skipfork test2-flag" ` +
		`uroot.initflags="systemd test-flag=3  test2-flag runlevel=2" ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`ipv6.autoconf=0 erst_disable nox2apic crashkernel=128M ` +
		`systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all console=tty0 ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	exampleCmdLineNoInitFlags := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	// Do this once, we'll over-write soon
	once.Do(cmdLineOpener)
	cmdLineReader := strings.NewReader(exampleCmdLine)
	procCmdLine = parse(cmdLineReader)

	if procCmdLine.Err != nil {
		t.Errorf("procCmdLine threw an error: %v", procCmdLine.Err)
	}

	if len(procCmdLine.Raw) != 393 {
		t.Errorf("procCmdLine.Raw wrong length: %v != 417",
			len(procCmdLine.Raw))
	}

	if len(FullCmdLine()) != 393 {
		t.Errorf("FullCmdLine() returned wrong length: %v != 417",
			len(FullCmdLine()))
	}

	if len(procCmdLine.AsMap) != 22 {
		t.Errorf("procCmdLine.Raw wrong length: %v != 22",
			len(procCmdLine.AsMap))
	}

	if ContainsFlag("biosdevname") == false {
		t.Error("couldn't find biosdevname in kernel flags")
	}

	if ContainsFlag("biosname") == true {
		t.Error("could find biosname in kernel flags, but shouldn't")
	}

	if security, present := Flag("security"); !present || security != "selinux" {
		t.Errorf("Flag 'security' is %v instead of 'selinux'", security)
	}

	uinitFlagMap := GetUinitFlagMap()

	if _, present := uinitFlagMap["skipfork"]; !present {
		t.Errorf("Can't find 'skipfork' flag in uinit flags: present == %v",
			present)
	}
	if _, present := uinitFlagMap["madeup"]; present {
		t.Error("Should not find a 'madeup' flag in uinit flags")
	}

	initFlagMap := GetInitFlagMap()
	if testflag, present := initFlagMap["test-flag"]; !present || testflag != "3" {
		t.Errorf("init test-flag == %v instead of test-flag == 3\nMAP: %v", testflag, initFlagMap)
	}

	cmdLineReader = strings.NewReader(exampleCmdLineNoInitFlags)
	procCmdLine = parse(cmdLineReader)
	if initFlagMap = GetInitFlagMap(); len(initFlagMap) != 0 {
		t.Errorf("initFlagMap should be empty, is actually %v", initFlagMap)
	}

}
