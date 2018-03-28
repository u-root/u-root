// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"fmt"
	"os"
	"strings"
)

// MultiCompleter is a Completer consisting of one or more Completers
// Why do this?
// We need it for paths, anyway, but consider a shell which
// has builtins and metacharacters such as >, &, etc.
// You can build a MultiCompleter which has a string completer
// and a set of file completers, so you don't need to special
// case anything.
type MultiCompleter struct {
	Completers []Completer
}

// NewMultiCompleter returns a MultiCompleter created from
// one or more Completers. It is perfectly legal to include a
// MultiCompleter.
func NewMultiCompleter(c Completer, cc ...Completer) Completer {
	return &MultiCompleter{append([]Completer{c}, cc...)}
}

// Complete Returns a []string consisting of the results
// of calling all the Completers.
func (m *MultiCompleter) Complete(s string) ([]string, error) {
	var files []string
	for _, c := range m.Completers {
		cc, err := c.Complete(s)
		if err != nil {
			Debug("MultiCompleter: %v: %v", c, err)
		}
		files = append(files, cc...)
	}
	return files, nil
}

// NewEnvCompleter creates a MultiCompleter consisting of one
// or more FileCompleters. It is given an environment variable,
// which it splits on :. If there are only zero entries,
// it returns an error; else it returns a MultiCompleter.
// N.B. it does *not* check for whether a directory exists
// or not; directories can come and go.
func NewEnvCompleter(s string) (Completer, error) {
	dirs := strings.Split(s, ":")
	if len(dirs) == 0 {
		return nil, fmt.Errorf("%s is empty", s)
	}
	c := make([]Completer, len(dirs))
	for i := range dirs {
		c[i] = NewFileCompleter(dirs[i])
	}
	return NewMultiCompleter(c[0], c[1:]...), nil
}

// NewPathCompleter calls NewEnvCompleter with "PATH" as the
// value. It can be used to create completers for shells.
func NewPathCompleter() (Completer, error) {
	// Getenv returns the same value ("") if a path is not found
	// or if it has the value "". Oh well.
	return NewEnvCompleter(os.Getenv("PATH"))
}
