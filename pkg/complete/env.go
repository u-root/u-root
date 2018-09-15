// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"fmt"
	"os"
	"strings"
)

// NewEnvCompleter creates a MultiCompleter consisting of one
// or more FileCompleters. It is given an environment variable,
// which it splits on :. If there are only zero entries,
// it returns an error; else it returns a MultiCompleter.
// N.B. it does *not* check for whether a directory exists
// or not; directories can come and go.
func NewEnvCompleter(s string) (Completer, error) {
	e := os.Getenv(s)
	Debug("NewEnvCompleter: path %q has value %q", s, e)
	if e == "" {
		return nil, ErrEmptyEnv
	}
	dirs := strings.Split(e, ":")
	if len(dirs) == 0 {
		return nil, fmt.Errorf("%s is empty", s)
	}
	Debug("Build completer for %d dirs: %v", len(dirs), dirs)
	c := make([]Completer, len(dirs))
	for i := range dirs {
		c[i] = NewFileCompleter(dirs[i])
	}
	return NewMultiCompleter(c[0], c[1:]...), nil
}

// NewPathCompleter calls NewEnvCompleter with "PATH" as the
// variable name. It can be used to create completers for shells.
func NewPathCompleter() (Completer, error) {
	// Getenv returns the same value ("") if a path is not found
	// or if it has the value "". Oh well.
	return NewEnvCompleter("PATH")
}
