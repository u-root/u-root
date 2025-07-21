// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dd converts and copies a file.
//
// Synopsis:
//
//	dd [OPTIONS...] [-inName FILE] [-outName FILE]
//
// Description:
//
//	dd is modeled after dd(1).
//
// Options:
//
//	-ibs n:   input block size (default=1)
//	-obs n:   output block size (default=1)
//	-bs n:    input and output block size (default=0)
//	-skip n:  skip n ibs-sized input blocks before reading (default=0)
//	-seek n:  seek n obs-sized output blocks before writing (default=0)
//	-conv s:  comma separated list of conversions (none|notrunc)
//	-count n: copy only n ibs-sized input blocks
//	-if:      defaults to stdin
//	-of:      defaults to stdout
//	-oflag:   comma separated list of out flags (none|sync|dsync)
//	-status:  print transfer stats to stderr, can be one of:
//	    none:     do not display
//	    xfer:     print on completion (default)
//	    progress: print throughout transfer (GNU)
//
// Notes:
//
//	Because UTF-8 clashes with block-oriented copying, `conv=lcase` and
//	`conv=ucase` will not be supported. Additionally, research showed these
//	arguments are rarely useful. Use tr instead.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/rck/unit"
	"github.com/u-root/u-root/pkg/progress"
)

type bitClearAndSet struct {
	clear int
	set   int
}

// N.B. The flags in os derive from syscall.
// They are, hence, mutually exclusive, on target
// kernels.
var convMap = map[string]bitClearAndSet{
	"notrunc": {clear: os.O_TRUNC},
}

var flagMap = map[string]bitClearAndSet{
	"sync": {set: os.O_SYNC},
}

var allowedFlags = os.O_TRUNC | os.O_SYNC

