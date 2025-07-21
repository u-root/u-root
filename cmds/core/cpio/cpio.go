// Copyright 2013-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cpio operates on cpio files using a cpio package
// It only implements basic cpio options.
//
// Synopsis:
//
//	cpio
//
// Description:
//
// Options:
//
//	o: output an archive to stdout given a pattern
//	i: output files from a stdin stream
//	t: print table of contents
//	-v: debug prints
//
// Bugs: in i mode, it can't use non-seekable stdin, i.e. a pipe. Yep, this sucks.
// But if we implement seek on such things, we have to do it by reading, which
// really sucks. It's doable, we'll do it if we have to, but for now I'd like
// to avoid the complexity. cpio is a 40 year old concept. If you want something
// better, see ../archive which has a VTOC and separates data from metadata (unlike cpio).
// We could test for ESPIPE and fix it that way ... later.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

var (
	debug  = func(string, ...any) {}
	d      = flag.Bool("v", false, "Debug prints")
	format = flag.String("H", "newc", "format")

	errInvalidArgs = errors.New("usage of the command:\ncpio o < name-list [> archive]\ncpio i [< archive]\ncpio p destination-directory < name-list\nOptions: -H format (default: newc) -v Debug prints ")
)

func run(args []string, stdin *os.File, stdout io.Writer, d bool, format string) error {
	if d {
		debug = log.Printf
	}

	debug("Args %v", args)
	if len(args) < 1 {
		return errInvalidArgs
	}
	op := args[0]

	archiver, err := cpio.Format(format)
	if err != nil {
		return fmt.Errorf("format %q not supported: %w", format, err)
	}

	switch op {
	case "i":
		var inums map[uint64]string
		inums = make(map[uint64]string)
		rr, err := archiver.NewFileReader(stdin)
		if err != nil {
			return err
		}
		for {
			rec, err := rr.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading records: %w", err)
			}
			debug("record name %s ino %d\n", rec.Name, rec.Ino)

			// A file with zero size could be a hard link to another file
			// in the archive. The file with content always comes first.
			//
			// But we should ignore files with Ino of 0; that's an illegal value.
			// The current most common use of this command is with u-root
			// initramfs cpio files on Linux and Harvey.
			// (nobody else cares about cpio any more save kernels).
			// Those always have Ino of zero for reproducible builds.
			// Hence doing the Ino != 0 test first saves a bit of work.
			if rec.Ino != 0 {
				switch rec.Mode & cpio.S_IFMT {
				// In any Unix past about V1, you can't do os.Link from user mode.
				// Except via mkdir of course :-).
				case cpio.S_IFDIR:
				default:
					// FileSize of non-zero means it is the first and possibly
					// only instance of this file.
					if rec.FileSize != 0 {
						break
					}
					// If the file is not in []inums it is a true zero-length file,
					// not a hard link to a file already seen.
					// (pedantic mode: on Unix all files are hard links;
					// so what this comment really means is "file with more than one
					// hard link).
					ino, ok := inums[rec.Ino]
					if !ok {
						break
					}
					err := os.Link(ino, rec.Name)
					debug("Hard linking %s to %s", ino, rec.Name)
					if err != nil {
						return err
					}
					continue
				}
				inums[rec.Ino] = rec.Name
			}
			debug("Creating file %s", rec.Name)
			if err := cpio.CreateFile(rec); err != nil {
				log.Printf("Creating %q failed: %v", rec.Name, err)
			}
		}

	case "o":
		rw := archiver.Writer(stdout)
		cr := cpio.NewRecorder()
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			name := scanner.Text()
			rec, err := cr.GetRecord(name)
			if err != nil {
				return fmt.Errorf("getting record of %q failed: %w", name, err)
			}
			if err := rw.WriteRecord(rec); err != nil {
				return fmt.Errorf("writing record %q failed: %w", name, err)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}
		if err := cpio.WriteTrailer(rw); err != nil {
			return fmt.Errorf("error writing trailer record: %w", err)
		}

	case "t":
		rr, err := archiver.NewFileReader(stdin)
		if err != nil {
			return err
		}
		for {
			rec, err := rr.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error reading records: %w", err)
			}
			fmt.Fprintln(stdout, rec)
		}

	default:
		return errInvalidArgs
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(flag.Args(), os.Stdin, os.Stdout, *d, *format); err != nil {
		log.Fatalf("cpio: %v", err)
	}
}
