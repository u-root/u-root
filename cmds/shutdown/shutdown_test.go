// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/sys/unix"
)

type TestGlobals struct {
	flags options
	op    string
	res   uint32
}

func TestAscertainOperationValidArgs(t *testing.T) {
	testFlags := []TestGlobals{
		{options{h: true, r: false}, "", unix.LINUX_REBOOT_CMD_POWER_OFF},
		{options{h: false, r: true}, "", unix.LINUX_REBOOT_CMD_RESTART},
		{options{h: false, r: false}, "suspend", unix.LINUX_REBOOT_CMD_SW_SUSPEND},
		{options{h: false, r: false}, "halt", unix.LINUX_REBOOT_CMD_POWER_OFF},
		{options{h: false, r: false}, "reboot", unix.LINUX_REBOOT_CMD_RESTART},
	}
	for _, test := range testFlags {
		flags, op = test.flags, test.op
		if c, err := ascertainOperation(); err != nil {
			t.Errorf("%s with %+v - %s", err, flags, op)
		} else if c != test.res {
			t.Errorf("expected %s, got %s, with %+v - %s", test.res, c, flags, op)
		}
	}

}

func TestAscertainOperationInvalidArgs(t *testing.T) {
	testFlags := []TestGlobals{
		{options{h: true, r: true}, "", 0},
		{options{h: false, r: true}, "reboot", 0},
		{options{h: true, r: false}, "reboot", 0},
		{options{h: false, r: true}, "suspend", 0},
		{options{h: true, r: false}, "suspend", 0},
		{options{h: false, r: true}, "halt", 0},
		{options{h: true, r: false}, "halt", 0},
	}
	for _, test := range testFlags {
		flags, op = test.flags, test.op
		if _, err := ascertainOperation(); err == nil {
			t.Errorf("expected error, got nothing with %+v - %s", flags, op)
		}
	}

}
