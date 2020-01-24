// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"strings"
	"testing"
)

func TestCmdlineDuplicates(t *testing.T) {
	cmdline := "test-flag=a test_flag=b test_flag=c"
	c := Parse(cmdline)

	want := "c"
	if got, _ := c.Flag("test_flag"); got != want {
		t.Errorf("(%s).Flag(test_flag) = %s, want %s", cmdline, got, want)
	}
	if got, _ := c.Flag("test-flag"); got != want {
		t.Errorf("(%s).Flag(test-flag) = %s, want %s", cmdline, got, want)
	}
	if _, ok := c.Flag("testflag"); ok {
		t.Errorf("(%s).Flag(testflag) = %t, want false", cmdline, ok)
	}
}

func TestCmdlineNoValue(t *testing.T) {
	cmdline := "ro foo bar=baz"
	c := Parse(cmdline)

	contains := []string{"ro", "foo", "bar"}
	for _, flag := range contains {
		if ok := c.Contains(flag); !ok {
			t.Errorf("(%s).Contains(%s) = %t, want true", cmdline, flag, ok)
		}
	}
}

func TestCmdlineQuoted(t *testing.T) {
	cmdline := `ro root=/dev/sda1 ipv6.autoconf=0 uroot.initflags="systemd test-flag=3  test2-flag runlevel=2"`
	c := Parse(cmdline)

	want := "systemd test-flag=3  test2-flag runlevel=2"
	if got, _ := c.Flag("uroot.initflags"); got != want {
		t.Errorf("(%s).Flag(uroot.initflags) = %v, want %v", cmdline, got, want)
	}
}

func TestAppendPrepend(t *testing.T) {
	cmdline := "test-flag=a test_flag=b test_flag=c"
	c := Parse(cmdline)

	c.Append("test-flag=3")

	want := "3"
	raw := "test-flag=a test_flag=b test_flag=c test-flag=3"
	if got := c.String(); got != raw {
		t.Errorf("String() = %v want %v", got, raw)
	}
	if got, _ := c.Flag("test_flag"); got != want {
		t.Errorf("Flag(test_flag) = %v want %v", got, want)
	}

	c.Prepend("test-flag=5")

	raw = "test-flag=5 test-flag=a test_flag=b test_flag=c test-flag=3"
	if got := c.String(); got != raw {
		t.Errorf("String() = %v want %v", got, raw)
	}
	// Still 3.
	if got, _ := c.Flag("test_flag"); got != want {
		t.Errorf("Flag(test_flag) = %v want %v", got, want)
	}

	nc := NewCmdline()
	nc.Append("test")
	if s := nc.String(); s != "test" {
		t.Errorf("String() = %v, want test", s)
	}

	nc2 := NewCmdline()
	nc2.Prepend("test")
	if s := nc2.String(); s != "test" {
		t.Errorf("String() = %v, want test", s)
	}
}

func TestCmdline(t *testing.T) {
	cmdline := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`test-flag test2-flag=8 ` +
		`uroot.initflags="systemd test-flag=3  test2-flag runlevel=2" ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`ipv6.autoconf=0 erst_disable nox2apic crashkernel=128M ` +
		`systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all console=tty0 ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	c := Parse(cmdline)
	// There's 20 args in there, but console appears twice.
	if len(c.asMap) != 19 {
		t.Errorf("asMap wrong length: %v != 19", len(c.asMap))
	}

	if c.Contains("biosdevname") == false {
		t.Error("couldn't find biosdevname in kernel flags")
	}

	if c.Contains("biosname") == true {
		t.Error("could find biosname in kernel flags, but shouldn't")
	}

	if security, ok := c.Flag("security"); !ok || security != "selinux" {
		t.Errorf("Flag 'security' is %v instead of 'selinux'", security)
	}
}

func TestInitMap(t *testing.T) {
	cmdline := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`test-flag test2-flag=8 ` +
		`uroot.initflags="systemd test-flag=3  test2-flag runlevel=2" ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`ipv6.autoconf=0 erst_disable nox2apic crashkernel=128M ` +
		`systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all console=tty0 ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	hostOnce.Do(func() {})
	hostCmdline = Parse(cmdline)

	initFlagMap := GetInitFlagMap()
	if testflag, present := initFlagMap["test-flag"]; !present || testflag != "3" {
		t.Errorf("init test-flag == %v instead of test-flag == 3\nMAP: %v", testflag, initFlagMap)
	}

	cmdlineNoInit := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`
	hostCmdline = Parse(cmdlineNoInit)
	if initFlagMap = GetInitFlagMap(); len(initFlagMap) != 0 {
		t.Errorf("initFlagMap should be empty, is actually %v", initFlagMap)
	}
}

func TestCmdlineModules(t *testing.T) {
	cmdline := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`my_module.flag1=8 my-module.flag2-string=hello ` +
		`otherMod.opt1=world otherMod.opt_2=22-22`
	c := Parse(cmdline)

	// Check flags using contains to not rely on map iteration order
	flags := c.FlagsForModule("my-module")
	if !strings.Contains(flags, "flag1=8") || !strings.Contains(flags, "flag2_string=hello") {
		t.Errorf("my-module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}
	flags = c.FlagsForModule("my_module")
	if !strings.Contains(flags, "flag1=8") || !strings.Contains(flags, "flag2_string=hello") {
		t.Errorf("my_module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}

	flags = c.FlagsForModule("otherMod")
	if !strings.Contains(flags, "opt1=world") || !strings.Contains(flags, "opt_2=22-22") {
		t.Errorf("otherMod flags got: %v, want opt1=world opt_2=22-22 ", flags)
	}
}
