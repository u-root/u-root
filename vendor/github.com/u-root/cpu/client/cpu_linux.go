// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"os"
	"syscall"

	"github.com/hugelgupf/p9/p9"
)

func osflags(fi os.FileInfo, mode p9.OpenFlags) int {
	flags := int(mode)
	if fi.IsDir() {
		flags |= syscall.O_DIRECTORY
	}
	return flags
}
