// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"strings"
	"testing"
)

func TestRemoveFilter(t *testing.T) {
	toRemove := []string{"remove-1", "remove_2"}

	cl := `keep=5 remove_1=wontbethere remove-2=nomore keep2`
	want := `keep=5 keep2`

	got := removeFilter(cl, toRemove)
	if got != want {
		t.Errorf("removeFilter(%v,%v) = %v, want %v", cl, toRemove, got, want)
	}
}

func TestUpdateFilter(t *testing.T) {
	exampleCmdLine := `BOOT_IMAGE=/vmlinuz-4.11.2 ro ` +
		`test-flag test2-flag=8 uroot.uinitflags="a=3 skipfork test2-flag" ` +
		`uroot.initflags="systemd test-flag=3  test2-flag runlevel=2" ` +
		`root=LABEL=/ biosdevname=0 net.ifnames=0 fsck.repair=yes ` +
		`ipv6.autoconf=0 erst_disable nox2apic crashkernel=128M ` +
		`systemd.unified_cgroup_hierarchy=1 cgroup_no_v1=all console=tty0 ` +
		`console=ttyS0,115200 security=selinux selinux=1 enforcing=0`

	c := parse(strings.NewReader(exampleCmdLine))

	toRemove := []string{"console", "earlyconsole"}
	toReuse := []string{"console", "not-present"}
	toAppend := "append=me"

	cl := `keep=5 console=ttyS1 keep2 earlyconsole=ttyS1`
	want := `keep=5 keep2 append=me console=ttyS0,115200`

	filter := NewUpdateFilter(toAppend, toRemove, toReuse)
	got := filter.Update(c, cl)
	if got != want {
		t.Errorf("Update(%q) = %q, want %q", cl, got, want)
	}
}
