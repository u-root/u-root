// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux || !amd64

package trampoline

import "errors"

func Setup(path string, magic, infoAddr, entryPoint uintptr) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}
