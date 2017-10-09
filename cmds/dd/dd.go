// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Convert and copy a file.
//
// Synopsis:
//     dd [OPTIONS...] [-inName FILE] [-outName FILE]
//
// Description:
//     dd is modeled after dd(1).
//
// Options:
//     -ibs n:   input block size (default=1)
//     -obs n:   output block size (default=1)
//     -bs n:    input and output block size (default=0)
//     -skip n:  skip n ibs-sized input blocks before reading (default=0)
//     -seek n:  seek n obs-sized output blocks before writing (default=0)
//     -conv s:  Convert the file on a specific way, like notrunc
//     -count n: copy only n ibs-sized input blocks
//     -inName:  defaults to stdin
//     -outName: defaults to stdout
//     -status:  print transfer stats to stderr, can be one of:
//         none:     do not display
//         xfer:     print on completion (default)
//         progress: print throughout transfer (GNU)
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ibs     = flag.Int64("ibs", 512, "Default input block size")
	obs     = flag.Int64("obs", 512, "Default output block size")
	bs      = flag.Int64("bs", 0, "Default input and output block size")
	skip    = flag.Int64("skip", 0, "skip N ibs-sized blocks before reading")
	seek    = flag.Int64("seek", 0, "seek N obs-sized blocks before writing")
	conv    = flag.String("conv", "none", "Convert the file on a specific way, like notrunc")
	count   = flag.Int64("count", math.MaxInt64, "copy only N input blocks")
	inName  = flag.String("if", "", "Input file")
	outName = flag.String("of", "", "Output file")
	status  = flag.String("status", "xfer", "display status of transfer (none|xfer|progress)")

	bytesWritten int64 // access atomically, must be global for correct alignedness
)

// intermediateBuffer is a buffer that one can write to and read from.
type intermediateBuffer interface {
	io.ReaderFrom
	io.WriterTo
}

// chunkedBuffer is an intermediateBuffer with a specific size.
type chunkedBuffer struct {
	outChunk  int64
	length    int64
	data      []byte
	transform func([]byte) []byte
}

