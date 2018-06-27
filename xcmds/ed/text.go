// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
)

// start and end are used like slices.
// so 1,1 defines a point. 1, 2 is line 1.
// 1, $ is all the lines where $ is len(lines)+1
type file struct {
	dot   int
	dirty bool
	lines []int
	data  []byte
	pat   string
}

func makeLines(n []byte) [][]byte {
	if len(n) == 0 {
		return nil
	}
	var (
		data  [][]byte
		sawnl bool
		pos   int
		i     int
	)
	for i = range n {
		if sawnl {
			sawnl = false
			l := make([]byte, i-pos)
			copy(l, n[pos:i])
			data = append(data, l)
			pos = i
		}
		if n[i] == '\n' {
			sawnl = true
		}
	}

	if pos <= i {
		l := make([]byte, len(n[pos:]))
		copy(l, n[pos:])
		data = append(data, l)
	}
	debug("makeLines: %v", data)
	return data
}

func (f *file) String() string {
	return fmt.Sprintf("%d %v %s", f.dot, f.lines, string(f.data))
}

func (f *file) fixLines() {
	f.lines = nil
	if len(f.data) == 0 {
		return
	}
	f.lines = []int{0}

	lines := 1
	var sawnl bool
	for i, v := range f.data {
		if sawnl {
			lines++
			f.lines = append(f.lines, i)
			sawnl = false
		}
		if v == '\n' {
			sawnl = true
		}
	}
	if f.dot < 1 {
		f.dot = 1
	}
	if f.dot > lines {
		f.dot = lines
	}
	//	f.start, f.end = f.dot, f.dot
}

// SliceXX returns half-closed indexes of the segment corresponding to [f.start, f.end]
// as in the Go slice style.
// The intent is that start and end act like slice indices,
// i.e. the slice returned should look as though
// you'd taken lines[f.start:f.end]
// Note that a zero length slice is fine, and is used
// for, e.g., inserting.
// Also, note, the indices need to be changed from 1-relative
// to 0-relative, since the API at all levels is 1-relative.
func (f *file) SliceX(startLine, endLine int) (int, int) {
	var end, start int
	debug("SliceX: f %v", f)
	defer debug("SliceX done: f %v, start %v, end %v", f, start, end)
	if startLine > len(f.lines) {
		//f.start = len(f.lines)
		//f.end = f.start
		return len(f.data), len(f.data)
	} else {
		if startLine > 0 {
			start = f.lines[startLine-1]
		}
	}
	end = start
	if endLine < startLine {
		endLine = startLine
	}
	if endLine != startLine {
		if endLine >= len(f.lines) {
			//f.end = len(f.lines)
			end = len(f.data)
			return start, end
		}
		end = f.lines[endLine-1]
	}
	debug("SliceX: f %v, start %v end %v\n", f, start, end)
	return start, end
}

func (f *file) Slice(startLine, endLine int) ([]byte, int, int) {
	s, e := f.SliceX(startLine, endLine)
	return f.data[s:e], s, e
}

func (f *file) Replace(n []byte, startLine, endLine int) (int, error) {
	defer debug("Replace done: f %v", f)
	if f.data != nil {
		start, end := f.SliceX(startLine, endLine)
		debug("replace: f %v start %v end %v", f, start, end)
		pre := f.data[0:start]
		post := f.data[end:]
		debug("replace: pre is %v, post is %v\n", pre, post)
		var b bytes.Buffer
		b.Write(pre)
		b.Write(n)
		b.Write(post)
		f.data = b.Bytes()
	} else {
		f.data = n
	}
	f.fixLines()
	return len(n), nil
}

// Read reads strings from the file after .
// If there are lines already the new lines are inserted after '.'
// dot is unchanged.
func (f *file) Read(r io.Reader, startLine, endLine int) (int, error) {
	debug("Read: r %v, startLine %v endLine %v", r, startLine, endLine)
	d, err := ioutil.ReadAll(r)
	debug("ReadAll returns %v, %v", d, err)
	if err != nil {
		return -1, err
	}
	// in ed, for text read, the line # is a point, not a line.
	// For other commands, it's a line. Hence for reading
	// line x, we need to make a point at line x+1 and
	// read there.
	startLine = startLine + 1
	return f.Replace(d, startLine, endLine)
}

