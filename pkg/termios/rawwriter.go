// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"bytes"
	"io"
)

// RawWriter wraps an io.Writer and translates \n to \r\n.
// This is necessary in raw mode where OPOST (output processing) is disabled.
// Without this translation, newlines cause a "staircase effect" where each
// line starts one column to the right of the previous line.
type RawWriter struct {
	w io.Writer
}

// NewRawWriter creates a writer that translates LF to CRLF.
func NewRawWriter(w io.Writer) *RawWriter {
	return &RawWriter{w: w}
}

// Write implements io.Writer, translating \n to \r\n.
// It returns the original byte count for transparency.
func (rw *RawWriter) Write(p []byte) (n int, err error) {
	// Fast path: no newlines
	if !bytes.ContainsRune(p, '\n') {
		return rw.w.Write(p)
	}

	// Replace all \n with \r\n
	translated := bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	_, err = rw.w.Write(translated)

	// Return original byte count for transparency
	// This ensures callers see the same byte count they passed in
	return len(p), err
}

// Ensure RawWriter implements io.Writer
var _ io.Writer = (*RawWriter)(nil)
