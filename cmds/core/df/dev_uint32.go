// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (mips || mips64 || mips64le || mipsle) && !plan9 && !windows

package main

import "golang.org/x/sys/unix"

func deviceNumber(path string) (uint64, error) {
	st := &unix.Stat_t{}
	err := unix.Stat(path, st)
	if err != nil {
		return 0, err
	}
	return uint64(st.Dev), nil
}
