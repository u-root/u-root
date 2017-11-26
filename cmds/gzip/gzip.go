// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/klauspost/pgzip"
	"github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/null"
)

const version string = "0.0.1"

type options struct {
	blocksize  int
	level      int
	processes  int
	decompress bool
	force      bool
	help       bool
	keep       bool
	license    bool
	quiet      bool
	stdin      bool
	stdout     bool
	test       bool
	verbose    bool
	version    bool
	suffix     string
}

func (o *options) parseArgs() error {
	var levels [10]bool

	pflag.IntVarP(&o.blocksize, "blocksize", "b", 128, "Set compression block size in KiB")
	pflag.BoolVarP(&o.decompress, "decompress", "d", false, "Decompress the compressed input")
	pflag.BoolVarP(&o.force, "force", "f", false, "Force overwrite of output file and compress links")
	pflag.BoolVarP(&o.help, "help", "h", false, "Display a help screen and quit")
	pflag.BoolVarP(&o.keep, "keep", "k", false, "Do not delete original file after processing")
	pflag.BoolVarP(&o.license, "license", "L", false, "Display license")
	// TODO: implement list option here
	pflag.IntVarP(&o.processes, "processes", "p", runtime.NumCPU(), "Allow up to n compression threads")
	pflag.BoolVarP(&o.quiet, "quiet", "q", false, "Print no messages, even on error")
	// TODO: implement recursive option here
	pflag.BoolVarP(&o.stdout, "stdout", "c", false, "Write all processed output to stdout (won't delete)")
	pflag.StringVarP(&o.suffix, "suffix", "S", ".gz", "Specify suffix for compression")
	pflag.BoolVarP(&o.test, "test", "t", false, "Test the integrity of the compressed input")
	pflag.BoolVarP(&o.verbose, "verbose", "v", false, "Produce more verbose output")
	pflag.BoolVarP(&o.version, "version", "V", false, "Print version")
	pflag.BoolVarP(&levels[1], "fast", "1", false, "Compression level 1")
	pflag.BoolVarP(&levels[2], "two", "2", false, "Compression level 2")
	pflag.BoolVarP(&levels[3], "three", "3", false, "Compression level 3")
	pflag.BoolVarP(&levels[4], "four", "4", false, "Compression level 4")
	pflag.BoolVarP(&levels[5], "five", "5", false, "Compression level 5")
	pflag.BoolVarP(&levels[6], "six", "6", false, "Compression level 6")
	pflag.BoolVarP(&levels[7], "seven", "7", false, "Compression level 7")
	pflag.BoolVarP(&levels[8], "eight", "8", false, "Compression level 8")
	pflag.BoolVarP(&levels[9], "best", "9", false, "Compression level 9")

	// Hide the compression levels 2 - 8 from usage.
	_ = pflag.CommandLine.MarkHidden("two")
	_ = pflag.CommandLine.MarkHidden("three")
	_ = pflag.CommandLine.MarkHidden("four")
	_ = pflag.CommandLine.MarkHidden("five")
	_ = pflag.CommandLine.MarkHidden("six")
	_ = pflag.CommandLine.MarkHidden("seven")
	_ = pflag.CommandLine.MarkHidden("eight")

	pflag.Parse()

	var err error
	o.level, err = parseLevels(levels)
	if err != nil {
		return err
	}

	return o.validate()
}

