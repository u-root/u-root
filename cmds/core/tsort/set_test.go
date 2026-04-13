// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"slices"
	"testing"
)

func TestSet(t *testing.T) {
	s := makeSet()

	if s.has("a") {
		t.Errorf(`set %#v: want to not have "a", but did have it`, s)
	}
	if len(s) != 0 {
		t.Errorf(`set %#v: want len of 0, got %d`, s, len(s))
	}

	s.add("b")
	s.add("d")
	s.add("c")
	s.add("c")
	s.add("e")
	s.add("a")

	if !s.has("a") {
		t.Errorf(`set %#v: want to have "a", but did not`, s)
	}
	if !s.has("b") {
		t.Errorf(`set %#v: want to have "b", but did not`, s)
	}
	if !s.has("c") {
		t.Errorf(`set %#v: want to have "c", but did not`, s)
	}
	if !s.has("d") {
		t.Errorf(`set %#v: want to have "d", but did not`, s)
	}
	if !s.has("e") {
		t.Errorf(`set %#v: want to have "e", but did not`, s)
	}
	if s.has("absent-value") {
		t.Errorf(`set %#v: want to not have "absent-value", but did have it`, s)
	}
	if diff := orderInsensitiveDiff(
		slices.Collect(s.all()),
		[]string{"a", "b", "c", "d", "e"},
	); diff != "" {
		t.Errorf("set iterator mismatch (-got +want):\n%s", diff)
	}

	s.remove("a")
	s.remove("e")
	s.remove("c")

	if s.has("a") {
		t.Errorf(`set %#v: want to not have "a", but did have it`, s)
	}
	if !s.has("b") {
		t.Errorf(`set %#v: want to have "b", but did not`, s)
	}
	if s.has("c") {
		t.Errorf(`set %#v: want to not have "c", but did have it`, s)
	}
	if !s.has("d") {
		t.Errorf(`set %#v: want to have "d", but did not`, s)
	}
	if s.has("e") {
		t.Errorf(`set %#v: want to not have "e", but did have it`, s)
	}
	if s.has("absent-value") {
		t.Errorf(`set %#v: want to not have "absent-value", but did have it`, s)
	}
	if diff := orderInsensitiveDiff(
		slices.Collect(s.all()),
		[]string{"b", "d"},
	); diff != "" {
		t.Errorf("set iterator mismatch (-got +want):\n%s", diff)
	}
}
