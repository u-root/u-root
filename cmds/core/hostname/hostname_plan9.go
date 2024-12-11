// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9

package main

import "os"

func Sethostname(n string) error {
	return os.WriteFile("#c/sysname", []byte(n), 0o644)
}
