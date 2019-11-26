// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	yes := "y"
	if len(args) > 0 {
		yes = strings.Join(args, " ")
	}
	for {
		if _, err := fmt.Println(yes); err != nil {
			break
		}
	}
}
