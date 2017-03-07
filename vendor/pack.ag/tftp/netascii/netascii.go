// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

/*
Package netascii implements reading and writing of netascii, as defined in RFC 764.

Netascii encodes LF to CRLF and CR to CRNUL.
CRLF is decoded to the platform's representation of a new line.
*/
package netascii // import "pack.ag/tftp/netascii"

import (
	"bufio"
	"io"
	"runtime"
)

const (
	cr  = '\r'
	lf  = '\n'
	nul = 0
)

// Reader is an io.Reader used to retrieve data in the
// local system's format from a netascii encoded source.
type Reader struct {
	r *bufio.Reader
}

// NewReader returns a Reader wrapping r.
func NewReader(reader io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(reader)}
}

// Read reads and decodes netascii from r.
func (d *Reader) Read(p []byte) (int, error) {
	// Encodes to netascii, processing 1 byte at a time
	bufLen := len(p)
	written := 0

	for written < bufLen {
		current, err := d.r.ReadByte()
		if err != nil {
			return written, err
		}

		if current == cr {
			b, err := d.r.ReadByte()
			if err != nil {
				return written, err
			}
			if runtime.GOOS != "windows" && b == lf {
				// CRLF becomes LF
				current = lf
			} else if b == nul {
				// CRNUL becomes CR
			} else {
				// Next byte isn't LF or NUL
				d.r.UnreadByte()
			}
		}

		p[written] = current
		written++
	}
	return written, nil
}

// Writer is an io.Writer. Writes to Writer are encoded into netascii and written to w.
type Writer struct {
	w    *bufio.Writer
	last byte
}

// NewWriter returns a Writer wrapping the w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Write encodes p as netascii and writes it to w. Writer must be flushed to
// guarantee that all data has been written to w.
func (e *Writer) Write(p []byte) (int, error) {
	written := 0
	var err error // Declare here and break to avoid duplication of written > len(p) logic

	for _, current := range p {
		if current == lf && e.last != cr {
			// LF becomes CRLF
			err = e.w.WriteByte(cr)
			if err != nil {
				break
			}
			e.last = cr
			written++
		} else if e.last == cr && current != lf && current != nul {
			// CR becomes CRNUL
			err = e.w.WriteByte(nul)
			if err != nil {
				break
			}
			e.last = nul
			written++
		}

		err = e.w.WriteByte(current)
		if err != nil {
			break
		}
		e.last = current
		written++
	}

	// We may have written more than p, which is an error,
	// return len(p)
	if written > len(p) {
		return len(p), err
	}
	return written, err
}

// Flush flushes any pending data to w.
func (e *Writer) Flush() error {
	return e.w.Flush()
}
