// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := set{}
	s.add("a")
	s.add("b")
	s.add("c")

	if !s.has("a") {
		t.Errorf(`set %#v: want to have "a", but did not`, s)
	}
	if !s.has("b") {
		t.Errorf(`set %#v: want to have "b", but did not`, s)
	}
	if !s.has("c") {
		t.Errorf(`set %#v: want to have "c", but did not`, s)
	}
	if s.has("absent-value") {
		t.Errorf(`set %#v: want to not have "absent-value", but did`, s)
	}
}
