// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Bind binds new on old.
//
// Synopsis:
//	bind [ option ... ] new old
//
// Description:
//	Bind modifies the name space of the current
//	process and other processes in the same name space group
//	(see https://9p.io/magic/man2html/1/bind).
//
// Options:
//	–b:	Both files must be directories. Add the new directory to the beginning of the union directory represented by the old file.
//	–a:	Both files must be directories. Add the new directory to the end of the union directory represented by the old file.
//	–c:	This can be used in addition to any of the above to permit creation in a union directory.
//		When a new file is created in a union directory, it is placed in the first element of the union that has been bound or mounted with the –c flag.
//		If that directory does not have write permission, the create fails.

package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/namespace"
)

func main() {
	mod, err := namespace.ParseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := mod.Modify(namespace.DefaultNamespace, &namespace.Builder{}); err != nil {
		log.Fatal(err)
	}
}
