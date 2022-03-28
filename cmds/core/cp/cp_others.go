// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

package main

import "golang.org/x/sys/unix"

const eNotDir = unix.ENOTDIR
