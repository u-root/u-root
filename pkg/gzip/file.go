package gzip

import (
	"fmt"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/null"
)

// File is a file path to be compressed or decompressed.
type File struct {
	Path    string
	Options *Options
}

// outputPath removes the path suffix on decompress and adds it on compress.
// In the case of when options stdout or test are enabled it returns the path
// as is.
func (f *File) outputPath() string {
	if f.Options.Stdout || f.Options.Test {
		return f.Path
	} else if f.Options.Decompress {
		return f.Path[:len(f.Path)-len(f.Options.Suffix)]
	}
	return f.Path + f.Options.Suffix
}

// CheckPath validates the input file path. Checks on compression
// if the path has the correct suffix, and on decompression checks
// that it doesn't have the suffix. Allows override by force option.
func (f *File) CheckPath() error {
	_, err := os.Stat(f.Path)
	if os.IsNotExist(err) {
		return &appError{level: skipping, path: f.Path, msg: "does not exist"}
	} else if os.IsPermission(err) {
		return &appError{level: skipping, path: f.Path, msg: "permission denied"}
	}

	if !f.Options.Force {
		if f.Options.Decompress {
			if !strings.HasSuffix(f.Path, f.Options.Suffix) {
				return &appError{level: skipping, path: f.Path, msg: fmt.Sprintf("does not have %s suffix", f.Options.Suffix)}
			}
		} else {
			if strings.HasSuffix(f.Path, f.Options.Suffix) {
				return &appError{level: skipping, path: f.Path, msg: fmt.Sprintf("already has %s suffix", f.Options.Suffix)}
			}
		}
	}
	return f.checkOutPath()
}

// checkOutPath checks if output is attempting to write binary to stdout if
// stdout is a device. Also checks if output path already exists. Allow
// override via force option.
func (f *File) checkOutPath() error {
	if f.Options.Stdout {
		stat, _ := os.Stdout.Stat()
		if !f.Options.Decompress && !f.Options.Force && (stat.Mode()&os.ModeDevice) != 0 {
			return &appError{level: fatal, msg: "trying to write compressed data to a terminal/device (use -f to force)"}
		}
		return nil
	}
	_, err := os.Stat(f.outputPath())
	if !os.IsNotExist(err) && !f.Options.Stdout && !f.Options.Test && !f.Options.Force {
		return &appError{level: skipping, path: f.outputPath(), msg: "already exist"}
	} else if os.IsPermission(err) {
		return &appError{level: skipping, path: f.outputPath(), msg: "permission denied"}
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
	i, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer i.Close()

	// Use the null.WriteNameCloser interface so both *os.File and
	// null.WriteNameClose can be assigned to var o without any type casting below.
	var o null.WriteNameCloser

	if f.Options.Test {
		o = null.WriteNameClose
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
		if err := decompress(i, o, f.Options.Blocksize, f.Options.Processes); err != nil {
			if !f.Options.Stdout {
				o.Close()
			}
			return err
		}
	} else {
		if err := compress(i, o, f.Options.Level, f.Options.Blocksize, f.Options.Processes); err != nil {
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
