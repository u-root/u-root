// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !(mips || mips64 || mips64le || mipsle || plan9 || windows)

package main

import (
	"fmt"
	"os"
	"syscall"
)

func deviceNumber(path string) (uint64, error) {
	// stat()
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	// retrieve device number
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("stat %v: error retrieving devno", path)
	}

	return st.Dev, nil
}
