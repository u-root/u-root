// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
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

type ConcurrentWriteCloser struct {
	ConcurrentWriter
}

func NewConcurrentWriteCloser(wc io.WriteCloser) *ConcurrentWriteCloser {
	return &ConcurrentWriteCloser{ConcurrentWriter{writer: wc}}
}

func (swc *ConcurrentWriteCloser) Close() error {
	swc.ConcurrentWriter.mu.Lock()
	defer swc.ConcurrentWriter.mu.Unlock()
	return swc.ConcurrentWriter.writer.(io.WriteCloser).Close()
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

// TeeWriteCloser is a simplistic crossing between io.MultiWriter and
// io.WriteCloser.
type TeeWriteCloser struct {
	wc1 io.WriteCloser
	wc2 io.WriteCloser
}

func NewTeeWriteCloser(wc1 io.WriteCloser, wc2 io.WriteCloser) *TeeWriteCloser {
	return &TeeWriteCloser{wc1: wc1, wc2: wc2}
}

// Write doesn't continue past the first write error, similarly to
// io.MultiWriter's Write.
func (twc *TeeWriteCloser) Write(b []byte) (n int, err error) {
	if n, err = twc.wc1.Write(b); err != nil {
		return
	}
	return twc.wc2.Write(b)
}

// Close closes both io.WriteCloser objects in all cases, and returns the first
// error (if any).
func (twc *TeeWriteCloser) Close() error {
	err1 := twc.wc1.Close()
	err2 := twc.wc2.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// CloseWrite is a convenience wrapper that calls CloseWrite on an io.Writer
// that is a TLS connection, a TCP connection, or a UNIX domain connection.
func CloseWrite(writerToClose io.Writer) error {
	writer := writerToClose

	if tlsConn, ok := writer.(*tls.Conn); ok {
		// If handshake has completed, shut down the writing side of
		// the TLS session.
		if tlsConn.ConnectionState().HandshakeComplete {
			if err := tlsConn.CloseWrite(); err != nil {
				return err
			}
		}
		// Fall through to shutting down the underlying transport.
		writer = tlsConn.NetConn()
	}

	if tcpConn, ok := writer.(*net.TCPConn); ok {
		return tcpConn.CloseWrite()
	}

	if unixConn, ok := writer.(*net.UnixConn); ok {
		return unixConn.CloseWrite()
	}

	return nil
}
