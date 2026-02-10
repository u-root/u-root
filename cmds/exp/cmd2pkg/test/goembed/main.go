// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"fmt"
)

//go:embed foo/*.txt
var s string

func main() {
	fmt.Printf(s)
}
