// Copyright 2009-2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lddfiles prints the arguments and all .so dependencies of those arguments
//
// Description:
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
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/ldd"
)

func usage() {
	log.Fatalf("usage: lddfiles file [file...]")
}

func ldd(o io.Writer, s ...string) error {
	l, err := uroot.Ldd(s)
	if err != nil {
		return err
	}
	for i := range s {
		fmt.Fprintf(o, "%s\n", s[i])
	}
	for i := range l {
		fmt.Fprintf(o, "%s\n", l[i].FullName)
	}
	return nil
}

func main() {

	if err := ldd(os.Stdout, os.Args[1:]...); err != nil {
		log.Fatalf("ldd: %v", err)
	}
}
