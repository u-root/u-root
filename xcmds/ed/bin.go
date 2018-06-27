// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

type bin struct {
	dot   int
	dirty bool
	data  []byte
}

func (f *bin) String() string {
	return fmt.Sprintf("%d %v", f.dot, len(f.data))
}

func (f *bin) Replace(n []byte, start, end int) (int, error) {
	defer debug("Replace done: f %v", f)
	if f.data != nil {
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
	return len(n), nil
}

// Read reads strings from the file after .
// If there are lines already the new lines are inserted after '.'
// dot is unchanged.
func (f *bin) Read(r io.Reader, start, end int) (int, error) {
	debug("Read: r %v, start %v end %v", r, start, end)
	d, err := ioutil.ReadAll(r)
	debug("ReadAll returns %v, %v", d, err)
	if err != nil {
		return -1, err
	}
	return f.Replace(d, start, end)
}

// Write writes the lines out from start to end, inclusive.
// dot is unchanged.
func (f *bin) Write(w io.Writer, start, end int) (int, error) {
	if end < start || start < 0 || end > len(f.data) {
		return -1, fmt.Errorf("file is %d lines and [start, end] is [%d, %d]", len(f.data), start, end)
	}

	amt, err := w.Write(f.data[start:end])

	return amt, err
}

func (f *bin) Print(w io.Writer, start, end int) (int, error) {
	if end > len(f.data) {
		end = len(f.data)
	}
	for l := start; l < end; l += 16 {
		fmt.Fprintf(w, "%08x ", l)
		for b := l; b < end && b < l+16; b++ {
			fmt.Fprintf(w, "%02x ", f.data[b])
		}
		for b := l; b < end && b < l+16; b++ {
			c := f.data[b]
			if c < 32 || c > 126 {
				fmt.Fprintf(w, ".")
			} else {
				fmt.Fprintf(w, "%c", c)
			}
		}
		fmt.Printf("\n")
	}
	f.dot = end
	return end - start, nil
}

// Sub replaces the regexp with a different one.
// I don't know how to do a regexp for bin. So, for now
// it's a simple substitution, and we'll start by assuming
// the same length
func (f *bin) Sub(x, y, opt string, start, end int) error {
	debug("Slice from [%v,%v] %v", start, end, string(f.data[start:end]))
	if len(x) != len(y) {
		return fmt.Errorf("For now, old and new must be same len")
	}
	l := len(x) / 2
	xb := make([]byte, l)
	yb := make([]byte, l)

	for i := range xb {
		n := x[i*2 : i*2+2]
		c, err := strconv.ParseUint(n, 16, 8)
		if err != nil {
			return fmt.Errorf("%s is not a hex number", n)
		}
		xb[i] = uint8(c)
		n = y[i*2 : i*2+2]
		c, err = strconv.ParseUint(n, 16, 8)
		if err != nil {
			return fmt.Errorf("%s is not a hex number", n)
		}
		yb[i] = uint8(c)
	}

	for i := start; i < end; i++ {
		if reflect.DeepEqual(f.data[i:i+l], xb) {
			copy(f.data[i:], yb)
		}
	}
	return nil
}

func (f *bin) Dirty(d bool) {
	f.dirty = d
}

func (f *bin) IsDirty() bool {
	return f.dirty
}

func (f *bin) Range() (int, int) {
	return 0, len(f.data)
}

func (f *bin) Equal(e Editor) error {
	g := e.(*file)
	// we should verify dot but let's not do that just yet, we don't
	// have it right.
	if /*f.dot != g.dot || */ !reflect.DeepEqual(f.data, g.data) {
		return fmt.Errorf("%v is not the same as %v", f.String(), g.String())
	}
	return nil
}

func (f *bin) Dot() int {
	return f.dot
}

func (f *bin) Move(dot int) {
	f.dot = dot
}

func NewBinEditor(a ...editorArg) (Editor, error) {
	var f = &bin{}
	for _, v := range a {
		if err := v(f); err != nil {
			return nil, err
		}
	}
	f.dot = 0
	return f, nil
}
