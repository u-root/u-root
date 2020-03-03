// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import "path/filepath"

// FileCompleter is used to implement a Completer for a single
// directory in a file system.
type FileCompleter struct {
	// Root is the starting point for this Completer.
	Root string
}

// NewFileCompleter returns a FileCompleter for a single directory.
func NewFileCompleter(s string) Completer {
	return &FileCompleter{Root: s}
}

// Complete implements complete for a file starting at a directory.
func (f *FileCompleter) Complete(s string) (string, []string, error) {
	// Check for an exact match. If so, that is good enough.
	var x string
	p := filepath.Join(f.Root, s)
	Debug("FileCompleter: Check %v with %v", s, p)
	g, err := filepath.Glob(p)
	Debug("FileCompleter: %s: matches %v, err %v", s, g, err)
	if len(g) > 0 {
		x = g[0]
	}
	p = filepath.Join(f.Root, s+"*")
	Debug("FileCompleter: Check %v* with %v", s, p)
	g, err = filepath.Glob(p)
	Debug("FileCompleter: %s*: matches %v, err %v", s, g, err)
	if err != nil || len(g) == 0 {
		// one last test: directory?
		p = filepath.Join(f.Root, s, "*")
		g, err = filepath.Glob(p)
		if err != nil || len(g) == 0 {
			return x, nil, err
		}
	}
	// Here's a complication: we don't want to repeat
	// the exact match in the g array
	var ret []string
	for i := range g {
		if g[i] == x {
			continue
		}
		ret = append(ret, g[i])
	}
	Debug("FileCompleter: %s: returns %v, %v", s, g, ret)
	return x, ret, err
}
