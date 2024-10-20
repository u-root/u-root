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
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"

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

// intermediateBuffer is a buffer that one can write to and read from.
type intermediateBuffer interface {
	io.ReaderFrom
	io.WriterTo
}

// chunkedBuffer is an intermediateBuffer with a specific size.
type chunkedBuffer struct {
	outChunk int64
	length   int64
	data     []byte
	flags    int
}

// newChunkedBuffer returns an intermediateBuffer that stores inChunkSize-sized
// chunks of data and writes them to writers in outChunkSize-sized chunks.
func newChunkedBuffer(inChunkSize int64, outChunkSize int64, flags int) intermediateBuffer {
	return &chunkedBuffer{
		outChunk: outChunkSize,
		length:   0,
		data:     make([]byte, inChunkSize),
		flags:    flags,
	}
}

// ReadFrom reads an inChunkSize-sized chunk from r into the buffer.
func (cb *chunkedBuffer) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(cb.data)
	cb.length = int64(n)

	// Convert to EOF explicitly.
	if n == 0 && err == nil {
		return 0, io.EOF
	}
	return int64(n), err
}

// WriteTo writes from the buffer to w in outChunkSize-sized chunks.
func (cb *chunkedBuffer) WriteTo(w io.Writer) (int64, error) {
	var i int64
	for i = 0; i < int64(cb.length); {
		chunk := cb.outChunk
		if i+chunk > cb.length {
			chunk = cb.length - i
		}
		block := cb.data[i : i+chunk]
		got, err := w.Write(block)

		// Ugh, Go cruft: io.Writer.Write returns (int, error).
		// io.WriterTo.WriteTo returns (int64, error). So we have to
		// cast.
		i += int64(got)
		if err != nil {
			return i, err
		}
		if int64(got) != chunk {
			return 0, io.ErrShortWrite
		}
	}
	return i, nil
}

// bufferPool is a pool of intermediateBuffers.
type bufferPool struct {
	f func() intermediateBuffer
	c chan intermediateBuffer
}

func newBufferPool(size int64, f func() intermediateBuffer) bufferPool {
	return bufferPool{
		f: f,
		c: make(chan intermediateBuffer, size),
	}
}

// Put returns a buffer to the pool for later use.
func (bp bufferPool) Put(b intermediateBuffer) {
	// Non-blocking write in case pool has filled up (too many buffers
	// returned, none being used).
	select {
	case bp.c <- b:
	default:
	}
}

// Get returns a buffer from the pool or allocates a new buffer if none is
// available.
func (bp bufferPool) Get() intermediateBuffer {
	select {
	case buf := <-bp.c:
		return buf
	default:
		return bp.f()
	}
}

func (bp bufferPool) Destroy() {
	close(bp.c)
}

func parallelChunkedCopy(r io.Reader, w io.Writer, inBufSize, outBufSize int64, bytesWritten *int64, flags int) error {
	if inBufSize == 0 {
		return fmt.Errorf("inBufSize is not allowed to be zero")
	}

	// Make the channels deep enough to hold a total of 1GiB of data.
	depth := (1024 * 1024 * 1024) / inBufSize
	// But keep it reasonable!
	if depth > 8192 {
		depth = 8192
	}

	readyBufs := make(chan intermediateBuffer, depth)
	pool := newBufferPool(depth, func() intermediateBuffer {
		return newChunkedBuffer(inBufSize, outBufSize, flags)
	})
	defer pool.Destroy()

	// Closing quit makes both goroutines below exit.
	quit := make(chan struct{})

	// errs contains the error value to be returned.
	errs := make(chan error, 1)
	defer close(errs)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Closing this unblocks the writing for-loop below.
		defer close(readyBufs)
		defer wg.Done()

		for {
			select {
			case <-quit:
				return
			default:
				buf := pool.Get()
				n, err := buf.ReadFrom(r)
				if n > 0 {
					readyBufs <- buf
				}
				if errors.Is(err, io.EOF) {
					return
				}
				if n == 0 || err != nil {
					errs <- fmt.Errorf("input error: %w", err)
					return
				}
			}
		}
	}()

	var writeErr error
	for buf := range readyBufs {
		if n, err := buf.WriteTo(w); err != nil {
			writeErr = fmt.Errorf("output error: %w", err)
			break
		} else {
			atomic.AddInt64(bytesWritten, n)
		}
		pool.Put(buf)
	}

	// This will force the goroutine to quit if an error occurred writing.
	close(quit)

	// Wait for goroutine to exit.
	wg.Wait()

	select {
	case readErr := <-errs:
		return readErr
	default:
		return writeErr
	}
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

	if max := s.limit - s.offset; int64(len(p)) > max {
		p = p[0:max]
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
	var f = flag.NewFlagSet(name, flag.ExitOnError)

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
		ibs = unit.MustNewUnit(ddUnits).MustNewValue(512, unit.None)
		obs = unit.MustNewUnit(ddUnits).MustNewValue(512, unit.None)
		bs  = unit.MustNewUnit(ddUnits).MustNewValue(512, unit.None)
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
		for _, c := range strings.Split(*conv, ",") {
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
		for _, f := range strings.Split(*oFlag, ",") {
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
	if err := parallelChunkedCopy(in, out, ibs.Value, obs.Value, &bytesWritten, flags); err != nil {
		return err
	}

	progress.End()
	return nil
}
