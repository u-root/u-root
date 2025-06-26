// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import "testing"

func TestRaw(t *testing.T) {
	// It's uncertain what github will ever give us for an environment.
	// If this fails, note it, but it is not an error.
	tm, err := GetTermios(0)
	t.Logf("GetTermios: %v, %v", tm, err)
	if err != nil {
		t.Skipf("No termios available to test on Plan 9: %v", err)
	}

	tm.Close()
	// If that did work, then this has to.
	tt, err := New()
	if err != nil {
		t.Fatalf("New: got %v, want nil", err)
	}

	// The only thing we currently care about
	if _, err := tt.Raw(); err != nil {
		t.Fatalf("Raw: got %v, want nil", err)
	}

	w, err := tt.GetWinSize()
	if err != nil {
		t.Fatalf("GetWinsize: got %v, want nil", err)
	}
	t.Logf("Winsize %d", w)
}
