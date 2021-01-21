// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "syscall"

func init() {
	flagMap["dsync"] = bitClearAndSet{set: syscall.O_DSYNC}
	allowedFlags |= syscall.O_DSYNC
}
