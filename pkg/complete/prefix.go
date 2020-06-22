// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import "strings"

// Prefix finds the longest common prefix string
// from a []string.
func Prefix(s []string) string {
	if len(s) == 0 {
		return ""
	}
	var a = s[0]
	for _, h := range s {
		if len(h) < len(a) {
			a = h
		}
	}
	var done bool
	for !done && len(a) > 0 {
		done = true
		for _, h := range s {
			if !strings.HasPrefix(h, a) {
				a = a[:len(a)-1]
				done = false
				break
			}

		}
	}
	return a
}
