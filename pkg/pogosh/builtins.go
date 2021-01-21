// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"strconv"
)

// DefaultBuiltins lists all the available builtins.
func DefaultBuiltins() map[string]Builtin {
	return map[string]Builtin{
		"exit": BuiltinExit,
	}
}

// BuiltinExit implements the "exit" builtin.
func BuiltinExit(s *State, cmd *Cmd) {
	if len(cmd.argv) >= 2 {
		code, err := strconv.Atoi(cmd.argv[1])
		if err == nil {
			s.Overrides.Exit(int(code))
		}
	}
	s.Overrides.Exit(0)
}
