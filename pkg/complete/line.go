// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

// LineReader has three things and returns on string.
// The three things are an io.Reader, an io.Writer, and a Completer
// Bytes are read one at a time, and depending on their value,
// are written to the io.Writer.
// If the completer returns 0 or 1 answers, a 0 or 1 length string is returned.
// If there are two or more answers, they are printed and the line is
// printed out again.
// Most characters are just echoed.
// Special handling: newline or space returns.
// tab tries to complete
// backspace erases. Since everything is ansi, we assume ansi.
//
// LineReader is used to implement input for a Completer.
// It uses a Completer, io.Reader, io.Writer, and bytes.Buffer.
// bytes are read from the reader, processed, held in the
// bytes.Buffer and, as a side effect, some information is
// written to the io.Writer.
type LineReader struct {
	// Completer for this LineReader
	C Completer
	// R is used for input. Most characters are stored in the
	// Line, while some initiate special processing.
	R io.Reader
	// W is used for output, usually for showing completions.
	W io.Writer
	// Lines holds incoming data as it is read.
	Line bytes.Buffer
}

// NewLineReader returns a LineReader.
func NewLineReader(c Completer, r io.Reader, w io.Writer) *LineReader {
	return &LineReader{C: c, R: r, W: w}
}

// ReadOne tries to read one choice from l.R, printing out progress to l.W.
// In the case of \t it will try to complete given what it has.
// It will show the result of trying to complete it.
// If there is only one possible completion, it will return the result.
// In the case of \r or \n, if there is more than one choice,
// it will return the list of choices, preprended with that has been typed so far.
func (l *LineReader) ReadOne() ([]string, error) {
	Debug("LineReader: start with %v", l)
	for {
		var b [1]byte
		n, err := l.R.Read(b[:])
		if err != nil {
			if err == io.EOF {
				ln := l.Line.String()
				if ln == "" {
					return []string{}, nil
				}
				return l.C.Complete(ln)
			}
			return nil, err
		}
		Debug("LineReader.ReadOne: got %s, %v, %v", b, n, err)
		if n == 0 {
			continue
		}
		switch b[0] {
		default:
			Debug("LineReader.Just add it to line and pipe")
			l.Line.Write(b[:])
			l.W.Write(b[:])
		case '\n', '\r':
			err = EOL
			fallthrough
		case ' ':
			ln := l.Line.String()
			if ln == "" {
				return []string{}, nil
			}
			s, _ := l.C.Complete(ln)
			// the choice to use is always the first element.
			// In the case too many elements, put what they
			// typed so far as the only choice
			if len(s) > 1 {
				s = append([]string{ln}, s...)
			}
			return s, err
		case '\t':
			Debug("LineReader.Try complete with %s", l.Line.String())
			s, err := l.C.Complete(l.Line.String())
			Debug("LineReader.Complete returns %v, %v", s, err)
			if err != nil {
				return nil, err
			}
			if len(s) < 2 {
				Debug("Return is %v", s)
			}
			if _, err := l.W.Write([]byte(fmt.Sprintf("\n\r%v\n\r%v", s, l.Line.String()))); err != nil {
				log.Printf("Write %v: %v", s, err)
			}

		}
	}
}
