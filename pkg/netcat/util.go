// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bytes"
	"io"
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
