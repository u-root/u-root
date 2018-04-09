// +build !linux
// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func canExecute(path string) bool {
	return true
}