// newChunkedBuffer returns an intermediateBuffer that stores inChunkSize-sized
// chunks of data and writes them to writers in outChunkSize-sized chunks.
func newChunkedBuffer(inChunkSize int64, outChunkSize int64, transform func([]byte) []byte) intermediateBuffer {
	return &chunkedBuffer{
		outChunk:  outChunkSize,
		length:    0,
		data:      make([]byte, inChunkSize),
		transform: transform,
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
		got, err := w.Write(cb.transform(cb.data[i : i+chunk]))
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

func parallelChunkedCopy(r io.Reader, w io.Writer, inBufSize, outBufSize int64, transform func([]byte) []byte) error {
	// Make the channels deep enough to hold a total of 1GiB of data.
	depth := (1024 * 1024 * 1024) / inBufSize
	// But keep it reasonable!
	if depth > 8192 {
		depth = 8192
	}

	readyBufs := make(chan intermediateBuffer, depth)
	pool := newBufferPool(depth, func() intermediateBuffer {
		return newChunkedBuffer(inBufSize, outBufSize, transform)
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
				if err == io.EOF {
					return
				}
				if n == 0 || err != nil {
					errs <- fmt.Errorf("input error: %v", err)
					return
				}
			}
		}
	}()

	var writeErr error
	for buf := range readyBufs {
		if n, err := buf.WriteTo(w); err != nil {
			writeErr = fmt.Errorf("output error: %v", err)
			break
		} else {
			atomic.AddInt64(&bytesWritten, n)
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
		if n, err := io.CopyN(ioutil.Discard, s.Reader, s.base); err != nil {
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
func inFile(name string, inputBytes int64, skip int64, count int64) (io.Reader, error) {
	maxRead := int64(math.MaxInt64)
	if count != math.MaxInt64 {
		maxRead = count * inputBytes
	}

	if name == "" {
		// os.Stdin is an io.ReaderAt, but you can't actually call
		// pread(2) on it, so use the copying section reader.
		return newStreamSectionReader(os.Stdin, inputBytes*skip, maxRead), nil
	}

	in, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("error opening input file %q: %v", name, err)
	}
	return io.NewSectionReader(in, inputBytes*skip, maxRead), nil
}

// outFile opens the output file and seeks to the right position.
func outFile(name string, outputBytes int64, seek int64) (io.Writer, error) {
	var out io.WriteSeeker
	var err error
	if name == "" {
		out = os.Stdout
	} else {
		if out, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
			return nil, fmt.Errorf("error opening output file %q: %v", name, err)
		}
	}
	if seek*outputBytes != 0 {
		if _, err := out.Seek(seek*outputBytes, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("error seeking output file: %v", err)
		}
	}
	return out, nil
}

type progressData struct {
	mode     string // one of: none, xfer, progress
	start    time.Time
	variable *int64 // must be aligned for atomic operations
	quit     chan struct{}
}

func progressBegin(mode string, variable *int64) (ProgressData *progressData) {
	p := &progressData{
		mode:     mode,
		start:    time.Now(),
		variable: variable,
	}
	if p.mode == "progress" {
		p.print()
		// Print progress in a separate goroutine.
		p.quit = make(chan struct{}, 1)
		go func() {
			ticker := time.Tick(1 * time.Second)
			for {
				select {
				case <-ticker:
					p.print()
				case <-p.quit:
					return
				}
			}
		}()
	}
	return p
}

func (p *progressData) end() {
	if p.mode == "progress" {
		// Properly synchronize goroutine.
		p.quit <- struct{}{}
		p.quit <- struct{}{}
	}
	if p.mode == "progress" || p.mode == "xfer" {
		// Print grand total.
		p.print()
		fmt.Fprint(os.Stderr, "\n")
	}
}

// With "status=progress", this is called from 3 places:
// - Once at the beginning to appear responsive
// - Every 1s afterwards
// - Once at the end so the final value is accurate
func (p *progressData) print() {
	elapse := time.Since(p.start)
	n := atomic.LoadInt64(p.variable)
	d := float64(n)
	const mib = 1024 * 1024
	const mb = 1000 * 1000
	// The ANSI escape may be undesirable to some eyes.
	if p.mode == "progress" {
		os.Stderr.Write([]byte("\033[2K\r"))
	}
	fmt.Fprintf(os.Stderr, "%d bytes (%.3f MB, %.3f MiB) copied, %.3f s, %.3f MB/s",
		n, d/mb, d/mib, elapse.Seconds(), float64(d)/elapse.Seconds()/mb)
}

func usage() {
	// If the conversions get more complex we can dump
	// the convs map. For now, it's not really worth it.
	log.Fatal(`Usage: dd [if=file] [of=file] [conv=lcase|ucase] [seek=#] [skip=#] [count=#] [bs=#] [ibs=#] [obs=#] [status=none|xfer|progress]
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
	// rather than, in essence, recreating all the apparatus of flag.xxxx
	// with the if= bits, including dup checking, conversion, etc. we just
	// convert the arguments and then run flag.Parse. Gross, but hey, it
	// works.
	os.Args = convertArgs(os.Args)
	flag.Parse()

	if len(flag.Args()) > 0 {
		usage()
	}

	convs := map[string]func([]byte) []byte{
		"none":  func(b []byte) []byte { return b },
		"ucase": bytes.ToUpper,
		"lcase": bytes.ToLower,
	}
	convert, ok := convs[*conv]
	if !ok {
		usage()
	}

	if *status != "none" && *status != "xfer" && *status != "progress" {
		usage()
	}
	progress := progressBegin(*status, &bytesWritten)

	// bs = both 'ibs' and 'obs' (IEEE Std 1003.1 - 2013)
	if *bs > 0 {
		*ibs = *bs
		*obs = *bs
	}

	in, err := inFile(*inName, *ibs, *skip, *count)
	if err != nil {
		log.Fatal(err)
	}
	out, err := outFile(*outName, *obs, *seek)
	if err != nil {
		log.Fatal(err)
	}
	if err := parallelChunkedCopy(in, out, *ibs, *obs, convert); err != nil {
		log.Fatal(err)
	}

	progress.end()
}
