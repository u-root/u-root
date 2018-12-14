// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!amd64

package multiboot

import "errors"

func setupTrampoline(path string, infoAddr, entryPoint uintptr) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}
