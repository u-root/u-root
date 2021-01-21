// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ts contains a Transform to prepend a timestamp in front of each line.
package ts

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// PrependTimestamp is an io.Reader which prepends a timestamp on each line.
type PrependTimestamp struct {
	R         io.Reader
	StartTime time.Time
	Format    func(startTime time.Time) string

	// ResetTimeOnNextRead is set to force the StartTime to reset to the
	// current time when data is next read. It is useful to set this to
	// true initially for all times to be relative to the first line.
	ResetTimeOnNextRead bool

	// noPrintTS is true to indicate the timestamp should be printed
	// in front of the next byte.
	noPrintTS bool

	// buf is a Buffer to store data in case the buffer passed to
	// PrependTimestamp.Read is not large enough.
	buf []byte

	// isEOF indicates the wrapped reader R returned an io.EOF. All future
	// calls to PrependTimestamp.Read should just read from buf until it is
	// empty, then return io.EOF.
	isEOF bool
}

// New creates a PrependTimestamp with default settings.
func New(r io.Reader) *PrependTimestamp {
	return &PrependTimestamp{
		R:         r,
		StartTime: time.Now(),
		Format:    DefaultFormat,
	}
}

// Read prepends a timestamp on each line.
func (t *PrependTimestamp) Read(p []byte) (n int, err error) {
	// Empty the buffer first.
	n = copy(p, t.buf)
	t.buf = t.buf[n:]
	if len(t.buf) != 0 {
		return // Buffer not yet empty.
	}
	if t.isEOF {
		err = io.EOF
		return // Buffer empty and reached EOF.
	}

	// Read new data.
	var m int
	scratch := p[n:]
	m, err = t.R.Read(scratch)
	scratch = scratch[:m]
	if m == 0 {
		return
	}
	if t.ResetTimeOnNextRead {
		t.StartTime = time.Now()
		t.ResetTimeOnNextRead = false
	}
	if err == io.EOF {
		// Do not return EOF immediately because there may be more than
		// one calls to PrependTimestamp.Read.
		t.isEOF = true
		err = nil
	}

	// Generate the timestamp.
	lfts := []byte("\n" + t.Format(t.StartTime))
	ts := lfts[1:]

	// Insert timestamps after newlines.
	t.buf = bytes.ReplaceAll(scratch, []byte{'\n'}, lfts)

	// If the input ends in a newline, defer the insertion of the timestamp
	// until the next byte is read.
	if !t.noPrintTS {
		t.buf = append(ts, t.buf...)
		t.noPrintTS = true
	}
	if scratch[len(scratch)-1] == '\n' {
		t.buf = t.buf[:len(t.buf)-len(ts)]
		t.noPrintTS = false
	}

	// Empty new data from the buffer.
	m = copy(p[n:], t.buf)
	n += m
	t.buf = t.buf[m:]
	if len(t.buf) == 0 && t.isEOF {
		err = io.EOF // Buffer empty and reached EOF.
	}
	return
}

// DefaultFormat formats in seconds since the startTime. Ex: [12.3456s]
func DefaultFormat(startTime time.Time) string {
	return fmt.Sprintf("[%06.4fs] ", time.Since(startTime).Seconds())
}

// NewRelativeFormat returns a format function which formats in seconds since
// the previous line. Ex: [+1.0050s]
func NewRelativeFormat() func(time.Time) string {
	firstLine := true
	var lastTime time.Time
	return func(startTime time.Time) string {
		if firstLine {
			firstLine = false
			lastTime = startTime
		}
		curTime := time.Now()
		s := fmt.Sprintf("[+%06.4fs] ", curTime.Sub(lastTime).Seconds())
		lastTime = curTime
		return s
	}
}
