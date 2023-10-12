// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lddfiles prints the arguments and all .so dependencies of those arguments
//
// Description:
//
//	lddfiles prints the arguments on the command line and all .so's
//	on which they depend. In some cases, those .so's are actually symlinks;
//	in that case, the symlink and its value are printed.
//	lddfiles can be used to package up a command for tranporation to
//	another machine, e.g.
//	scp `lddfiles /usr/bin/*` remotehost:/
//	will let you copy all of /usr/bin, and all needed libraries. to a remote
//	host.
//	lddfiles /usr/bin/* | cpio -H newc -o > /tmp/x.cpio
//	lets you easily prepare cpio archives, which can be included in a kernel
//	or similarly scp'ed to another machine.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/ldd"
)

func main() {
	l, err := ldd.FList(os.Args[1:]...)
	if err != nil {
		log.Fatalf("ldd: %v", err)
	}

	for _, dep := range l {
		fmt.Printf("%s\n", dep)
	}
}
