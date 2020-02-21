// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// NewerLineReader reads a char and changes state. It does no I/O.
// It maintains information such that a caller can figure out
// what to do with a line.
type NewerLineReader struct {
	// Completer for this LineReader
	C Completer
	F Completer
	// Prompt is a prompt
	Prompt string
	// Lines holds data as it is read.
	Line string
	// FullLine holds the last full line read
	FullLine string
	// Exact is the exact match.
	// It can be "" if there is not one.
	Exact string
	// Candidates are any completion candidates.
	// The UI can decide how to handle them.
	Candidates []string
	EOF        bool
	Fields     int
	// we only try to do completions if Tabbed is true
	Tabbed bool
}

// NewNewerLineReader returns a LineReader.
func NewNewerLineReader(c, f Completer) *NewerLineReader {
	return &NewerLineReader{C: c, F: f}
}

// ReadChar reads one character and processes it. It is inflexible by design.
func (l *NewerLineReader) ReadChar(b byte) (err error) {
	defer func() {
		l.Fields = len(strings.Fields(l.Line))
	}()
	c := l.C
	Debug("l.Fields %d", l.Fields)
	if l.Fields > 1 || strings.Trim(l.Line, " ") != l.Line {
		c = l.F
	}
	Debug("NewerLineReader: start with %v", l)
	l.Exact, l.Candidates = "", []string{}
	switch b {
	case '\t':
		if !l.Tabbed {
			l.Tabbed = true
			return nil
		}
	default:
		l.Tabbed = false
	}
	switch b {
	default:
		Debug("NewerLineReader.Just add it to line and pipe")
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
		Debug("ERREOL")
		return ErrEOL
	case ' ':
		l.Line += string(b)
		return nil
	case '\t':
		ll := len(l.Line)
		// Special case: if there's nothing in the line yet,
		// just return.
		if ll == 0 {
			return nil
		}
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
		x, cmpl, err := c.Complete(cc)
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
		// see if there is enough for a common prefix.
		if x == "" {
			p := Prefix(cmpl)
			if p != "" {
				l.Line = l.Line[:bl+1] + p
			}
		}
	}
	return nil
}

// ReadLine reads until an error occurs
func (l *NewerLineReader) ReadLine(r io.Reader, w io.Writer) error {
	fmt.Print(l.Prompt)
	for {
		Debug("ReadLine: start with %v", l)
		var b [1]byte
		n, err := r.Read(b[:])
		if err != nil && err != io.EOF {
			return err
		}
		Debug("ReadLine: got %s, %v, %v", b, n, err)
		if n == 0 {
			return io.EOF
		}
		p := l.Line
		if err := l.ReadChar(b[0]); err != nil {
			Debug("Readline: %v", l)
			if err == io.EOF || err == ErrEOL {
				Debug("Readline: %v", err)
				l.FullLine = l.Line
				l.Line = ""
				return nil
			}
			return err
		}
		Debug("readline, exact %q, cand %q", l.Exact, l.Candidates)
		if len(l.Candidates) > 1 {
			c := l.Candidates
			if l.Exact != "" {
				c = append([]string{l.Exact}, l.Candidates...)
			}
			if _, err := fmt.Fprintf(w, "\r\n%s\r\n%s%s", c, l.Prompt, l.Line); err != nil {
				log.Printf("Showing completions: %v", err)
			}
			continue
		}
		// How we handle this depends on whether it is a directory
		if s, err := os.Stat(l.Exact); err == nil && s.IsDir() {
			l.Line += "/"
		}
		if p != l.Line {
			if _, err := w.Write([]byte("\r" + strings.Repeat(" ", len(l.Prompt+p)) + "\r" + l.Prompt + l.Line)); err != nil {
				log.Print(err)
			}
		}
	}
}
