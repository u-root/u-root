// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func DisabledTestTTY(t *testing.T) {
	tty()
	foreground()
	t.Logf("tty testing done")
}
