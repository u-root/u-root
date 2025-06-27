// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package ls

import (
	"fmt"
	"strings"
)

func (c cmd) printFile(stringer Stringer, f file) {
	if f.err != nil {
		fmt.Fprintln(c.w, f.err)
		return
	}
	// Hide .files unless -a was given
	if c.all || !strings.HasPrefix(f.lsfi.Name, ".") {
		// Print the file in the proper format.
		if c.classify {
			f.lsfi.Name = f.lsfi.Name + indicator(f.lsfi)
		}
		fmt.Fprintln(c.w, stringer.FileString(f.lsfi))
	}
}
