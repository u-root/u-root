// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/ls"
)

func printFile(w io.Writer, stringer ls.Stringer, f file) {
	// Hide .files unless -a was given
	if *all || f.lsfi.Name[0] != '.' {
		// Print the file in the proper format.
		if *classify {
			f.lsfi.Name = f.lsfi.Name + indicator(f.lsfi)
		}
		fmt.Fprintln(w, stringer.FileString(f.lsfi))
	}
}
