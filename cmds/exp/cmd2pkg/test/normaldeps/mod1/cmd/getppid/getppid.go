// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// getppid is a package that has one external dependency.
package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println(unix.Getppid())
}
