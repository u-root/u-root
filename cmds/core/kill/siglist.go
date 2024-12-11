// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package main

import "fmt"

func siglist() (s string) {
	for i, sig := range signames {
		s = s + fmt.Sprintf("%d: %v\n", i, sig)
	}
	return
}
