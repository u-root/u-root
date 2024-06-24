// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
)

// EOLReader is a reader that reads until the end of line and appends the EOL sequence.
type EOLReader struct {
	scanner *bufio.Scanner
	eol     []byte
}

func (cr *EOLReader) Read(p []byte) (n int, err error) {
	if !cr.scanner.Scan() {
		if err := cr.scanner.Err(); err != nil {
			return 0, err
		}
	}

	buf := cr.scanner.Bytes()

	// Remove the last element of buf
	if len(buf) > 0 {
		buf = buf[:len(buf)-1]
	}

	// Append the EOL sequence
	buf = append(buf, cr.eol...)

	copy(p, buf)

	n = len(buf)

	return n, nil
}

func NewEOLReader(reader io.Reader, eol []byte) *EOLReader {
	return &EOLReader{
		scanner: bufio.NewScanner(reader),
		eol:     eol,
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

// Logf logs a message if the verbose flag is set.
func Logf(nc *NetcatConfig, format string, args ...interface{}) {
	if nc.Output.Verbose {
		log.Printf(LOG_PREFIX+format, args...)
	}
}

// Logf logs a message if the verbose flag is set.
// The output is written to the provided writer instead of os.stderr
func FLogf(nc *NetcatConfig, w io.Writer, format string, args ...interface{}) {
	if nc.Output.Verbose {
		fmt.Fprintf(w, LOG_PREFIX+format, args...)
	}
}
