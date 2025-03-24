// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// EOLReader is a reader that reads until the end of line and appends the EOL sequence.
type EOLReader struct {
	reader io.Reader
	eol    []byte
}

func (cr *EOLReader) Read(p []byte) (n int, err error) {
	buffer := make([]byte, len(p))
	bytesRead, readErr := cr.reader.Read(buffer)

	if bytesRead == 0 {
		return 0, readErr
	}

	buffer = buffer[:bytesRead]

	// Check and replace the EOL sequence if necessary.
	if len(buffer) > 0 {
		buffer = bytes.ReplaceAll(buffer, []byte("\n"), cr.eol)
	}

	// Copy the processed data into p and calculate the number of bytes copied.
	n = copy(p, buffer)
	return n, readErr
}

func NewEOLReader(reader io.Reader, eol []byte) *EOLReader {
	return &EOLReader{
		reader: reader,
		eol:    eol,
	}
}

// ConcurrentWriter wraps an io.Writer with a mutex to ensure safe concurrent writes.
type ConcurrentWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

// NewConcurrentWriter creates a new ConcurrentWriter.
func NewConcurrentWriter(w io.Writer) *ConcurrentWriter {
	return &ConcurrentWriter{writer: w}
}

// Write writes data to the underlying io.Writer, locking the writer during the operation.
func (sw *ConcurrentWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.writer.Write(p)
}

// StdoutWriteCloser wraps os.Stdout such that it satisfy the io.WriteCloser
// interface, *but* with a special Close. Per POSIX, fd#1 should always remain
// open:
//
// https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_01_05
// https://pubs.opengroup.org/onlinepubs/9799919799/functions/exec.html#tag_17_129_03
// https://pubs.opengroup.org/onlinepubs/9799919799/functions/close.html#tag_17_81_07
//
// Thus, we can't close os.Stdout, even temporarily.
//
// What we actually need to do is decrement the reference count on the file
// description that underlies file descriptor #1, without closing file
// descriptor #1. That's only doable with dup2(). Therefore, duplicate
// "/dev/null" over fd#1.
//
// On any GOOS that does not conform to POSIX (by certificate or by intent),
// StdoutWriteCloser.Close() is a no-op, as the well-known file handles
// shouldn't be closed in any case. (See e.g. <https://pkg.go.dev/os#Stderr>.)
type StdoutWriteCloser struct{}

// Write implements io.WriteCloser.Write.
func (swc *StdoutWriteCloser) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}
