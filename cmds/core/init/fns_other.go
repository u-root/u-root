// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !(darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows)

package main

func quiesce() error {
	return nil
}
