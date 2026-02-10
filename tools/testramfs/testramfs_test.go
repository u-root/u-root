// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestTestRamFS(t *testing.T) {
	if err := Execute(nil, nil, nil, "NO", "amd64", []string{"bad"}); err == nil {
		t.Errorf("execute with bad OS: got nil, want err")
	}
	if err := Execute(nil, nil, nil, "linux", "NOARCH", []string{"bad"}); err == nil {
		t.Errorf("execute with bad ARCH: got nil, want err")
	}

	if err := Execute(nil, nil, nil, "linux", "amd64", []string{"bad"}); err == nil {
		t.Errorf("execute with too few args: got nil, want err")
	}

	if err := Execute(nil, nil, nil, "linux", "amd64", []string{"testramfs", "tmp"}); err == nil {
		t.Errorf("execute with bad u-root directory: got nil, want err")
	}

	// Don't quite know how to do this, as it needs to run u-root
	if false {
		if err := Execute(nil, nil, nil, "linux", "amd64", []string{"testramfs", "../.."}); err != nil {
			t.Errorf("execute with linux/amd64: got %v, want nil", err)
		}
	}
}
