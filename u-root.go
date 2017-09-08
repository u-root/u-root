// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

func main() {
	fmt.Printf(`Welcome to u-root.
To get the dynamically compiled initramfs, run scripts/ramfs.go.
To get the 'busybox' mode, cd bb, go build, and run ./bb
For more information, please see the web page.
`)
}
