// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"testing"
)

// TestHelloWorld runs an init which prints the string "HELLO WORLD" and exits.
func TestHelloWorld(t *testing.T) {
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/helloworld/uinit",
			"github.com/u-root/u-root/cmds/core/init",
		},
	})
	defer cleanup()

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(`expected "HELLO WORLD", got error: `, err)
	}
}

// TestHelloWorldNegative runs an init which does not print the string "HELLO WORLD".
func TestHelloWorldNegative(t *testing.T) {
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/helloworld/uinit",
			"github.com/u-root/u-root/cmds/core/init",
		},
	})
	defer cleanup()

	if err := q.Expect("GOODBYE WORLD"); err == nil {
		t.Fatal(`expected error, but matched "GOODBYE WORLD"`)
	}
}

func TestScript(t *testing.T) {
	q, cleanup := QEMUTest(t, &Options{
		Name: "ShellScript",
		Cmds: []string{
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/shutdown",
			"github.com/u-root/u-root/cmds/core/echo",
		},
		Uinit: []string{
			"echo HELLO WORLD",
			"shutdown -h",
		},
	})
	defer cleanup()

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(`expected "HELLO WORLD", got error: `, err)
	}
}
