// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"io"
	"os"
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

	c := parse(strings.NewReader(exampleCmdLine))
	wantLen := len(exampleCmdLine)
	if len(c.Raw) != wantLen {
		t.Errorf("c.Raw wrong length: %v != %d", len(c.Raw), wantLen)
	}

	if len(c.AsMap) != 21 {
		t.Errorf("c.AsMap wrong length: %v != 21", len(c.AsMap))
	}

	if c.ContainsFlag("biosdevname") == false {
		t.Errorf("couldn't find biosdevname in kernel flags: map is %v", c.AsMap)
	}

	if c.ContainsFlag("biosname") == true {
		t.Error("could find biosname in kernel flags, but shouldn't")
	}

	if security, present := c.Flag("security"); !present || security != "selinux" {
		t.Errorf("Flag 'security' is %v instead of 'selinux'", security)
	}

	initFlagMap := c.GetInitFlagMap()
	if testflag, present := initFlagMap["test-flag"]; !present || testflag != "3" {
		t.Errorf("init test-flag == %v instead of test-flag == 3\nMAP: %v", testflag, initFlagMap)
	}

	c = parse(strings.NewReader(exampleCmdLineNoInitFlags))
	if initFlagMap = c.GetInitFlagMap(); len(initFlagMap) != 0 {
		t.Errorf("initFlagMap should be empty, is actually %v", initFlagMap)
	}
}

func TestCmdlineModules(t *testing.T) {
	exampleCmdlineModules := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`my_module.flag1=8 my-module.flag2-string=hello ` +
		`otherMod.opt1=world otherMod.opt_2=22-22`

	c := parse(strings.NewReader(exampleCmdlineModules))

	// Check flags using contains to not rely on map iteration order
	flags := c.FlagsForModule("my-module")
	if !strings.Contains(flags, "flag1=8 ") || !strings.Contains(flags, "flag2_string=hello ") {
		t.Errorf("my-module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}
	flags = c.FlagsForModule("my_module")
	if !strings.Contains(flags, "flag1=8 ") || !strings.Contains(flags, "flag2_string=hello ") {
		t.Errorf("my_module flags got: %v, want flag1=8 flag2_string=hello ", flags)
	}

	flags = c.FlagsForModule("otherMod")
	if !strings.Contains(flags, "opt1=world ") || !strings.Contains(flags, "opt_2=22-22 ") {
		t.Errorf("my_module flags got: %v, want opt1=world opt_2=22-22 ", flags)
	}
}

// Functional tests are done elsewhere. This test is purely to
// call the package level functions.
func TestCmdLineClassic(t *testing.T) {
	t.Skipf("This fails in integration for reasons still unknown")
	c := getCmdLine()
	if c.Err != nil {
		t.Skipf("getCmdLine(): got %v, want nil, skipping test", c.Err)
	}

	c = cmdLine("/proc/cmdlinexyzzy")
	// There is no good reason for an open like this to succeed.
	// But, in virtual environments, it seems to at times.
	// Just log it.
	if c.Err == nil {
		t.Skipf(`cmdLine("/proc/cmdlinexyzzy"): got nil, want %v, skipping test`, os.ErrNotExist)
	}
	NewCmdLine()
	FullCmdLine()
	// These functions call functions that are already tested, but
	// this is our way of boosting coverage :-)
	FlagsForModule("something")
	GetUinitArgs()
	GetInitFlagMap()
	Flag("noflag")
	ContainsFlag("noflag")
}

type badreader struct{}

// Read implements io.Reader, always returning io.ErrClosedPipe
func (*badreader) Read([]byte) (int, error) {
	// Interesting. If you return a -1 for the length,
	// it tickles a bug in io.ReadAll. It uses the returned
	// length BEFORE seeing if there was an error.
	// Note to self: file an issue on Go.
	return 0, io.ErrClosedPipe
}

func TestBadRead(t *testing.T) {
	if err := parse(&badreader{}); err == nil {
		t.Errorf("parse(&badreader{}): got nil, want %v", io.ErrClosedPipe)
	}
}
