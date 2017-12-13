// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import "testing"

var m = map[string]int64{
	"B":   1,
	"KiB": 1024,
	"MiB": 1024 * 1024,
}

var unitTests = []struct {
	in, out      string
	value        int64
	explicitSign Sign
	err          bool
}{
	{
		// Invalid
		in:  "",
		err: true,
	}, {
		// With invalid unit
		in:  "23K",
		err: true,
	}, {
		// With valid unit
		in:    "23KiB",
		out:   "23KiB",
		value: 23 * 1024,
	}, {
		// With valid without unit
		in:    "23",
		out:   "23B",
		value: 23,
	}, {
		// With valid unit, mutliple of another unit
		in:    "1024KiB",
		out:   "1MiB",
		value: 1024 * 1024,
	}, {
		// With valid unit, mutliple of another unit, explicit +
		in:           "+1024KiB",
		out:          "+1MiB",
		value:        1024 * 1024,
		explicitSign: Positive,
	}, {
		// With valid unit, mutliple of another unit, explicit -
		in:           "-1024KiB",
		out:          "-1MiB",
		value:        -1024 * 1024,
		explicitSign: Negative,
	},
}

// TestUnit implements a table-driven test.
func TestUnit(t *testing.T) {
	u, err := NewUnit(m)
	if err != nil {
		t.Fatalf("Expected that a valid unit map does not return an error: %v\n", err)
	}

	for _, test := range unitTests {
		v, err := u.ValueFromString(test.in)
		if test.err && err == nil {
			t.Fatalf("Expected that rest '%s' returns an error\n", test.in)
		} else if !test.err && err != nil {
			t.Fatalf("Expected that rest '%s' does not return an error\n", test.in)
		} else if test.err && err != nil {
			continue
		}

		if v.Value != test.value {
			t.Fatalf("Expected that rest '%s' returns %d as bytes, but got %d\n", test.in, test.value, v.Value)
		}

		if v.String() != test.out {
			t.Fatalf("Expected that rest '%s' returns %s as string, but got %s\n", test.in, test.out, v)
		}

		if v.ExplicitSign != test.explicitSign {
			t.Fatalf("Expected that rest '%s' has explicit Sign %d, but got %d\n", test.in, test.explicitSign, v.ExplicitSign)
		}
	}
}
