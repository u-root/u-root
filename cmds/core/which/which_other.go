// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !darwin && !dragonfly && !freebsd && !linux && !nacl && !netbsd && !openbsd && !solaris

package main

func canExecute(path string) bool {
	return true
}
