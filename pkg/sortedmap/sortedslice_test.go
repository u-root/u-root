// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortedmap

import (
	"testing"
)

func TestIntSliceInsert(t *testing.T) {
	cases := []struct {
		insert   []int64
		expected []int64
	}{
		{insert: []int64{1, 2, 3}, expected: []int64{1, 2, 3}},
		{insert: []int64{3, 2, 1}, expected: []int64{1, 2, 3}},
		{insert: []int64{-4, 20, -10, 5}, expected: []int64{-10, -4, 5, 20}},
	}

	for _, c := range cases {
		l := make(sortedSlice, 0)

		for _, v := range c.insert {
			l.Insert(v)
		}

		if len(l) != len(c.expected) {
			t.Errorf("Bad length, got %d, expected %d. %v vs %v", len(l), len(c.expected), l, c.expected)
		}

		for i, e := range c.expected {
			if l[i] != e {
				t.Errorf("Got %v, expected %v", l, c.expected)
				break
			}
		}
	}
}

func TestIntSliceDelete(t *testing.T) {
	cases := []struct {
		slice    sortedSlice
		del      []int64
		expected []int64
	}{
		{slice: sortedSlice{1, 2, 3}, del: []int64{2}, expected: []int64{1, 3}},
		{slice: sortedSlice{1, 2, 3}, del: []int64{3, 2, 1}, expected: []int64{}},
		{slice: sortedSlice{-10, -4, 5, 20}, del: []int64{-4}, expected: []int64{-10, 5, 20}},
	}

	for _, c := range cases {
		for _, v := range c.del {
			c.slice.Delete(v)
		}

		if len(c.slice) != len(c.expected) {
			t.Errorf("Bad length, got %d, expected %d. %v vs %v", len(c.slice), len(c.expected), c.slice, c.expected)
		}

		for i, e := range c.expected {
			if c.slice[i] != e {
				t.Errorf("Got %v, expected %v", c.slice, c.expected)
				break
			}
		}
	}
}
