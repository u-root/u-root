// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortedmap

import (
	"testing"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		data map[int64]int64
		keys []int64
	}{
		{data: map[int64]int64{1: 2, 3: 4}, keys: []int64{1, 3}},
		{data: map[int64]int64{4: 0, 1: 0}, keys: []int64{1, 4}},
	}

	for _, c := range cases {
		m := NewMap()

		for k, v := range c.data {
			m.Insert(k, v)
		}

		// All expected entries exist
		for k, e := range c.data {
			v, ok := m.m[k]
			if !ok {
				t.Errorf("%d not found in %v", k, m.m)
			}
			if v != e {
				t.Errorf("got %d want %d in %v", v, e, m.m)
			}
		}

		if len(m.k) != len(c.keys) {
			t.Errorf("Bad length, got %d, expected %d. %v vs %v", len(m.k), len(c.keys), m.k, c.keys)
		}

		// Key slice is correct
		for i, e := range c.keys {
			if m.k[i] != e {
				t.Errorf("Got %v, expected %v", m.k, c.keys)
				break
			}
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		before Map
		del    []int64
		after  Map
	}{
		{
			before: Map{
				m: map[int64]int64{1: 0, 2: 0},
				k: sortedSlice{1, 2},
			},
			del: []int64{1},
			after: Map{
				m: map[int64]int64{2: 0},
				k: sortedSlice{2},
			},
		},
		{
			before: Map{
				m: map[int64]int64{1: 0, 2: 0},
				k: sortedSlice{1, 2},
			},
			del: []int64{1, 2},
			after: Map{
				m: map[int64]int64{},
				k: sortedSlice{},
			},
		},
	}

	for i := range cases {
		for _, k := range cases[i].del {
			cases[i].before.Delete(k)
		}

		// All expected entries exist
		for k, e := range cases[i].after.m {
			v, ok := cases[i].before.m[k]
			if !ok {
				t.Errorf("%d not found in %v", k, cases[i].before.m)
			}
			if v != e {
				t.Errorf("got %d want %d in %v", v, e, cases[i].before.m)
			}
		}

		if len(cases[i].before.k) != len(cases[i].after.k) {
			t.Errorf("Bad length, got %d, expected %d. %v vs %v", len(cases[i].before.k), len(cases[i].after.k), cases[i].before.k, cases[i].after.k)
		}

		// Key slice is correct
		for i, e := range cases[i].after.k {
			if cases[i].before.k[i] != e {
				t.Errorf("Got %v, expected %v", cases[i].before.k, cases[i].after.k)
				break
			}
		}
	}
}

func TestGet(t *testing.T) {
	m := Map{
		m: map[int64]int64{2: 20, 4: 40},
		k: sortedSlice{2, 4},
	}

	// Simple lookup
	var v int64
	var ok bool
	if v, ok = m.Get(2); !ok {
		t.Errorf("want ok got false for Get(2)")
	}
	if v != 20 {
		t.Errorf("want 20 got %d for Get(2)", v)
	}

	// Non-existent key
	if _, ok = m.Get(3); ok {
		t.Errorf("want not ok got ok for Get(3)")
	}
}

func TestNearestLessEqual(t *testing.T) {
	m := Map{
		m: map[int64]int64{2: 20, 4: 40},
		k: sortedSlice{2, 4},
	}

	// Nothing less than smallest
	_, _, err := m.NearestLessEqual(1)
	if err != ErrNoSuchKey {
		t.Errorf("want err got nil for NLE(1)")
	}

	// Exact match
	k, v, err := m.NearestLessEqual(2)
	if err != nil {
		t.Errorf("want nil got err for NLE(2): %v", err)
	}
	if k != 2 {
		t.Errorf("bad key for NLE(2): want 2 got %d", k)
	}
	if v != 20 {
		t.Errorf("bad value for NLE(2): want 20 got %d", k)
	}

	// One above
	k, v, err = m.NearestLessEqual(3)
	if err != nil {
		t.Errorf("want nil got err for NLE(3): %v", err)
	}
	if k != 2 {
		t.Errorf("bad key for NLE(3): want 2 got %d", k)
	}
	if v != 20 {
		t.Errorf("bad value for NLE(3): want 20 got %d", k)
	}

	// Way above
	k, v, err = m.NearestLessEqual(1000)
	if err != nil {
		t.Errorf("want nil got err for NLE(1000): %v", err)
	}
	if k != 4 {
		t.Errorf("bad key for NLE(1000): want 4 got %d", k)
	}
	if v != 40 {
		t.Errorf("bad value for NLE(1000): want 40 got %d", k)
	}
}

func TestNearestGreaterEqual(t *testing.T) {
	m := Map{
		m: map[int64]int64{2: 20, 4: 40},
		k: sortedSlice{2, 4},
	}

	// Nothing bigger than biggest
	_, _, err := m.NearestGreater(5)
	if err != ErrNoSuchKey {
		t.Errorf("want ErrNoSuchKey got %v for NG(5)", err)
	}

	// One below
	k, v, err := m.NearestGreater(3)
	if err != nil {
		t.Errorf("want nil got err for NG(3): %v", err)
	}
	if k != 4 {
		t.Errorf("bad key for NG(3): want 4 got %d", k)
	}
	if v != 40 {
		t.Errorf("bad value for NG(3): want 40 got %d", k)
	}

	// Way below
	k, v, err = m.NearestGreater(-1000)
	if err != nil {
		t.Errorf("want nil got err for NG(-1000): %v", err)
	}
	if k != 2 {
		t.Errorf("bad key for NG(-1000): want 2 got %d", k)
	}
	if v != 20 {
		t.Errorf("bad value for NG(-1000): want 20 got %d", k)
	}
}
