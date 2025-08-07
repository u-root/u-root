// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gzip implements the gzip command.
package gzip

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	pkggzip "github.com/u-root/u-root/pkg/gzip"
)

// Gzip implements the gzip command.
type Gzip struct {
	core.Base
	cmdLine *flag.FlagSet
}

// New returns a new Gzip command.
func New() core.Command {
	g := &Gzip{
		cmdLine: flag.NewFlagSet("gzip", flag.ContinueOnError),
	}
	g.Init()
	return g
}

// Run executes the gzip command with the given arguments.
func (g *Gzip) Run(args ...string) error {
	return g.RunContext(context.Background(), args...)
}

// RunContext executes the gzip command with the given arguments and context.
func (g *Gzip) RunContext(ctx context.Context, args ...string) error {
	// Create a new flag set for each run to avoid flag redefinition errors
	g.cmdLine = flag.NewFlagSet("gzip", flag.ContinueOnError)
	g.cmdLine.SetOutput(g.Stderr)
	g.cmdLine.Usage = g.usage

	// Handle the case when args is empty (for tests)
	if len(args) == 0 {
		args = []string{"gzip"}
	}

	var opts pkggzip.Options
	if err := opts.ParseArgs(args, g.cmdLine); err != nil {
		if errors.Is(err, pkggzip.ErrStdoutNoForce) {
			return fmt.Errorf("gzip: %w", err)
		}
		if errors.Is(err, pkggzip.ErrHelp) {
			g.cmdLine.Usage()
			return nil
		}
		fmt.Fprintf(g.Stderr, "%s\n", err)
		g.cmdLine.Usage()
		return err
	}

	return g.run(opts, g.cmdLine.Args())
}

func (g *Gzip) usage() {
	fmt.Fprintf(g.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
	g.cmdLine.PrintDefaults()
}

func (g *Gzip) run(opts pkggzip.Options, args []string) error {
	var input []pkggzip.File
	if len(args) == 0 {
		// no args given, compress stdin to stdout
		input = append(input, pkggzip.File{Options: &opts})
	} else {
		for _, arg := range args {
			// Resolve path relative to working directory
			resolvedPath := g.ResolvePath(arg)
			input = append(input, pkggzip.File{Path: resolvedPath, Options: &opts})
		}
	}

	for i := range input {
		// We need to use a pointer to the File struct to modify it
		f := &input[i]

		// Custom implementation of CheckPath that respects working directory
		if err := g.checkPath(f); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(g.Stderr, "%s\n", err)
			}
			continue
		}

		if err := f.CheckOutputStdout(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(g.Stderr, "%s\n", err)
			}
			return err
		}

		// Custom implementation of CheckOutputPath that respects working directory
		if err := g.checkOutputPath(f); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(g.Stderr, "%s\n", err)
			}
			continue
		}

		if err := g.processFile(f); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(g.Stderr, "%s\n", err)
			}
		}

		if err := f.Cleanup(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(g.Stderr, "%s\n", err)
			}
			continue
		}
	}

	return nil
}

// checkPath is a custom implementation of CheckPath that respects working directory
func (g *Gzip) checkPath(f *pkggzip.File) error {
	if f.Options.Stdin {
		return nil
	}

	// Note: on Darwin, this permission test is not that reliable.
	_, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		return err
	} else if os.IsPermission(err) {
		return err
	}

	if !f.Options.Force {
		if f.Options.Decompress {
			if !strings.HasSuffix(f.Path, f.Options.Suffix) {
				return fmt.Errorf("%q does not have %q suffix", f.Path, f.Options.Suffix)
			}
		} else {
			if strings.HasSuffix(f.Path, f.Options.Suffix) {
				return fmt.Errorf("%q already has %q suffix", f.Path, f.Options.Suffix)
			}
		}
	}
	return nil
}

// checkOutputPath is a custom implementation of CheckOutputPath that respects working directory
func (g *Gzip) checkOutputPath(f *pkggzip.File) error {
	outputPath := resolveOutputPath(g, f)
	_, err := os.Stat(outputPath)
	if !os.IsNotExist(err) && !f.Options.Stdout && !f.Options.Test && !f.Options.Force {
		return err
	} else if os.IsPermission(err) {
		return err
	}
	return nil
}

// processFile is a modified version of the File.Process method that uses our IO streams
func (g *Gzip) processFile(f *pkggzip.File) error {
	var i *os.File
	var err error

	if f.Options.Stdin {
		i = os.Stdin
	} else {
		i, err = os.Open(f.Path)
		if err != nil {
			return err
		}
		defer i.Close()
	}

	// Use our own stdout/stderr
	var output io.Writer
	var outputCloser io.Closer
	var outputName string

	if f.Options.Test {
		output = io.Discard
		outputName = "discard"
	} else if f.Options.Stdout {
		output = g.Stdout
		outputName = "stdout"
	} else {
		// Resolve the output path relative to the working directory
		outputPath := resolveOutputPath(g, f)
		o, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		output = o
		outputCloser = o
		outputName = outputPath
	}

	if f.Options.Verbose && !f.Options.Quiet {
		fmt.Fprintf(g.Stderr, "%s to %s\n", i.Name(), outputName)
	}

	if f.Options.Decompress {
		if err := pkggzip.Decompress(i, output, f.Options.Blocksize, f.Options.Processes); err != nil {
			if outputCloser != nil {
				outputCloser.Close()
			}
			return err
		}
	} else {
		if err := pkggzip.Compress(i, output, f.Options.Level, f.Options.Blocksize, f.Options.Processes); err != nil {
			if outputCloser != nil {
				outputCloser.Close()
			}
			return err
		}
	}

	if outputCloser != nil {
		return outputCloser.Close()
	}
	return nil
}
