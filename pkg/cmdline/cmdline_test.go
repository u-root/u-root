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
		`test-flag test2-flag=8 ` +
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

	wantLen := len(exampleCmdLine)
	if len(procCmdLine.Raw) != wantLen {
		t.Errorf("procCmdLine.Raw wrong length: %v != %d",
			len(procCmdLine.Raw), wantLen)
	}

	if len(FullCmdLine()) != wantLen {
		t.Errorf("FullCmdLine() returned wrong length: %v != %d",
			len(FullCmdLine()), wantLen)
	}

	if len(procCmdLine.AsMap) != 21 {
		t.Errorf("procCmdLine.AsMap wrong length: %v != 21",
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

func TestCmdlineModules(t *testing.T) {
	exampleCmdlineModules := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`my_module.flag1=8 my-module.flag2-string=hello ` +
		`otherMod.opt1=world otherMod.opt_2=22-22`

	once.Do(cmdLineOpener)
	cmdLineReader := strings.NewReader(exampleCmdlineModules)
	procCmdLine = parse(cmdLineReader)

	if procCmdLine.Err != nil {
		t.Errorf("procCmdLine threw an error: %v", procCmdLine.Err)
	}

	// Check flags using contains to not rely on map iteration order
	flags := FlagsForModule("my-module")
	if !strings.Contains(flags, "flag1=8 ") || !strings.Contains(flags, "flag2_string=hello ") {
		t.Errorf("my-module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}
	flags = FlagsForModule("my_module")
	if !strings.Contains(flags, "flag1=8 ") || !strings.Contains(flags, "flag2_string=hello ") {
		t.Errorf("my_module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}

	flags = FlagsForModule("otherMod")
	if !strings.Contains(flags, "opt1=world ") || !strings.Contains(flags, "opt_2=22-22 ") {
		t.Errorf("my_module flags got: %v, want opt1=world opt_2=22-22 ", flags)
	}
}
