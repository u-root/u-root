// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"fmt"
	"strings"
)

// RemoveFilter filters out variable for a given space-separated kernel commandline
func RemoveFilter(input string, variables []string) string {
	var newCl []string

	// kernel variables must allow '-' and '_' to be equivalent in variable
	// names. We will replace dashes with underscores for processing as
	// `doParse` is doing.
	for i, v := range variables {
		variables[i] = strings.Replace(v, "-", "_", -1)
	}

	doParse(input, func(flag, key, canonicalKey, value, trimmedValue string) {
		skip := false
		for _, v := range variables {
			if canonicalKey == v {
				skip = true
				break
			}
		}
		if skip {
			return
		}
		newCl = append(newCl, flag)
	})
	return strings.Join(newCl, " ")
}

// UpdateFilter get the kernel command line parameters and filter it:
// it removes parameters listed in 'remove' and append extra parameters from
// the 'append' and 'reuse' flags
func UpdateFilter(cl, append string, remove, reuse []string) string {
	acl := ""
	if len(append) > 0 {
		acl = " " + append
	}
	for _, f := range reuse {
		value, present := Flag(f)
		if present {
			acl = fmt.Sprintf("%s %s=%s", acl, f, value)
		}
	}

	return RemoveFilter(cl, remove) + acl
}
