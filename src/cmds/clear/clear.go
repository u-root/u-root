// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"fmt"
)

// using escape codes /033 => escape code
// The "\033[1;1H" part moves the cursor to position (1,1)
// "\033[2J" part clears the screen.
const magicPosixCleaner = "\033[1;1H\033[2J"

func clear() {
	fmt.Printf(magicPosixCleaner)
}

func main() {
	clear()
}
