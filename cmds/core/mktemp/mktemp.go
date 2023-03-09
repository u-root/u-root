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
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
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

func usage() {
	log.Fatalf("Usage: mktemp [options] [template]\n%v", flag.CommandLine.FlagUsages())
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

func command(stdout io.Writer, f flags, args []string) *cmd {
	return &cmd{
		stdout: stdout,
		flags:  f,
		args:   args,
	}
}

func (c *cmd) run() error {
	switch len(c.args) {
	case 1:
		c.flags.prefix = c.flags.prefix + strings.Split(c.args[0], "X")[0] + c.flags.suffix
	case 0:
	default:
		usage()
	}

	fileName, err := c.mktemp()
	if err != nil && !c.flags.q {
		return err
	}

	fmt.Fprintf(c.stdout, "%s\n", fileName)
	return nil
}

func (f *flags) register(fs *flag.FlagSet) {
	fs.BoolVarP(&f.d, "directory", "d", false, "Make a directory")
	fs.BoolVarP(&f.u, "dry-run", "u", false, "Do everything save the actual create")
	fs.BoolVarP(&f.v, "quiet", "q", false, "Quiet: show no errors")
	fs.StringVarP(&f.prefix, "prefix", "s", "", "add a prefix -- the s flag is for compatibility with GNU mktemp")
	fs.StringVarP(&f.suffix, "suffix", "", "", "add a suffix to the prefix (rather than the end of the mktemp file)")
	fs.StringVarP(&f.dir, "tmpdir", "p", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp")
}

func main() {
	flags := flags{}
	flags.register(flag.CommandLine)
	flag.Parse()
	if err := command(os.Stdout, flags, flag.Args()).run(); err != nil {
		log.Fatal(err)
	}
}
