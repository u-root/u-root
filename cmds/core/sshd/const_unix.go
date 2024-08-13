// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

var (
	shells = [...]string{"bash", "zsh", "gosh"}
	shell  = "/bin/sh"
)
