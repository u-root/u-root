// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trampoline

import (
	"testing"
)

// TODO: add a unit test that actually verifies that Setup() is
// doing the right thing.
func TestSetupExecutesWithoutCrashing(t *testing.T) {
	var magic, infoAddr, entryPoint uintptr
	somebytes, err := Setup("<path unused>", magic, infoAddr, entryPoint)
	if len(somebytes) != 0 {
		// Expected behavior for linux/amd64
	} else if err != nil {
		// Expected behavior for all other arch/OS's
	} else {
		t.Fatal("didn't expect this")
	}
}
