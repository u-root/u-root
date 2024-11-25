// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !arm && !386 && !mips && !mipsle

package brctl

import "golang.org/x/sys/unix"

func readTimeVal(p string, obj string) (unix.Timeval, error) {
	var tval unix.Timeval

	valRaw, err := readInt(p, obj)
	if err != nil {
		return tval, err
	}

	tvusec := 10000 * valRaw

	tval.Sec = int64(tvusec / 1000000)
	tval.Usec = int64(tvusec - (1000000 * (tvusec / 1000000)))

	return tval, nil
}
