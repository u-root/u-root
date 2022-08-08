// Package limitio brings `io.Reader` and `io.Writer` with limit.
package limitio

import (
	"errors"
	"fmt"
	"io"
)

// ErrThresholdExceeded indicates stream size exceeds threshold
var ErrThresholdExceeded = errors.New("stream size exceeds threshold")

// AtMostFirstNBytes takes at more n bytes from s
func AtMostFirstNBytes(s []byte, n int) []byte {
	if len(s) <= n {
		return s
	}

	return s[:n]
}

// NewReader creates a Reader works like io.LimitedReader
// but may be tuned to report oversize read as
// a distinguishable error other than io.EOF.
func NewReader(r io.Reader, limit int, regardOverSizeEOF bool) *Reader {
	return &Reader{
		r:                 r,
		left:              limit,
		originalLimit:     limit,
		regardOverSizeEOF: regardOverSizeEOF,
	}
}

var _ io.Reader = (*Reader)(nil)

// Reader works like io.LimitedReader but may be tuned to report
// oversize read as a distinguishable error other than io.EOF.
//
// Use NewReader() to create Reader.
type Reader struct {
	r                 io.Reader
	left              int
	originalLimit     int
	regardOverSizeEOF bool
}

// Read implements io.Reader
func (lr *Reader) Read(p []byte) (n int, err error) {
	if lr.left <= 0 {
		if lr.regardOverSizeEOF {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("threshold is %d bytes: %w", lr.originalLimit, ErrThresholdExceeded)
	}
	if len(p) > lr.left {
		p = p[0:lr.left]
	}
	n, err = lr.r.Read(p)
	lr.left -= n
	return
}

var _ io.ReadCloser = (*ReadCloser)(nil)

// ReadCloser works like io.LimitedReader( but with a Close() method)
// but may be tuned to report oversize read as
// a distinguishable error other than io.EOF.
//
// User NewReadCloser() to create ReadCloser.
type ReadCloser struct {
	*Reader
	io.Closer
}

// NewReadCloser creates a ReadCloser that
// works like io.LimitedReader(but with a Close() method)
// and it may be tuned to report oversize read as
// a distinguishable error other than io.EOF.
func NewReadCloser(r io.ReadCloser, limit int, regardOverSizeEOF bool) *ReadCloser {
	return &ReadCloser{
		Closer: r,
		Reader: NewReader(r, limit, regardOverSizeEOF),
	}
}

var _ io.Writer = (*Writer)(nil)

// Writer wraps w with writing length limit.
//
// To create Writer, use NewWriter().
type Writer struct {
	w                    io.Writer
	written              int
	limit                int
	regardOverSizeNormal bool
}

// NewWriter create a writer that writes at most n bytes.
//
// regardOverSizeNormal controls whether Writer.Write() returns error
// when writing totally more bytes than n, or do no-op to inner w,
// pretending writing is processed normally.
func NewWriter(w io.Writer, n int, regardOverSizeNormal bool) *Writer {
	return &Writer{
		w:                    w,
		written:              0,
		limit:                n,
		regardOverSizeNormal: regardOverSizeNormal,
	}
}

// Writer implements io.Writer
func (lw *Writer) Write(p []byte) (n int, err error) {
	if lw.written >= lw.limit {
		if lw.regardOverSizeNormal {
			n = len(p)
			lw.written += n
			return
		}

		err = fmt.Errorf("threshold is %d bytes: %w", lw.limit, ErrThresholdExceeded)
		return
	}

	var (
		overSized   bool
		originalLen int
	)

	left := lw.limit - lw.written
	if originalLen = len(p); originalLen > left {
		overSized = true
		p = p[0:left]
	}
	n, err = lw.w.Write(p)
	lw.written += n
	if overSized && err == nil {
		// Write must return a non-nil error if it returns n < len(p).
		if lw.regardOverSizeNormal {
			return originalLen, nil
		}

		err = fmt.Errorf("threshold is %d bytes: %w", lw.limit, ErrThresholdExceeded)
		return
	}

	return
}

// Written returns number of bytes written
func (lw *Writer) Written() int {
	return lw.written
}

var _ io.WriteCloser = (*WriteCloser)(nil)

// WriteCloser wraps w with writing length limit.
//
// To create WriteCloser, use NewWriteCloser().
type WriteCloser struct {
	*Writer
	io.Closer
}

// NewWriteCloser create a WriteCloser that writes at most n bytes.
//
// regardOverSizeNormal controls whether Writer.Write() returns error
// when writing totally more bytes than n, or do no-op to inner w,
// pretending writing is processed normally.
func NewWriteCloser(w io.WriteCloser, n int, silentWhenOverSize bool) *WriteCloser {
	return &WriteCloser{
		Writer: NewWriter(w, n, silentWhenOverSize),
		Closer: w,
	}
}
