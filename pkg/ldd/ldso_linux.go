// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldd

import (
	"fmt"
	"path/filepath"
)

const ldso = "/lib*/ld-linux-*.so.*"

func LdSo() (string, error) {
	n, err := filepath.Glob(ldso)
	if err != nil {
		return "", err
	}
	if len(n) == 0 {
		return "", fmt.Errorf("No ld.so matches %v", ldso)
	}
	return n[0], nil
}
