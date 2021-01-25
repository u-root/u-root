// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import "testing"

func TestMagicToName(t *testing.T) {
	for _, m := range magics {
		s, err := MagicToName(m.magic)
		if err != nil {
			t.Errorf("MagicToName(%#x): got (%q, %v), want (%q, nil)", m.magic, s, err, m.name)
		}
		if s != m.name {
			t.Errorf("MagicToName(%#x): got %q, want %q", m.magic, s, m.name)
		}
	}
	// Test something bogus
	_, err := MagicToName(0)
	if err == nil {
		t.Errorf("MagicToName(0): got nil, want err")
	}
}

func TestNameToMagic(t *testing.T) {
	for _, m := range magics {
		s, err := NameToMagic(m.name)
		if err != nil {
			t.Errorf("NameToMagic(%q): got (%#x, %v), want (%#x, nil)", m.name, s, err, m.magic)
		}
		if s != m.magic {
			t.Errorf("NameToMagic(%q): got %#x, want %#x", m.name, s, m.magic)
		}
	}
	// Test something bogus
	_, err := NameToMagic("bogus")
	if err == nil {
		t.Errorf("NameToMagic(\"bogus\"): got nil, want err")
	}
}
