// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msr

import (
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/testutil"
)

func TestParseCPUs(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		cpus   CPUs
		errStr string
	}{
		{
			name:   "one core",
			input:  "0",
			cpus:   []uint64{0},
			errStr: "",
		},
		{
			name:   "eight cores",
			input:  "0-7",
			cpus:   []uint64{0, 1, 2, 3, 4, 5, 6, 7},
			errStr: "",
		},
		{
			name:   "split cores",
			input:  "0-2,4-7",
			cpus:   []uint64{0, 1, 2, 4, 5, 6, 7},
			errStr: "",
		},
		{
			name:   "duplicates",
			input:  "0-2,0-2",
			cpus:   []uint64{0, 1, 2},
			errStr: "",
		},
		{
			name:   "invalid range",
			input:  "7-0",
			cpus:   nil,
			errStr: "invalid cpu range, upper bound greater than lower: 7-0",
		},
		{
			name:   "non numbers",
			input:  "a",
			cpus:   nil,
			errStr: "unknown cpu range: a, failed to parse strconv.ParseUint: parsing \"a\": invalid syntax",
		},
		{
			name:   "non numbers 2",
			input:  "a-b",
			cpus:   nil,
			errStr: "unknown cpu range: a-b, failed to parse strconv.ParseUint: parsing \"a\": invalid syntax",
		},
		{
			name:   "weird range",
			input:  "1-2-3",
			cpus:   nil,
			errStr: "unknown cpu range: 1-2-3",
		},
		{
			name:   "weirder range",
			input:  "-1",
			cpus:   nil,
			errStr: "unknown cpu range: -1, failed to parse strconv.ParseUint: parsing \"\": invalid syntax",
		},
		{
			name:   "empty string",
			input:  "",
			cpus:   nil,
			errStr: "no cpus found, input was ",
		},
		{
			name:   "comma string",
			input:  ",,,,",
			cpus:   nil,
			errStr: "no cpus found, input was ,,,,",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o, err := parseCPUs(test.input)
			if e := testutil.CheckError(err, test.errStr); e != nil {
				t.Error(e)
			}
			if test.cpus != nil && o == nil {
				t.Errorf("Expected cpus %v, got nil", test.cpus)
			} else if test.cpus == nil && o != nil {
				t.Errorf("Expected nil cpus, got %v", o)
			} else if test.cpus != nil && o != nil {
				if len(test.cpus) != len(o) {
					t.Errorf("Mismatched output: got %v, want %v", o, test.cpus)
				} else {
					for i := range o {
						if test.cpus[i] != o[i] {
							t.Errorf("Mismatched output, got %v, want %v", o, test.cpus)
						}
					}
				}
			}
		})
	}
}

func TestCPUsString(t *testing.T) {
	tests := []struct {
		name   string
		cpus   CPUs
		output string
	}{
		{
			name:   "no cpus",
			cpus:   []uint64{},
			output: "nil",
		},
		{
			name:   "1 cpu",
			cpus:   []uint64{1},
			output: "1",
		},
		{
			name:   "2 cpu",
			cpus:   []uint64{1, 2},
			output: "1-2",
		},
		{
			name:   "3 cpu",
			cpus:   []uint64{1, 2},
			output: "1-2",
		},
		{
			name:   "3 cpu",
			cpus:   []uint64{1, 2, 3},
			output: "1-3",
		},
		{
			name:   "3 noncontinuous cpu",
			cpus:   []uint64{1, 2, 4},
			output: "1-2,4",
		},
		{
			name:   "large sequence",
			cpus:   []uint64{1, 2, 4, 6, 7, 8, 9, 10, 11, 12, 13},
			output: "1-2,4,6-13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.cpus.String()
			if s != test.output {
				t.Errorf("CPUs(%v).String() == %q; want %q", []uint64(test.cpus), s, test.output)
			}
		})
	}
}

// This is a hard one to test. But for many systems, 0x3a is a good bet.
// but this is of necessity not a complete test! We don't want to set an MSR
// as part of a test, it might cause real trouble.
func TestTestAndSet(t *testing.T) {
	guest.SkipIfNotInVM(t)

	c, err := AllCPUs()
	if err != nil {
		t.Fatalf("AllCPUs: got %v,want nil", err)
	}
	r := IntelIA32FeatureControl
	vals, errs := r.Read(c)
	if errs != nil {
		t.Skipf("Skipping test, can't read %s", r)
	}

	if errs = r.Test(c, 0, 0); errs != nil {
		t.Errorf("Test with good val: got %v, want nil", errs)
	}

	// clear every set bit, set every clear bit, this should not go well.
	if errs = r.Test(c, vals[0], ^vals[0]); errs == nil {
		t.Errorf("Test with bad clear/set/val 0x%#x: got nil, want err", vals[0])
	}
}

func TestLocked(t *testing.T) {
	// Passing this is basically optional, but at least try.
	if err := Locked(); err != nil {
		t.Logf("(warning only) Verify GenuineIntel: got %v, want nil", err)
	}
}
