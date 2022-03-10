// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

// Ensure 0 is returned.
func TestTrue(t *testing.T) {
	if err := runTrue(); err != nil {
		t.Errorf("runTrue():=%q, want nil", err)
	}
}
