// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"testing"
)

// TestRun test watchdog cli against a regular file, most tests expected to
// return an error, except of "keepalive"
func TestRun(t *testing.T) {
	dev, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatalf("can't create temp file: %v", err)
	}

	tests := []struct {
		name               string
		dev                string
		args               []string
		expectedError      bool
		expectedUsageError bool
	}{
		{
			name:               "no args",
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "device not exists",
			args:          []string{"keepalive"},
			expectedError: true,
		},
		{
			name: "keepalive", // expect no error even file is not watchdog device
			args: []string{"keepalive"},
			dev:  dev.Name(),
		},
		{
			name:               "keepalive wrong args", // expect no error even file is not watchdog device
			args:               []string{"keepalive", "arg"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "settimeout wrong duration",
			args:          []string{"settimeout", "2"},
			dev:           dev.Name(),
			expectedError: true,
		},
		{
			name:          "settimeout",
			args:          []string{"settimeout", "2s"},
			dev:           dev.Name(),
			expectedError: true, // correct command file is not watchdog device
		},
		{
			name:               "settimeout wrong args",
			args:               []string{"settimeout"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "setpretimeout wrong duration",
			args:          []string{"settimeout", "2"},
			dev:           dev.Name(),
			expectedError: true,
		},
		{
			name:               "setpretimeout wrong args",
			args:               []string{"setpretimeout", "2s", "h"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "setpretimeout wrong duration",
			args:          []string{"setpretimeout", "2"},
			dev:           dev.Name(),
			expectedError: true, // correct command, but file is not watchdog device
		},
		{
			name:          "setpretimeout",
			args:          []string{"setpretimeout", "2s"},
			dev:           dev.Name(),
			expectedError: true, // correct command, but file is not watchdog device
		},
		{
			name:               "gettimeout wrong args",
			args:               []string{"gettimeout", "arg"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "gettimeout",
			args:          []string{"gettimeout"},
			dev:           dev.Name(),
			expectedError: true, // correct command, but file is not watchdog device
		},
		{
			name:          "getpretimeout",
			args:          []string{"getpretimeout"},
			dev:           dev.Name(),
			expectedError: true,
		},
		{
			name:               "getpretimeout wrong args",
			args:               []string{"getpretimeout", "arg"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "gettimeleft",
			args:          []string{"gettimeleft"},
			dev:           dev.Name(),
			expectedError: true,
		},
		{
			name:               "gettimeleft wrong args",
			args:               []string{"gettimeleft", "arg"},
			dev:                dev.Name(),
			expectedError:      true,
			expectedUsageError: true,
		},
		{
			name:          "unknown command",
			args:          []string{"unknown"},
			dev:           dev.Name(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		err := runCommand(test.dev, test.args...)
		if test.expectedError && err == nil {
			t.Error("expected error, got nil")
		}
		if !test.expectedError && err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if test.expectedUsageError && !errors.Is(err, errUsage) {
			t.Errorf("expected %v, got %v", errUsage, err)
		}
	}
}
