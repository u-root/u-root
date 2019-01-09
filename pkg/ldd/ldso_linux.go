// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldd

import (
	"fmt"
	"path/filepath"
)

// This is a real mess. Ubuntu keeps makimg it worse.
func LdSo(bit64 bool) (string, error) {
	bits := 32
	if bit64 {
		bits = 64
	}
	choices := []string{fmt.Sprintf("/lib%d/ld-*.so.*", bits), "/lib/ld-*.so.*"}
	for _, d := range choices {
		n, err := filepath.Glob(d)
		if err != nil {
			return "", err
		}
		if len(n) > 0 {
			return n[0], nil
		}
	}
	return "", fmt.Errorf("No ld.so matches %v", choices)
}
