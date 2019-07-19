// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"strings"
)

// A StringCompleter performs completions against an array of strings.
// It can be used for, e.g., shell builtins.
type StringCompleter struct {
	// Names is the list of possible completions.
	Names []string
}

// NewStringCompleter returns a StringCompleter from the
// []string.
func NewStringCompleter(s []string) Completer {
	return &StringCompleter{Names: s}
}

// Complete returns a []string for each string of which the
// passed in string is a prefix. The error for now is always nil.
// If there is an exact match, only that match is returned,
// which is arguably wrong.
func (f *StringCompleter) Complete(s string) (string, []string, error) {
	var names []string
	for _, n := range f.Names {
		if n == s {
			return s, []string{}, nil
		}
		Debug("Check %v against %v", n, s)
		if strings.HasPrefix(n, s) {
			Debug("Add %v", n)
			names = append(names, n)
		}
	}
	return "", names, nil
}
