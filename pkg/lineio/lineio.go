// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lineio

import (
	"io"
	"regexp"

	"github.com/u-root/u-root/pkg/sortedmap"
)

type LineReader struct {
	src io.ReaderAt

	// offsetCache remembers the offset of various lines in src.
	// At minimum, line 1 must be prepopulated.
	offsetCache sortedmap.Map
}

// scanForLine reads from curOffset (which is on curLine), looking for line,
// returning the offset of line.  If err == io.EOF, offset is the offset of
// the last valid byte.
func (l *LineReader) scanForLine(line, curLine, curOffset int64) (offset int64, err error) {
	lastGoodOffset := int64(-1)

	for {
		buf := make([]byte, 128)

		n, err := l.src.ReadAt(buf, curOffset)
		// Keep looking as long as *something* is returned
		if n == 0 && err != nil {
			// In the event of EOF, callers want to know the last
			// byte read, to find the last byte in the last line.
			return lastGoodOffset, err
		}

		buf = buf[:n]

		for i, b := range buf {
			if b != '\n' {
				continue
			}

			offset := curOffset + int64(i) + 1

			// We haven't read this offset yet; it may not exist.
			// Read-ahead by a byte to double-check this offset
			// exists before adding a new line.
			// This is a common case for files ending in a newline.
			if offset >= curOffset+int64(len(buf)) {
				t := make([]byte, 1)
				_, err = l.src.ReadAt(t, offset)
				if err != nil {
					continue
				}
			}

			curLine++

			l.offsetCache.Insert(curLine, offset)

			if curLine == line {
				return offset, nil
			}
		}

		curOffset += int64(len(buf))
		// The last byte in the buffer must have been good if we read it.
		lastGoodOffset = curOffset - 1
	}
}

// Populate scans the file, populating the offsetCache, so that future
// lookups will be faster.
func (l *LineReader) Populate() {
	// Scan from the start of the file to the maximum possible line,
	// populating the offsetCache along the way.
	l.scanForLine(int64(0x7fffffffffffffff), 1, 0)
}

// findLine returns the offset of start of line.
func (l *LineReader) findLine(line int64) (offset int64, err error) {
	nearest, offset, err := l.offsetCache.NearestLessEqual(line)
	if err != nil {
		return 0, err
	}

	// Is this the line we want?
	if nearest == line {
		return offset, nil
	}

	return l.scanForLine(line, nearest, offset)
}

// findLineRange returns the offset of the first and last bytes in line.
func (l *LineReader) findLineRange(line int64) (start, end int64, err error) {
	start, err = l.findLine(line)
	if err != nil {
		return 0, 0, err
	}

	end, err = l.findLine(line + 1)
	// EOF means there is no next line.  End is the last byte in the file,
	// if it is positive.
	if err == io.EOF && end >= 0 {
		return start, end, nil
	} else if err != nil {
		return 0, 0, err
	}

	// The caller expects end to be the last character in the line,
	// but findLine returns the start of the next line.  Subtract
	// first character in next line and newline at end of previous line.
	end -= 2

	return start, end, nil
}

// LineExists returns true if the given line is in the file.
func (l *LineReader) LineExists(line int64) bool {
	_, err := l.findLine(line)
	return err == nil
}

// ReadLine reads up to len(p) bytes from line number line from the source.
// It returns the numbers of bytes written and any error encountered.
// If n < len(p), err is set to a non-nil value explaining why.
// See io.ReaderAt for full description of return values.
func (l *LineReader) ReadLine(p []byte, line int64) (n int, err error) {
	start, end, err := l.findLineRange(line)
	if err != nil {
		return 0, err
	}

	var shrunk bool
	// Only read one line worth of data.
	size := end - start + 1
	if size < int64(len(p)) {
		p = p[:size]
		shrunk = true
	}

	n, err = l.src.ReadAt(p, start)
	// We used less than len(p), we must return EOF.
	if err == nil && shrunk {
		err = io.EOF
	}

	return n, err
}

// SearchLine runs Regexp.FindAllIndex on the given line, providing the same
// return value.
func (l *LineReader) SearchLine(r *regexp.Regexp, line int64) ([][]int, error) {
	start, end, err := l.findLineRange(line)
	if err != nil {
		return nil, err
	}

	size := end - start + 1
	buf := make([]byte, size)

	_, err = l.src.ReadAt(buf, start)
	// TODO(prattmic): support partial reads
	if err != nil {
		return nil, err
	}

	return r.FindAllIndex(buf, -1), nil
}

func NewLineReader(src io.ReaderAt) *LineReader {
	l := &LineReader{
		src:         src,
		offsetCache: sortedmap.NewMap(),
	}

	// Line 1 starts at the beginning of the file!
	l.offsetCache.Insert(1, 0)

	return l
}
