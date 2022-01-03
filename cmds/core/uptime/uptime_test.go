// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"
)

var testTime = time.Date(0o001, 1, 15, 5, 35, 49, 0, time.UTC)

func invalidDurationError(d string) string {
	_, err := time.ParseDuration(d)
	return err.Error()
}

func TestUptime(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		uptime *time.Time
		err    string
	}{
		{
			name:   "goodInput",
			input:  "1229749 1422244",
			uptime: &testTime,
			err:    "",
		},
		{
			name:   "badDataInput",
			input:  "string",
			uptime: nil,
			err:    invalidDurationError("strings"),
		},
		{
			name:   "emptyDataInput",
			input:  "",
			uptime: nil,
			err:    "error:the contents of proc/uptime we are trying to read are empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUptime, err := uptime(test.input)
			if err == nil && test.err != "" {
				t.Errorf("uptime(%q) err = nil, want %q", test.input, test.err)
			} else if err != nil && err.Error() != test.err {
				t.Errorf("uptime(%q) err = %q, want %q", test.input, err.Error(), test.err)
			}
			if gotUptime == nil && test.uptime != nil {
				t.Errorf("uptime(%q) = nil, want %v", test.input, *test.uptime)
			} else if gotUptime != nil && test.uptime != nil && *gotUptime != *test.uptime {
				t.Errorf("uptime(%q) = %v, want %v", test.input, *gotUptime, *test.uptime)
			} else if gotUptime != nil && test.uptime == nil {
				t.Errorf("uptime(%q) = %v, want nil", test.input, *gotUptime)
			}
		})
	}
}

func TestLoadAverage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		loadAverage string
		err         string
	}{
		{"goodInput", "0.60 0.70 0.74", "0.60,0.70,0.74", ""},
		{"badDataInput", "1.00 2.00", "", "error:invalid contents:the contents of proc/loadavg we are trying to process contain less than the required 3 loadavgs"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loadAverage, err := loadavg(test.input)
			if err == nil && test.err != "" {
				t.Errorf("loadavg(%q) err = nil, want %q", test.input, test.err)
			} else if err != nil && err.Error() != test.err {
				t.Errorf("loadavg(%q) err = %q, want %q", test.input, err.Error(), test.err)
			}
			if loadAverage == "" && test.loadAverage != "" {
				t.Errorf("loadavg(%q) = \"\", want %v", test.input, test.loadAverage)
			} else if loadAverage != "" && test.loadAverage != "" && loadAverage != test.loadAverage {
				t.Errorf("loadavg(%q) = %v, want %v", test.input, loadAverage, test.loadAverage)
			} else if loadAverage != "" && test.loadAverage == "" {
				t.Errorf("loadavg(%q) = %v, want \"\"", test.input, loadAverage)
			}
		})
	}
}
