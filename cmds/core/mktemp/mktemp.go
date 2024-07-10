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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type cmd struct {
	stdout io.Writer
	flags  flags
	args   []string
}

type flags struct {
	d      bool
	u      bool
	q      bool
	v      bool
	prefix string
	suffix string
	dir    string
}

func (c *cmd) mktemp() (string, error) {
	if c.flags.dir == "" {
		c.flags.dir = os.TempDir()
	}

	if c.flags.u {
		if !c.flags.q {
			log.Printf("Not doing anything but dry-run is an inherently unsafe concept")
		}
		return "", nil
	}

	if c.flags.d {
		d, err := os.MkdirTemp(c.flags.dir, c.flags.prefix)
		return d, err
	}
	f, err := os.CreateTemp(c.flags.dir, c.flags.prefix)
	return f.Name(), err
}

func (c *cmd) run() error {
	fileName, err := c.mktemp()
	if err != nil && !c.flags.q {
		return err
	}

	fmt.Fprintf(c.stdout, "%s\n", fileName)
	return nil
}

func (f *flags) register(fs *flag.FlagSet) {
	fs.BoolVar(&f.d, "directory", false, "Make a directory")
	fs.BoolVar(&f.d, "d", false, "Make a directory (shorthand)")

	fs.BoolVar(&f.u, "dry-run", false, "Do everything save the actual create")
	fs.BoolVar(&f.u, "u", false, "Do everything save the actual create (shorthand)")

	fs.BoolVar(&f.v, "quiet", false, "Quiet: show no errors")
	fs.BoolVar(&f.v, "q", false, "Quiet: show no errors (shorthand)")

	fs.StringVar(&f.prefix, "prefix", "", "add a prefix")
	fs.StringVar(&f.prefix, "s", "", "add a prefix (shorthand, 's' is for compatibility with GNU mktemp")

	fs.StringVar(&f.suffix, "suffix", "", "add a suffix to the prefix (rather than the end of the mktemp file)")

	fs.StringVar(&f.dir, "tmpdir", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp")
	fs.StringVar(&f.dir, "p", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp (shorthand)")

}

func command(stdout io.Writer, args []string) *cmd {
	var c cmd
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.BoolVar(&c.flags.d, "directory", false, "Make a directory")
	fs.BoolVar(&c.flags.d, "d", false, "Make a directory (shorthand)")

	fs.BoolVar(&c.flags.u, "dry-run", false, "Do everything save the actual create")
	fs.BoolVar(&c.flags.u, "u", false, "Do everything save the actual create (shorthand)")

	fs.BoolVar(&c.flags.v, "quiet", false, "Quiet: show no errors")
	fs.BoolVar(&c.flags.v, "q", false, "Quiet: show no errors (shorthand)")

	fs.StringVar(&c.flags.prefix, "prefix", "", "add a prefix")
	fs.StringVar(&c.flags.prefix, "s", "", "add a prefix (shorthand, 's' is for compatibility with GNU mktemp")

	fs.StringVar(&c.flags.suffix, "suffix", "", "add a suffix to the prefix (rather than the end of the mktemp file)")

	fs.StringVar(&c.flags.dir, "tmpdir", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp")
	fs.StringVar(&c.flags.dir, "p", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mktemp [options] [template]\n")
		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	c.stdout = stdout
	c.args = fs.Args()

	switch len(c.args) {
	case 1:
		c.flags.prefix = c.flags.prefix + strings.Split(c.args[0], "X")[0] + c.flags.suffix
	case 0:
	default:
		fs.Usage()
		os.Exit(1)
	}

	return &c
}

func main() {
	if err := command(os.Stdout, os.Args).run(); err != nil {
		log.Fatal(err)
	}
}
