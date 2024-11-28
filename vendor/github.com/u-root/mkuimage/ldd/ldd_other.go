// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux && !freebsd
// +build !linux,!freebsd

package ldd

// List returns nothing.
//
// It's not an error for a file to not be an ELF.
func List(names ...string) ([]string, error) {
	return nil, nil
}

// FList returns nothing.
//
// It's not an error for a file to not be an ELF.
func FList(names ...string) ([]string, error) {
	return nil, nil
}
