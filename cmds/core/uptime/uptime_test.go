// Copyright 2019-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUptime(t *testing.T) {
	testTime := time.Date(0o001, 1, 15, 5, 35, 49, 0, time.UTC)

	tests := []struct {
		err    error
		uptime *time.Time
		name   string
		input  string
	}{
		{
			name:   "goodInput",
			input:  "1229749 1422244",
			uptime: &testTime,
		},
		{
			name:   "emptyDataInput",
			uptime: nil,
			err:    errUptimeFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUptime, err := uptime(test.input)
			if !errors.Is(err, test.err) {
				t.Errorf("uptime(%q) err = %v, want %v", test.input, err, test.err)
				return
			}

			if test.err != nil {
				return
			}

			if !gotUptime.Equal(*test.uptime) {
				t.Errorf("uptime(%q) = %v, want %v", test.input, *gotUptime, *test.uptime)
			}
		})
	}
}

func TestLoadAverage(t *testing.T) {
	tests := []struct {
		err         error
		name        string
		input       string
		loadAverage string
	}{
		{
			name:        "goodInput",
			input:       "0.60 0.70 0.74",
			loadAverage: "0.60, 0.70, 0.74",
			err:         nil,
		},
		{
			name:  "badDataInput",
			input: "1.00 2.00",
			err:   errAverageFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loadAverage, err := loadavg(test.input)
			if !errors.Is(err, test.err) {
				t.Errorf("loadavg(%q) err = %v, want %v", test.input, err, test)
				return
			}
			if loadAverage != test.loadAverage {
				t.Errorf("loadavg(%q) = \"\", want %v", test.input, test.loadAverage)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()
	notExists := filepath.Join(tmpDir, "filenotexists")
	t.Run("error uptime open", func(t *testing.T) {
		err := run(nil, notExists, notExists)
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("error uptime empty file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "uptime-error-0")
		err := os.WriteFile(path, nil, 0o664)
		if err != nil {
			t.Fatal(err)
		}
		err = run(nil, path, notExists)
		if !errors.Is(err, errUptimeFormat) {
			t.Fatal(err)
		}
	})
	t.Run("error uptime parsing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "/uptime-error-1")
		err := os.WriteFile(path, []byte("wrong"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		err = run(nil, path, notExists)
		if err == nil {
			t.Error("expected err got nil")
		}
	})
	t.Run("error average open", func(t *testing.T) {
		path := filepath.Join(tmpDir, "uptime-error-2")
		err := os.WriteFile(path, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}

		err = run(nil, path, notExists)
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("error average empty file", func(t *testing.T) {
		pathUptime := filepath.Join(tmpDir, "/uptime-error-3")
		err := os.WriteFile(pathUptime, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		pathAvg := filepath.Join(tmpDir, "/avg-error-0")
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
		pathUptime := filepath.Join(tmpDir, "uptime")
		err := os.WriteFile(pathUptime, []byte("1462.14 5746.97"), 0o664)
		if err != nil {
			t.Fatal(err)
		}
		pathAvg := filepath.Join(tmpDir, "avg")
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
