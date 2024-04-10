// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
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
			err:    errUptimeFormat.Error(),
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
		{"goodInput", "0.60 0.70 0.74", "0.60, 0.70, 0.74", ""},
		{"badDataInput", "1.00 2.00", "", errAverageFormat.Error()},
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

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()
	t.Run("error uptime open", func(t *testing.T) {
		err := run(nil, "filenotexists", "filenotexists")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("error uptime empty file", func(t *testing.T) {
		path := tmpDir + "/uptime-error-0"
		err := os.WriteFile(path, nil, 0o664)
		if err != nil {
			t.Fatal(err)
		}
		err = run(nil, path, "filenotexists")
		if !errors.Is(err, errUptimeFormat) {
			t.Fatal(err)
		}
	})
	t.Run("error uptime parsin file", func(t *testing.T) {
		path := tmpDir + "/uptime-error-1"
		err := os.WriteFile(path, []byte("wrong"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		err = run(nil, path, "filenotexists")
		if err == nil {
			t.Error("expected err got nil")
		}
	})
	t.Run("error average open", func(t *testing.T) {
		path := tmpDir + "/uptime-error-2"
		err := os.WriteFile(path, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}

		err = run(nil, path, "filenotexists")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("error average empty file", func(t *testing.T) {
		pathUptime := tmpDir + "/uptime-error-3"
		err := os.WriteFile(pathUptime, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		pathAvg := tmpDir + "/avg-error-0"
		err = os.WriteFile(pathAvg, nil, 0o664)
		if err != nil {
			t.Fatal(err)
		}
		err = run(nil, pathUptime, pathAvg)
		if !errors.Is(err, errAverageFormat) {
			t.Errorf("expected %v, got %v", errAverageFormat, err)
		}
	})
	t.Run("success", func(t *testing.T) {
		pathUptime := tmpDir + "/uptime"
		err := os.WriteFile(pathUptime, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		pathAvg := tmpDir + "/avg"
		err = os.WriteFile(pathAvg, []byte("0.08 0.09 0.10"), 0o664)
		if err != nil {
			t.Fatal(err)
		}

		stdout := &bytes.Buffer{}
		err = run(stdout, pathUptime, pathAvg)
		if err != nil {
			t.Errorf("expected nil got %v", err)
		}

		if !strings.Contains(stdout.String(), "0.08, 0.09, 0.10") {
			t.Error("expected to see load averages")
		}
	})
}
