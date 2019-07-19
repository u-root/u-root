// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"io"
	"path/filepath"
	"strings"
)

const (
	controlD  = 4
	backSpace = 8
	del       = 127
)

// LineReader has three things and returns one string.
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
	Line string
	// Exact is the exact match.
	// It can be "" if there is not one.
	Exact string
	// Candidates are any completion candidates.
	// The UI can decide how to handle them.
	Candidates []string
	EOF        bool
	Fields     int
}

// NewLineReader returns a LineReader.
func NewLineReader(c Completer, r io.Reader, w io.Writer) *LineReader {
	return &LineReader{C: c, R: r, W: w}
}

// ReadChar reads one character and processes it. It is inflexible by design.
func (l *LineReader) ReadChar(b byte) (err error) {
	defer func() {
		l.Fields = len(strings.Fields(l.Line))
	}()
	Debug("LineReader: start with %v", l)
	l.Exact, l.Candidates = "", []string{}
	switch b {
	default:
		Debug("LineReader.Just add it to line and pipe")
		l.Line += string(b)
	case controlD:
		l.EOF = true
		return io.EOF
	case backSpace, del:
		s := l.Line
		if len(s) > 0 {
			s = s[:len(s)-1]
			l.Line = s
		}
	case '\n', '\r':
		return ErrEOL
	case ' ':
		l.Line += string(b)
		return nil
	case '\t':
		ll := len(l.Line)
		s := l.Line
		bl := strings.LastIndexAny(s, " ")
		flds := strings.Fields(s)
		Debug("fields of %q is %v", s, flds)
		var cc string
		// The rules are complex.
		// It might be zero length.
		// It might end in whitespace
		// It might be one or more things
		switch {
		case ll == 0 || bl == ll-1:
		case len(flds) == 0:
			cc = flds[0]
		default:
			cc = flds[len(flds)-1]
		}

		Debug("ReadChar.Try complete with %s", cc)
		x, cmpl, err := l.C.Complete(cc)
		Debug("ReadChar.Complete returns %q, %v, %v", x, cmpl, err)
		if err != nil {
			return err
		}
		if len(cmpl) == 1 && x == "" {
			Debug("Readchar: only one candidate, so use it")
			x, cmpl = cmpl[0], []string{}
		}
		l.Exact = x
		l.Candidates = cmpl
		if x != "" && filepath.Clean(cc) != x {
			// Paste the completion over where we found the candidate.
			Debug("ReadChar: l.Line %q bl %d cmpl[0] %q", l.Line, bl, x)
			// In the case of multicompleters, x might have multiple
			// matches. So return the base, not the whole thing.
			if !filepath.IsAbs(cc) {
				x = filepath.Base(x)
			}
			l.Line = l.Line[:bl+1] + x
			Debug("Return is %v", l)
			l.Exact = x
			return nil
		}
	}
	return nil
}

// ReadLine reads until an error occurs
func (l *LineReader) ReadLine() error {
	for {
		Debug("ReadLine: start with %v", l)
		var b [1]byte
		n, err := l.R.Read(b[:])
		if err != nil && err != io.EOF {
			return err
		}
		Debug("ReadLine: got %s, %v, %v", b, n, err)
		if n == 0 {
			return io.EOF
		}
		if err := l.ReadChar(b[0]); err != nil {
			Debug("Readline: %v", l)
			if err == io.EOF || err == ErrEOL {
				Debug("Readline: %v", err)
				return nil
			}
			return err
		}
	}
}
