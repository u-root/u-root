// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mktemp makes a temporary file (or directory)
//
// Synopsis:
//       mktemp [OPTION]... [TEMPLATE]
//
//       Create  a  temporary  file or directory, safely, and print its name.  TEMPLATE must contain at least 3 consecutive 'X's in last component.  If TEMPLATE is not specified, use tmp.XXXXXXXXXX, and --tmpdir is implied.  Files are
//       created u+rw, and directories u+rwx, minus umask restrictions.
//
//       -d, --directory
//              create a directory, not a file
//
//       -u, --dry-run
//              do not create anything; merely print a name (unsafe)
//
//       -q, --quiet
//              suppress diagnostics about file/dir-creation failure
//
//       --suffix=SUFF
//              append SUFF to TEMPLATE; SUFF must not contain a slash.  This option is implied if TEMPLATE does not end in X
//
//       -p DIR, --tmpdir[=DIR]
//              interpret TEMPLATE relative to DIR; if DIR is not specified, use $TMPDIR if set, else /tmp.  With this option, TEMPLATE must not be an absolute name; unlike with -t, TEMPLATE may contain  slashes,  but  mktemp  creates
//              only the final component
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type mktempflags struct {
	d      bool
	u      bool
	q      bool
	v      bool
	prefix string
	suffix string
	dir    string
}

var (
	flags mktempflags
)

func init() {
	flag.BoolVarP(&flags.d, "directory", "d", false, "Make a directory")
	flag.BoolVarP(&flags.u, "dry-run", "u", false, "Do everything save the actual create")
	flag.BoolVarP(&flags.v, "quiet", "q", false, "Quiet: show no errors")
	flag.StringVarP(&flags.prefix, "prefix", "s", "", "add a prefix -- the s flag is for compatibility with GNU mktemp")
	flag.StringVarP(&flags.suffix, "suffix", "", "", "add a suffix to the prefix (rather than the end of the mktemp file)")
	flag.StringVarP(&flags.dir, "tmpdir", "p", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp")
}

func usage() {
	log.Fatalf("Usage: mktemp [options] [template]\n%v", flag.CommandLine.FlagUsages())
}

func mktemp() (string, error) {
	if flags.dir == "" {
		flags.dir = os.Getenv("TMPDIR")
	}

	if flags.u {
		if !flags.q {
			log.Printf("Not doing anything but dry-run is an inherently unsafe concept")
		}
		return "", nil
	}

	if flags.d {
		d, err := ioutil.TempDir(flags.dir, flags.prefix)
		return d, err
	}
	f, err := ioutil.TempFile(flags.dir, flags.prefix)
	return f.Name(), err
}

func main() {
	flag.Parse()

	args := flag.Args()

	switch len(args) {
	case 1:
		// To make this work, we strip the trailing X's, since the Go runtime doesn't work
		// as old school mktemp(3) does. Just split on the first X.
		// If they also specified a suffix, well, add that to the prefix I guess.
		flags.prefix = flags.prefix + strings.Split(args[0], "X")[0] + flags.suffix
	case 0:
	default:
		usage()
	}

	fileName, err := mktemp()
	if err != nil && !flags.q {
		log.Fatalf("%v", err)
	}
	fmt.Println(fileName)
}
