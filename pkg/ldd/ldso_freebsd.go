// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldd

import (
	"fmt"
	"path/filepath"
)

// LdSo finds the loader binary.
func LdSo(bit64 bool) (string, error) {
	path := "/libexec/ld-elf32.so.*"
	if bit64 {
		path = "/libexec/ld-elf.so.*"
	}
	n, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	if len(n) > 0 {
		return n[0], nil
	}
	return "", fmt.Errorf("could not find ld.so in %v", path)
}
