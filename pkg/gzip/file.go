// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/u-root/uio/uio"
)

// File is a file path to be compressed or decompressed.
type File struct {
	Options *Options
	Path    string
}

// Compression and Decompression functions used by the File.Process method.
// Bare metal support can be enabled, for example, witht the `tinygo` build tag.
// Setting these build tags will result in the use of pure go libraries for compression/decompression.
func Compress(r io.Reader, w io.Writer, level int, blocksize int, processes int) error {
	return compress(r, w, level, blocksize, processes)
}

// Decompress takes gzip compressed input from io.Reader and expands it using
// pgzip or gzip, depending on the build tags.
// When the `tinygo` build tag is set, the `compress/gzip` package is used and the `processes`
// argument is ignored.
func Decompress(r io.Reader, w io.Writer, blocksize int, processes int) error {
	return decompress(r, w, blocksize, processes)
}

// outputPath removes the path suffix on decompress and adds it on compress.
// In the case of when options stdout or test are enabled it returns the path
// as is.
func (f *File) outputPath() string {
	if f.Options.Stdout || f.Options.Test {
		return f.Path
	} else if f.Options.Decompress {
		return strings.TrimSuffix(f.Path, f.Options.Suffix)
	}
	return f.Path + f.Options.Suffix
}

// CheckPath validates the input file path. Checks on compression
// if the path has the correct suffix, and on decompression checks
// that it doesn't have the suffix. Allows override by force option.
// Skip if the input is a Stdin.
func (f *File) CheckPath() error {
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

// CheckOutputPath checks if output is attempting to write binary to stdout if
// stdout is a device. Also checks if output path already exists. Allow
// override via force option.
func (f *File) CheckOutputPath() error {
	_, err := os.Stat(f.outputPath())
	if !os.IsNotExist(err) && !f.Options.Stdout && !f.Options.Test && !f.Options.Force {
		return err
	} else if os.IsPermission(err) {
		return err
	}
	return nil
}

// CheckOutputStdout checks if output is attempting to write binary to stdout
// if stdout is a device.
func (f *File) CheckOutputStdout() error {
	if f.Options.Stdout {
		stat, _ := os.Stdout.Stat()
		if !f.Options.Decompress && !f.Options.Force && (stat.Mode()&os.ModeDevice) != 0 {
			return fmt.Errorf("can not write compressed data to a terminal/device (use -f to force)")
		}
	}
	return nil
}

// Cleanup removes input file. Overrided with keep option. Skipped if
// stdout or test option is true.
func (f *File) Cleanup() error {
	if !f.Options.Keep && !f.Options.Stdout && !f.Options.Test {
		return os.Remove(f.Path)
	}
	return nil
}

// Process either compresses or decompressed the input file based on
// the associated file.options.
func (f *File) Process() error {
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

	// Use the uio.WriteNameCloser interface so both *os.File and
	// uio.WriteNameClose can be assigned to var o without any type casting below.
	var o uio.WriteNameCloser

	if f.Options.Test {
		o = uio.Discard
	} else if f.Options.Stdout {
		o = os.Stdout
	} else {
		if o, err = os.Create(f.outputPath()); err != nil {
			return err
		}
	}

	if f.Options.Verbose && !f.Options.Quiet {
		fmt.Fprintf(os.Stderr, "%s to %s\n", i.Name(), o.Name())
	}

	if f.Options.Decompress {
		if err := Decompress(i, o, f.Options.Blocksize, f.Options.Processes); err != nil {
			if !f.Options.Stdout {
				o.Close()
			}
			return err
		}
	} else {
		if err := Compress(i, o, f.Options.Level, f.Options.Blocksize, f.Options.Processes); err != nil {
			if !f.Options.Stdout {
				o.Close()
			}
			return err
		}
	}

	if f.Options.Stdout {
		return nil
	}
	return o.Close()
}