// Write writes the lines out from start to end, inclusive.
// dot is unchanged.
func (f *file) Write(w io.Writer, startLine, endLine int) (int, error) {
	if endLine < startLine || startLine < 1 || endLine > len(f.lines)+1 {
		return -1, fmt.Errorf("file is %d lines and [start, end] is [%d, %d]", len(f.lines), startLine, endLine)
	}

	start, end := f.SliceX(startLine, endLine)
	amt, err := w.Write(f.data[start:end])

	return amt, err
}

func (f *file) Print(w io.Writer, start, end int) (int, error) {
	debug("Print %v %v %v", f.data, start, end)
	i, err := f.Write(w, start, end)
	if err != nil && start < end {
		f.dot = end + 1
	}
	return i, err
}

// Write writes the lines out from start to end, inclusive.
// dot is unchanged.
func (f *file) WriteFile(n string, startLine, endLine int) (int, error) {
	out, err := os.Create(n)
	if err != nil {
		return -1, err
	}
	defer out.Close()
	return f.Write(out, startLine, endLine)
}

// Sub replaces the regexp with a different one.
func (f *file) Sub(re, n, opt string, startLine, endLine int) error {
	debug("Sub re %s n %s opt %s\n", re, n, opt)
	if re == "" {
		return fmt.Errorf("Empty RE")
	}
	r, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	var opts = make(map[byte]bool)
	for i := range opt {
		opts[opt[i]] = true
	}

	o, start, end := f.Slice(startLine, endLine)
	debug("Slice from [%v,%v] is [%v, %v] %v", startLine, endLine, start, end, string(o))
	// Lines can turn into two lines. All kinds of stuff can happen.
	// So
	// Break it into lines
	// make copies of those lines.
	// for each line, do the sub
	// paste it back together
	// put it back in to f.data
	b := makeLines(o)
	if b == nil {
		return nil
	}
	for i := range b {
		var replaced bool
		debug("Sub: before b[i] is %v", b[i])
		b[i] = r.ReplaceAllFunc(b[i], func(o []byte) []byte {
			debug("Sub: func called with %v", o)
			if opts['g'] || !replaced {
				f.dirty = true
				replaced = true
				return []byte(n)
			}
			return o
		})
		debug("Sub: after b[i] is %v", b[i])
	}

	debug("replaced o %v with n %v\n", o, b)
	var repl = make([]byte, start)
	copy(repl, f.data[0:start])
	for _, v := range b {
		repl = append(repl, v...)
	}

	f.data = append(repl, f.data[end:]...)
	f.fixLines()
	if opts['p'] {
		_, err = f.Write(os.Stdout, startLine, endLine)
	}
	return err
}

func (f *file) Dirty(d bool) {
	f.dirty = d
}

func (f *file) IsDirty() bool {
	return f.dirty
}

// Range returns the closed range
// of lines.
func (f *file) Range() (int, int) {
	if f.lines == nil || len(f.lines) == 0 {
		return 0, 0
	}
	return 1, len(f.lines)
}

func (f *file) Equal(e Editor) error {
	var errors string
	g := e.(*file)
	// we should verify dot but let's not do that just yet, we don't
	// have it right.
	if !reflect.DeepEqual(f.lines, g.lines) {
		errors += fmt.Sprintf("%v vs %v: lines don't match", f.lines, g.lines)
	}
	if len(f.data) != len(g.data) {
		errors += fmt.Sprintf("data len  differs: %v vs %v", len(f.data), len(g.data))
	} else {
		for i := range f.data {
			if f.data[i] != g.data[i] {
				errors += fmt.Sprintf("%d: %v vs %v: data doesn't match", i, f.data[i], g.data[i])
			}
		}
	}
	if errors != "" {
		return fmt.Errorf(errors)
	}
	return nil
}

func (f *file) Dot() int {
	return f.dot
}

func (f *file) Move(dot int) {
	f.dot = dot
}

func NewTextEditor(a ...editorArg) (Editor, error) {
	var f = &file{}

	for _, v := range a {
		if err := v(f); err != nil {
			return nil, err
		}
	}
	f.dot = 1
	return f, nil
}
