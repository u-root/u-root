// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

var tests = []struct {
	args    []string
	retCode int
}{
	// Too many args is an error.
	{[]string{"halt", "police"}, 1},
	{[]string{"halt"}, 2},
	{[]string{"-h"}, 2},
	// No args means halt.
	{[]string{}, 2},
	{[]string{"reboot"}, 3},
	{[]string{"-r"}, 3},
	{[]string{"suspend"}, 4},
	{[]string{"-s"}, 4},
}

func TestShutdown(t *testing.T) {
	for i, tt := range tests {
		var retCode int
		c := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--"}, tt.args...)...)
		c.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		_, err := c.Output()
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Errorf("%d. Error running shutdown: %v", i, err)
				continue
			}
			retCode = exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		}
		if retCode != tt.retCode {
			t.Errorf("%d. Want: %d; Got: %d", i, tt.retCode, retCode)
		}
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		t.Logf("just a helper")
		return
	}

	reboot = func(i int) error {
		xval := 1
		switch uint32(i) {
		case unix.LINUX_REBOOT_CMD_POWER_OFF:
			xval = 2
		case unix.LINUX_REBOOT_CMD_RESTART:
			xval = 3
		case unix.LINUX_REBOOT_CMD_SW_SUSPEND:
			xval = 4
		}

		t.Logf("Exit with %#x", i)
		os.Exit(xval)
		return nil
	}

	os.Args = append([]string{"shutdown"}, os.Args[3:]...)
	main()
}
