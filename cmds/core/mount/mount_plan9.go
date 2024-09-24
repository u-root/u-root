// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9

// Mount mounts servename on old, with an optional keypattern spec
//
// Synopsis:
//
//	mount [ option ... ] servename old [ spec ]
//
// Description:
//
//	Mount modifies the name space of the current
//	process and other processes in the same name space group
//	(see https://9p.io/magic/man2html/1/bind).
//
// Options:
//
//	–b:	Both files must be directories. Add the new directory to the beginning of the union directory represented by the old file.
//	–a:	Both files must be directories. Add the new directory to the end of the union directory represented by the old file.
//	–c:	This can be used in addition to any of the above to permit creation in a union directory.
//		When a new file is created in a union directory, it is placed in the first element of the union that has been bound or mounted with the –c flag.
//		If that directory does not have write permission, the create fails.
//	-C:	By default, file contents are always retrieved from the server.
//		With this option, the kernel may instead use a local cache to satisfy read(5) requests for files accessible through this mount point.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/namespace"
)

func main() {
	if len(os.Args) == 1 {
		n := fmt.Sprintf("/proc/%d/ns", os.Getpid())
		if b, err := os.ReadFile(n); err == nil {
			fmt.Print(string(b))
			os.Exit(0)
		}
		log.Fatalf("Could not read %s to get namespace", n)
	}
	mod, err := namespace.ParseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := mod.Modify(namespace.DefaultNamespace, &namespace.Builder{}); err != nil {
		log.Fatal(err)
	}
}
