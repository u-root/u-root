// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

package main

import (
	"fmt"
	"io"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ls"
)

var final = flag.BoolP("print-last", "p", false, "Print only the final path element of each file name")

func (c cmd) printFile(w io.Writer, stringer ls.Stringer, f file) {
	if f.err != nil {
		fmt.Fprintln(w, f.err)
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
		fmt.Fprintln(w, stringer.FileString(f.lsfi))
	}
}
