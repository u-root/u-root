// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mktemp makes a temporary file (or directory)
//
// Synopsis:
//
//	mktemp [OPTION]... [TEMPLATE]
//
//	Create  a  temporary  file or directory, safely, and print its name.  TEMPLATE must contain at least 3 consecutive 'X's in last component.  If TEMPLATE is not specified, use tmp.XXXXXXXXXX, and --tmpdir is implied.  Files are
//	created u+rw, and directories u+rwx, minus umask restrictions.
//
//	-d, --directory
//	       create a directory, not a file
//
//	-u, --dry-run
//	       do not create anything; merely print a name (unsafe)
//
//	-q, --quiet
//	       suppress diagnostics about file/dir-creation failure
//
//	--suffix=SUFF
//	       append SUFF to TEMPLATE; SUFF must not contain a slash.  This option is implied if TEMPLATE does not end in X
//
//	-p DIR, --tmpdir[=DIR]
//	       interpret TEMPLATE relative to DIR; if DIR is not specified, use $TMPDIR if set, else /tmp.  With this option, TEMPLATE must not be an absolute name; unlike with -t, TEMPLATE may contain  slashes,  but  mktemp  creates
//	       only the final component
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/mktemp"
)

func main() {
	cmd := mktemp.New()
	if err := cmd.Run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
