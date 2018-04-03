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
func (f *FileCompleter) Complete(s string) ([]string, error) {
	// Check for an exact match. If so, that is good enough.
	p := filepath.Join(f.Root, s)
	n, _ := filepath.Glob(p)
	if len(n) == 1 {
		return n, nil
	}
	p = filepath.Join(f.Root, s+"*")
	Debug("FileCompleter: Check %v with %v", s, p)
	n, err := filepath.Glob(p)
	Debug("FileCompleter: %s: matches %v, err %v", s, n, err)
	if err != nil || len(n) == 0 {
		return n, err
	}
	Debug("FileCompleter: %s: returns %v", s, n)
	return n, err
}
