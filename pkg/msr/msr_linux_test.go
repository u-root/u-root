// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package msr

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestParseCPUs(t *testing.T) {
	var tests = []struct {
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