// Validate checks options and handles help, version, and license
// output to the user. If forces decompression to be enabled when
// test mode is enabled. It further modifies options if the running
// binary is named gunzip or gzcat to allow for expected behavor for
// those binaries. Finally it checks if there is piped stdin data.
func (o *options) validate() error {

	if o.help {
		pflag.Usage()
		os.Exit(0)
	}

	if o.version {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if o.license {
		fmt.Printf("%s\n%s\n%s\n",
			version,
			"Copyright (c) 2012-2017 The u-root Authors. All rights reserved.",
			"Subject to the terms of the BSD 3-Clause license.")
		os.Exit(0)
	}

	if o.test {
		o.decompress = true
	}

	// Support gunzip and gzcat symlinks
	if filepath.Base(os.Args[0]) == "gunzip" {
		o.decompress = true
	} else if filepath.Base(os.Args[0]) == "gzcat" {
		o.decompress = true
		o.stdout = true
	}

	// Stat os.Stdin and ignore errors. stat will be nil FileInfo if there is an
	// error.
	stat, _ := os.Stdin.Stat()

	// No files passed and arguments and Stdin piped data found.
	// Stdin piped data is ignored if arguments are found.
	if len(pflag.Args()) == 0 && (stat.Mode()&os.ModeNamedPipe) != 0 {
		o.stdin = true
		// Enable force to ignore suffix checks
		o.force = true
		// Since there's no filename to derive the output path from, only support
		// outputting to stdout when data is piped from stdin
		o.stdout = true
	} else if len(pflag.Args()) == 0 {
		// No stdin piped data found and no files passed as arguments
		pflag.Usage()
		os.Exit(0)
	}

	return nil
}

// parseLevels loops through a [10]bool and returns the index of the element
// thats true. If more than one element is true return an error. If no
// element is true, return the constant pgzip.DefaultCompression (-1).
func parseLevels(levels [10]bool) (int, error) {
	var level int

	for i, l := range levels {
		if l && level != 0 {
			return 0, errors.New("Multiple compression levels specified")
		} else if l {
			level = i
		}
	}

	if level == 0 {
		level = pgzip.DefaultCompression
	}
	return level, nil
}

// file is a file path to be compressed or decompressed.
type file struct {
	path    string
	options *options
}

// outputPath removes the path suffix on decompress and adds it on compress.
// In the case of when options stdout or test are enabled it returns the path
// as is.
func (f *file) outputPath() string {
	if f.options.stdout || f.options.test {
		return f.path
	} else if f.options.decompress {
		return f.path[:len(f.path)-len(f.options.suffix)]
	}
	return f.path + f.options.suffix
}

// checkPath validates the input file path. Checks on compression
// if the path has the correct suffix, and on decompression checks
// that it doesn't have the suffix. Allows override by force option.
func (f *file) checkPath() error {
	_, err := os.Stat(f.path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", f.path)
	} else if os.IsPermission(err) {
		return fmt.Errorf("%s permission denied", f.path)
	}

	if !f.options.force {
		if f.options.decompress {
			if !strings.HasSuffix(f.path, f.options.suffix) {
				return fmt.Errorf("%s does not have %s suffix", f.path, f.options.suffix)
			}
		} else {
			if strings.HasSuffix(f.path, f.options.suffix) {
				return fmt.Errorf("%s already has %s suffix", f.path, f.options.suffix)
			}
		}
	}
	return f.checkOutPath()
}

// checkOutPath checks if output is attempting to write binary to stdout if
// stdout is a device. Also checks if output path already exists. Allow
// override via force option.
func (f *file) checkOutPath() error {
	if f.options.stdout {
		stat, _ := os.Stdout.Stat()
		if !f.options.decompress && !f.options.force && (stat.Mode()&os.ModeDevice) != 0 {
			return errors.New("trying to write compressed data to a terminal/device (use -f to force)")
		}
		return nil
	}
	_, err := os.Stat(f.outputPath())
	if !os.IsNotExist(err) && !f.options.stdout && !f.options.test && !f.options.force {
		return fmt.Errorf("%s already exist", f.outputPath())
	} else if os.IsPermission(err) {
		return fmt.Errorf("%s permission denied", f.outputPath())
	}
	return nil
}

// Cleanup removes input file. Overrided with keep option. Skipped if
// stdout or test option is true.
func (f *file) cleanup() error {
	if !f.options.keep && !f.options.stdout && !f.options.test {
		return os.Remove(f.path)
	}
	return nil
}

// Process either compresses or decompressed the input file based on
// the associated file.options.
func (f *file) process() error {
	i, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer i.Close()

	// Use the null.WriteNameCloser interface so both *os.File and
	// null.WriteNameClose can be assigned to var o without any type casting below.
	var o null.WriteNameCloser

	if f.options.test {
		o = null.WriteNameClose
	} else if f.options.stdout {
		o = os.Stdout
	} else {
		if o, err = os.Create(f.outputPath()); err != nil {
			return err
		}
	}

	if f.options.verbose && !f.options.quiet {
		fmt.Fprintf(os.Stderr, "%s to %s\n", i.Name(), o.Name())
	}

	if f.options.decompress {
		if err := decompressFile(i, o, f.options.blocksize, f.options.processes); err != nil {
			if !f.options.stdout {
				o.Close()
			}
			return err
		}
	} else {
		if err := compressFile(i, o, f.options.level, f.options.blocksize, f.options.processes); err != nil {
			if !f.options.stdout {
				o.Close()
			}
			return err
		}
	}

	if f.options.stdout {
		return nil
	}

	return o.Close()
}

func compressFile(r io.Reader, w io.Writer, level int, blocksize int, processes int) error {
	zw, err := pgzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}

	if err := zw.SetConcurrency(blocksize*1024, processes); err != nil {
		zw.Close()
		return err
	}

	if _, err := io.Copy(zw, r); err != nil {
		zw.Close()
		return err
	}

	return zw.Close()
}

func decompressFile(r io.Reader, w io.Writer, blocksize int, processes int) error {
	zr, err := pgzip.NewReaderN(r, blocksize*1024, processes)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, zr); err != nil {
		zr.Close()
		return err
	}

	return zr.Close()
}

func main() {

	var opts options

	if err := opts.parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Argument error: %s\n", err)
		os.Exit(1)
	}

	var files []string

	if opts.stdin {
		files = []string{"/dev/stdin"}
	} else {
		files = pflag.Args()
	}

	for _, path := range files {

		f := file{path: path, options: &opts}
		if err := f.checkPath(); err != nil {
			if !opts.quiet {
				fmt.Fprintf(os.Stderr, "skipping, %s\n", err)
			}
			continue
		}

		if err := f.process(); err != nil {
			if !opts.quiet {
				fmt.Fprintf(os.Stderr, "error, %s %s\n", f.path, err)
			}
			os.Exit(1)
		}

		if err := f.cleanup(); err != nil {
			if !opts.quiet {
				fmt.Fprintf(os.Stderr, "warning, %s %s\n", f.path, err)
			}
			continue
		}
	}
}
