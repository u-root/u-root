// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Archive archives files. The VTOC is at the front; we're not modeling tape drives or
// streams as in tar and cpio. This will greatly speed up listing the archive,
// modifying it, and so on. We think.
// Why a new tool?
package main

import (
	"fmt"
)

func toc(files ...string) error {
	for _, v := range files {
		_, vtoc, err := loadVTOC(v)
		if err != nil {
			fmt.Printf("%v", err)
		}
		for _, vv := range vtoc {
			fmt.Printf("%v\n", vv)
		}
	}
	return nil
}