func dd(r io.Reader, w io.Writer, inBufSize, outBufSize int64, bytesWritten *int64) error {
	if inBufSize == 0 {
		return fmt.Errorf("inBufSize is not allowed to be zero")
	}
	// There is an optimization in the Go runtime for zero-copy,
	// which we can use when inBufSize == outBufSize.
	for inBufSize == outBufSize {
		amt, err := io.CopyN(w, r, inBufSize)
		*bytesWritten += int64(amt)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
	dat := &bytes.Buffer{}

	for {
		for int64(dat.Len()) > outBufSize {
			if _, err := io.CopyN(w, dat, outBufSize); err != nil {
				return err
			}
		}
		if _, err := io.CopyN(dat, r, inBufSize); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

	}
	for dat.Len() > 0 {
		if _, err := io.CopyN(w, dat, outBufSize); err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

// sectionReader implements a SectionReader on an underlying implementation of
// io.Reader (as opposed to io.SectionReader which uses io.ReaderAt).
type sectionReader struct {
	base   int64
	offset int64
	limit  int64
	io.Reader
}

// newStreamSectionReader uses an io.Reader to implement an io.Reader that
// seeks to offset and reads at most n bytes.
//
// This is useful if you want to use a NewSectionReader with stdin or other
// types of pipes (things that can't be seek'd or pread from).
func newStreamSectionReader(r io.Reader, offset int64, n int64) io.Reader {
	limit := offset + n
	if limit < 0 {
		limit = math.MaxInt64
	}
	return &sectionReader{offset, 0, limit, r}
}

// Read implements io.Reader.
func (s *sectionReader) Read(p []byte) (int, error) {
	if s.offset == 0 && s.base != 0 {
		if n, err := io.CopyN(io.Discard, s.Reader, s.base); err != nil {
			return 0, err
		} else if n != s.base {
			// Can't happen.
			return 0, fmt.Errorf("error skipping input bytes, short write")
		}
		s.offset = s.base
	}

	if s.offset >= s.limit {
		return 0, io.EOF
	}

	if mx := s.limit - s.offset; int64(len(p)) > mx {
		p = p[0:mx]
	}

	n, err := s.Reader.Read(p)
	s.offset += int64(n)

	// Convert to io.EOF explicitly.
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return n, err
}

// inFile opens the input file and seeks to the right position.
func inFile(stdin io.Reader, name string, inputBytes int64, skip int64, count int64) (io.Reader, error) {
	maxRead := int64(math.MaxInt64)
	if count != math.MaxInt64 {
		maxRead = count * inputBytes
	}

	if name == "" {
		// os.Stdin is an io.ReaderAt, but you can't actually call
		// pread(2) on it, so use the copying section reader.
		return newStreamSectionReader(stdin, inputBytes*skip, maxRead), nil
	}

	in, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %q: %w", name, err)
	}
	return io.NewSectionReader(in, inputBytes*skip, maxRead), nil
}

// outFile opens the output file and seeks to the right position.
func outFile(stdout io.WriteSeeker, name string, outputBytes int64, seek int64, flags int) (io.Writer, error) {
	var out io.WriteSeeker
	var err error
	if name == "" {
		out = stdout
	} else {
		perm := os.O_CREATE | os.O_WRONLY | (flags & allowedFlags)
		if out, err = os.OpenFile(name, perm, 0o666); err != nil {
			return nil, fmt.Errorf("error opening output file %q: %w", name, err)
		}
	}
	if seek*outputBytes != 0 {
		if _, err := out.Seek(seek*outputBytes, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("error seeking output file: %w", err)
		}
	}
	return out, nil
}

func usage() {
	log.Fatal(`Usage: dd [if=file] [of=file] [conv=none|notrunc] [seek=#] [skip=#]
			     [count=#] [bs=#] [ibs=#] [obs=#] [status=none|xfer|progress] [oflag=none|sync|dsync]
		options may also be invoked Go-style as -opt value or -opt=value
		bs, if specified, overrides ibs and obs`)
}

func convertArgs(osArgs []string) []string {
	// EVERYTHING in dd follows x=y. So blindly split and convert.
	var args []string
	for _, v := range osArgs {
		l := strings.SplitN(v, "=", 2)

		// We only fix the exact case for x=y.
		if len(l) == 2 {
			l[0] = "-" + l[0]
		}

		args = append(args, l...)
	}
	return args
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Args[0], os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(stdin io.Reader, stdout io.WriteSeeker, stderr io.Writer, name string, args []string) error {
	f := flag.NewFlagSet(name, flag.ExitOnError)

	var (
		skip    = f.Int64("skip", 0, "skip N ibs-sized blocks before reading")
		seek    = f.Int64("seek", 0, "seek N obs-sized blocks before writing")
		conv    = f.String("conv", "none", "comma separated list of conversions (none|notrunc)")
		count   = f.Int64("count", math.MaxInt64, "copy only N input blocks")
		inName  = f.String("if", "", "Input file")
		outName = f.String("of", "", "Output file")
		oFlag   = f.String("oflag", "none", "comma separated list of out flags (none|sync|dsync)")
		status  = f.String("status", "xfer", "display status of transfer (none|xfer|progress)")
	)
	ddUnits := unit.DefaultUnits
	ddUnits["c"] = 1
	ddUnits["w"] = 2
	ddUnits["b"] = 512
	delete(ddUnits, "B")

	var (
		ibs = unit.MustNewUnit(ddUnits).MustNewValue(math.MaxInt64, unit.None)
		obs = unit.MustNewUnit(ddUnits).MustNewValue(math.MaxInt64, unit.None)
		bs  = unit.MustNewUnit(ddUnits).MustNewValue(math.MaxInt64, unit.None)
	)
	f.Var(ibs, "ibs", "Default input block size")
	f.Var(obs, "obs", "Default output block size")
	f.Var(bs, "bs", "Default input and output block size")

	// rather than, in essence, recreating all the apparatus of flag.xxxx
	// with the if= bits, including dup checking, conversion, etc. we just
	// convert the arguments and then run flag.Parse. Gross, but hey, it
	// works.
	args = convertArgs(args)
	f.Parse(args)

	if len(f.Args()) > 0 {
		usage()
	}

	// Convert conv argument to bit set.
	flags := os.O_TRUNC
	if *conv != "none" {
		for c := range strings.SplitSeq(*conv, ",") {
			if v, ok := convMap[c]; ok {
				flags &= ^v.clear
				flags |= v.set
			} else {
				log.Printf("unknown argument conv=%s", c)
				usage()
			}
		}
	}

	// Convert oflag argument to bit set.
	if *oFlag != "none" {
		for f := range strings.SplitSeq(*oFlag, ",") {
			if v, ok := flagMap[f]; ok {
				flags &= ^v.clear
				flags |= v.set
			} else {
				log.Printf("unknown argument oflag=%s", f)
				usage()
			}
		}
	}

	if *status != "none" && *status != "xfer" && *status != "progress" {
		usage()
	}

	var bytesWritten int64
	progress := progress.New(stderr, *status, &bytesWritten)
	progress.Begin()

	// bs = both 'ibs' and 'obs' (IEEE Std 1003.1 - 2013)
	if bs.IsSet {
		ibs = bs
		obs = bs
	}

	in, err := inFile(stdin, *inName, ibs.Value, *skip, *count)
	if err != nil {
		return err
	}
	out, err := outFile(stdout, *outName, obs.Value, *seek, flags)
	if err != nil {
		return err
	}
	if err := dd(in, out, ibs.Value, obs.Value, &bytesWritten); err != nil {
		return err
	}

	progress.End()
	return nil
}
