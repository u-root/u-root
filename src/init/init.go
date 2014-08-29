// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// assumptions
// we've been booted into a ramfs with all this stuff unpacked and ready.
// we don't need a loop device mount because it's all there.
// So we run /go/bin/go build installcommand
// and then exec /buildbin/sh

package main

import (
	"fmt"
	"os"
)

var urpath = "/go/bin:/buildbin:/bin:/usr/local/bin:"

func main() {
o		run := exec.Command(argv[0], argv[1:]...)
		run.Env = e
		out, err := run.CombinedOutput()
		if err != nil {
			fmt.Printf("%v: Path %v\n", err, os.Getenv("PATH"))
		}
		fmt.Printf("%s", out)
		fmt.Printf("%% ")
}
