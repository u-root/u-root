// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb"
)

func run() {
	name := filepath.Base(os.Args[0])
	if err := bb.Run(name); err != nil {
		log.Fatalf("%s: %v", name, err)
	}
}

func main() {
	arg1 := os.Args[0]
	for s, err := os.Readlink(arg1); err == nil && filepath.Base(s) != "bb"; s, err = os.Readlink(arg1) {
		arg1 = s
	}
	os.Args[0] = arg1

	run()
}

func init() {
	m := func() {
		if len(os.Args) == 0 {
			log.Fatal("Arg len is 0. This is impossible")
		}
		if len(os.Args) == 1 {
			// This might be a symlink, and have been invoked by an sshd.
			// Let's try this: readlink until we get a terminal link.
			// If the final link is "", then forget it.
			var arg1 string
			for s, err := os.Readlink(os.Args[0]); err == nil && filepath.Base(s) != "bb"; s, err = os.Readlink(arg1) {
				arg1 = s
			}
			if arg1 == "" {
				log.Fatalf("os.Args is %v: you need to specify which command to invoke.", os.Args)
			}
			os.Args = append(os.Args, arg1)
		}
		// Use argv[1] as the name.
		os.Args = os.Args[1:]
		run()
	}
	bb.Register("bb", bb.Noop, m)
	bb.RegisterDefault(bb.Noop, m)
}
