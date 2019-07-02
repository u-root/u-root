// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"
)

var (
	testTime = time.Date(0001, 1, 15, 5, 35, 49, 0, time.UTC)
)

func TestUptime(t *testing.T) {
	var tests = []struct {
		name   string
		input  string
		uptime *time.Time
		err    string
	}{
		{"goodInput", "1229749 1422244", &testTime, ""},
		{"badDataInput", "string", nil, "error time: invalid duration strings"},
		{"emptyDataInput", "", nil, "error:the contents of proc/uptime we are trying to read are empty"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUptime, err := uptime(test.input)
			if err == nil && test.err != "" {
				t.Errorf("Error was not returned:got nil,want %v", test.err)
			} else if err != nil && err.Error() != test.err {
				t.Errorf("Mismatched Error returned :got %v ,want %v", err.Error(), test.err)
			}
			if gotUptime == nil && test.uptime != nil {
				t.Errorf("Error mismatched uptime :got nil ,want %v", *test.uptime)
			} else if gotUptime != nil && test.uptime != nil && *gotUptime != *test.uptime {
				t.Errorf("Error mismatched uptime :got %v , want %v", *gotUptime, *test.uptime)
			} else if gotUptime != nil && test.uptime == nil {
				t.Errorf("Error mismatched uptime :got %v , want nil", *gotUptime)
			}
		})
	}
}

func TestLoadAverage(t *testing.T) {
	var tests = []struct {
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
				t.Errorf("Error was not returned:got nil, want %s", test.err)
			} else if err != nil && err.Error() != test.err {
				t.Errorf("Mismatched Error returned,got %s ,want %s", err.Error(), test.err)
			}
			if loadAverage == "" && test.loadAverage != "" {
				t.Errorf("Error mismatched loadaverage: got '', want %s", test.loadAverage)
			} else if loadAverage != "" && test.loadAverage != "" && loadAverage != test.loadAverage {
				t.Errorf("Error mismatched loadaverage :got %s ,want %s", loadAverage, test.loadAverage)
			} else if loadAverage != "" && test.loadAverage == "" {
				t.Errorf("Error mismatched loadaverage :got %s ,want ''", loadAverage)
			}
		})
	}
}
