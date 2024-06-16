// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/u-root/u-root/pkg/ls"
)

var final = flag.Bool("p", false, "Print only the final path element of each file name")

func (c cmd) printFile(stringer ls.Stringer, f file) {
	if f.err != nil {
		fmt.Fprintln(c.w, f.err)
		return
	}
	// Hide .files unless -a was given
	if c.all || !strings.HasPrefix(f.lsfi.Name, ".") {
		// Unless they said -p, we always print the full path
		if !*final {
			f.lsfi.Name = f.path
		}
		if c.classify {
			f.lsfi.Name = f.lsfi.Name + indicator(f.lsfi)
		}
		fmt.Fprintln(c.w, stringer.FileString(f.lsfi))
	}
}
