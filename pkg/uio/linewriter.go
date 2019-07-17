// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"io"
)

// FullLineWriter returns an io.Writer that waits for a full line of prints
// before calling w.Write on one line each.
func FullLineWriter(w io.Writer) io.WriteCloser {
	return &fullLineWriter{w: w}
}

type fullLineWriter struct {
	w      io.Writer
	buffer []byte
}

func (fsw *fullLineWriter) printBuf() {
	fsw.w.Write(fsw.buffer)
	fsw.buffer = nil
}

// Write implements io.Writer and buffers p until at least one full line is
// received.
func (fsw *fullLineWriter) Write(p []byte) (int, error) {
	i := bytes.LastIndexByte(p, '\n')
	if i == -1 {
		fsw.buffer = append(fsw.buffer, p...)
	} else {
		fsw.buffer = append(fsw.buffer, p[:i]...)
		fsw.printBuf()
		fsw.buffer = append([]byte{}, p[i:]...)
	}
	return len(p), nil
}

// Close implements io.Closer and flushes the buffer.
func (fsw *fullLineWriter) Close() error {
	fsw.printBuf()
	return nil
}
