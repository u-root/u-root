// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"testing"

	"github.com/u-root/u-root/pkg/vmtest"
)

func TestScript(t *testing.T) {
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "ShellScript",
		TestCmds: []string{
			"echo HELLO WORLD",
			"shutdown -h",
		},
	})
	defer cleanup()

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(`expected "HELLO WORLD", got error: `, err)
	}
}
