// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tar implements the tar command.
package tar

import (
	"archive/tar"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/tarutil"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// Tar implements the tar command.
type Tar struct {
	core.Base
	params
}

type params struct {
	file        string
	create      bool
	extract     bool
	list        bool
	noRecursion bool
	verbose     bool
}

var (
	errCreateAndExtract     = fmt.Errorf("cannot supply both -c and -x")
	errCreateAndList        = fmt.Errorf("cannot supply both -c and -t")
	errExtractAndList       = fmt.Errorf("cannot supply both -x and -t")
	errEmptyFile            = fmt.Errorf("file is required")
	errMissingMandatoryFlag = fmt.Errorf("must supply at least one of: -c, -x, -t")
	errExtractArgsLen       = fmt.Errorf("args length should be 1")
)

// New returns a new Tar command.
func New() core.Command {
	t := &Tar{}
	t.Init()
	return t
}

// Run executes the tar command with the given arguments.
func (t *Tar) Run(args ...string) error {
	return t.RunContext(context.Background(), args...)
}

// RunContext executes the tar command with the given arguments and context.
func (t *Tar) RunContext(ctx context.Context, args ...string) error {
	f := flag.NewFlagSet("tar", flag.ContinueOnError)
	f.SetOutput(t.Stderr)

	f.BoolVar(&t.create, "create", false, "create a new tar archive from the given directory")
	f.BoolVar(&t.create, "c", false, "create a new tar archive from the given directory (shorthand)")

	f.BoolVar(&t.extract, "extract", false, "extract a tar archive from the given directory")
	f.BoolVar(&t.extract, "x", false, "extract a tar archive from the given directory (shorthand)")

	f.StringVar(&t.file, "file", "", "tar file")
	f.StringVar(&t.file, "f", "", "tar file (shorthand)")

	f.BoolVar(&t.list, "list", false, "list the contents of an archive")
	f.BoolVar(&t.list, "t", false, "list the contents of an archive (shorthand)")

	f.BoolVar(&t.noRecursion, "no-recursion", false, "do not automatically recurse into directories")

	f.BoolVar(&t.verbose, "verbose", false, "print each filename")
	f.BoolVar(&t.verbose, "v", false, "print each filename (shorthand)")

	if err := f.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := t.validate(f.Args()); err != nil {
		f.Usage()
		return err
	}

	return t.execute(f.Args())
}

func (t *Tar) validate(args []string) error {
	if t.create && t.extract {
		return errCreateAndExtract
	}
	if t.create && t.list {
		return errCreateAndList
	}
	if t.extract && t.list {
		return errExtractAndList
	}
	if t.extract && len(args) != 1 {
		return errExtractArgsLen
	}
	if !t.extract && !t.create && !t.list {
		return errMissingMandatoryFlag
	}
	if t.file == "" {
		return errEmptyFile
	}
	return nil
}

func (t *Tar) execute(args []string) error {
	opts := &tarutil.Opts{
		NoRecursion: t.noRecursion,
	}
	if t.verbose {
		opts.Filters = []tarutil.Filter{tarutil.VerboseFilter}
	}

	// Resolve file path relative to working directory
	filePath := t.ResolvePath(t.file)

	switch {
	case t.create:
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		// Resolve all input paths relative to working directory
		resolvedArgs := make([]string, len(args))
		for i, arg := range args {
			resolvedArgs[i] = t.ResolvePath(arg)
		}

		if err := tarutil.CreateTar(f, resolvedArgs, opts); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	case t.extract:
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Resolve extract directory path
		extractDir := t.ResolvePath(args[0])
		if err := tarutil.ExtractDir(f, extractDir, opts); err != nil {
			return err
		}
	case t.list:
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Use our own implementation to list the archive to Stdout
		tr := tar.NewReader(f)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			fmt.Fprintln(t.Stdout, hdr.Name)
		}
	}

	return nil
}
