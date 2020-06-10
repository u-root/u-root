// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"errors"
	"io"
)

// SectionWriter turns an io.WriterAt into an io.Writer.
//
// SectionWriter is not goroutine-safe.
type SectionWriter struct {
	w     io.WriterAt
	base  int64
	off   int64
	limit int64
}

// NewSectionWriter allows writes to w from offset off to off+n.
func NewSectionWriter(w io.WriterAt, off, n int64) *SectionWriter {
	return &SectionWriter{
		w:     w,
		base:  off,
		off:   off,
		limit: off + n,
	}
}

func (s *SectionWriter) Write(p []byte) (int, error) {
	if s.off >= s.limit {
		return 0, io.EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
	}
	n, err := s.w.WriteAt(p, s.off)
	s.off += int64(n)
	return n, err
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (s *SectionWriter) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case io.SeekStart:
		offset += s.base
	case io.SeekCurrent:
		offset += s.off
	case io.SeekEnd:
		offset += s.limit
	}
	if offset < s.base {
		return 0, errOffset
	}
	s.off = offset
	return offset - s.base, nil
}
