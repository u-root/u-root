// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/watchdogd"
)

func TestRun(t *testing.T) {
	// it requires root permissions to access the /dev/watchdog device. so the tests need it as well
	guest.SkipIfNotInVM(t)

	// verify there is a watchdog device to pet
	if _, err := os.Stat("/dev/watchdog"); err != nil {
		t.Skip("No /dev/watchdog")
	}

	for _, tt := range []struct {
		name          string
		argv          []string
		expectedError error
	}{
		{
			name:          "none",
			argv:          []string{},
			expectedError: watchdogd.ErrNoCommandSpecified,
		},
		{
			name: "invalid",
			argv: []string{"invalid"},
		},
		{
			name:          "run_empty",
			argv:          []string{"run"},
			expectedError: nil,
		},
		{
			name:          "run_arg1",
			argv:          []string{"run", "dummy"},
			expectedError: nil,
		},
		{
			name:          "run_dev",
			argv:          []string{"run", "--dev", "/dev/watchdog"},
			expectedError: nil,
		},
		{
			name:          "run_dev_invalid_device",
			argv:          []string{"run", "--dev", "DEV"},
			expectedError: os.ErrNotExist,
		},
		{
			name:          "run_timeout",
			argv:          []string{"run", "--timeout", "1s"},
			expectedError: nil,
		},
		{
			name:          "run_pre_timeout",
			argv:          []string{"run", "--pre_timeout", "1s"},
			expectedError: nil,
		},
		{
			name:          "run_keep_alive",
			argv:          []string{"run", "--keep_alive", "1s"},
			expectedError: nil,
		},
		{
			name:          "run_monitors",
			argv:          []string{"run", "--monitors", "oops"},
			expectedError: nil,
		},
		{
			name:          "run_monitors_invalid",
			argv:          []string{"run", "--monitors", "invalid"},
			expectedError: fmt.Errorf("%w: %v", watchdogd.ErrInvalidMonitor, "invalid"),
		},
		{
			name:          "run_full",
			argv:          []string{"run", "--dev", "/dev/watchdog", "--timeout", "1s", "--pre_timeout", "1s", "--keep_alive", "1s", "--monitors", "oops"},
			expectedError: nil,
		},
		{
			name:          "stop",
			argv:          []string{"stop"},
			expectedError: nil,
		},
		{
			name:          "continue",
			argv:          []string{"continue"},
			expectedError: nil,
		},
		{
			name:          "arm",
			argv:          []string{"arm"},
			expectedError: nil,
		},
		{
			name:          "disarm",
			argv:          []string{"disarm"},
			expectedError: nil,
		},
		{
			name:          "invalid_command",
			argv:          []string{"invalid"},
			expectedError: watchdogd.ErrNoCommandSpecified,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.argv); !errors.Is(err, tt.expectedError) {
				t.Errorf("run(): %v, not: %v", err, tt.expectedError)
			}
		})
	}
}
