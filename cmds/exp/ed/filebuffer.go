// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// filebuffer.go - defines the FileBuffer object
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// A FileBuffer manages a file being edited.
// A FileBuffer never deletes/modifies anything directly until it is replaced.
// It keeps a map of known lines to the current buffer.
// Note: FileBuffer is 0-addressed lines, so off-by-one from what `ed` expects.
type FileBuffer struct {
	cbuf      []string // cut buffer
	buffer    []string // all lines we know about, they never get delited
	file      []int    // sequence of buffer lines
	lastFile  []int    // used for undo capability
	tmpFile   []int    // used for undo capability
	dirty     bool     // tracks if the file has been modifed
	lastDirty bool     // used for undo capability
	tmpDirty  bool     // used for undo capability
	mod       bool     // mod is like dirty, but can be reset for transactions
	addr      int      // current file address
	lastAddr  int      // last address (for undo)
	tmpAddr   int      // last address (for undo)
	marks     map[byte]int
}

// NewFileBuffer creats a new FileBuffer object
func NewFileBuffer(in []string) *FileBuffer {
	f := &FileBuffer{
		buffer: in,
		file:   []int{},
		dirty:  false,
		mod:    false,
		addr:   0,
		marks:  make(map[byte]int),
	}
	for i := range f.buffer {
		f.file = append(f.file, i)
	}
	return f
}

// ErrOOB line is out of bounds
var ErrOOB = fmt.Errorf("line is out of bounds")

// ErrINV address is invalid
var ErrINV = fmt.Errorf("invalid address")

// OOB checks if a line is out of bounds
func (f *FileBuffer) OOB(l int) bool {
	if l < 0 || l >= f.Len() {
		return true
	}
	return false
}

// GetMust gets a specified line, to be used when we know it's safe (no error return)
// if "set" is true, sets the current line pointer
func (f *FileBuffer) GetMust(line int, set bool) string {
	if set {
		f.addr = line
	}
	return f.buffer[f.file[line]]
}

// Get a specified line range
// this updates the current line pointer
func (f *FileBuffer) Get(r [2]int) (lines []string, e error) {
	if f.OOB(r[0]) || f.OOB(r[1]) {
		e = ErrOOB
		return
	}
	for l := r[0]; l <= r[1]; l++ {
		lines = append(lines, f.buffer[f.file[l]])
		f.addr = l
	}
	return
}

// Copy lines into the cut buffer
func (f *FileBuffer) Copy(r [2]int) (e error) {
	var lines []string
	if lines, e = f.Get(r); e != nil {
		return
	}
	f.cbuf = lines
	return
}

// Paste lines from cut buffer insert at line
func (f *FileBuffer) Paste(line int) (e error) {
	e = f.Insert(line, f.cbuf)
	return
}

// Delete unmaps lines from the file
func (f *FileBuffer) Delete(r [2]int) (e error) {
	blines := []int{}
	for l := r[0]; l <= r[1]; l++ {
		if f.OOB(l) {
			return ErrOOB
		}
		blines = append(blines, f.file[l])
	}
	f.cbuf, _ = f.Get(r) // this shouldn't fail here, if it does we've got a bigger problem
	for _, b := range blines {
		for i, l := range f.file {
			if l == b {
				f.file = append(f.file[:i], f.file[i+1:]...)
				break
			}
		}
	}
	f.Touch()
	f.addr = r[0] + 1
	if f.OOB(f.addr) {
		f.addr = 0
	}
	return
}

// Insert adds nlines to buffer and inserts them at line
func (f *FileBuffer) Insert(line int, nlines []string) (e error) {
	if line != f.Len() && f.OOB(line) { // if line == f.Len() we append to the end
		return ErrOOB
	}
	if len(nlines) == 0 {
		return
	}
	first := len(f.buffer)
	f.buffer = append(f.buffer, nlines...)
	nf := []int{}
	for i := first; i < len(f.buffer); i++ {
		nf = append(nf, i)
	}
	f.file = append(f.file[:line], append(nf, f.file[line:]...)...)
	f.Touch()
	f.addr = line + len(nlines) - 1
	return
}

// Len returns the current file length
func (f *FileBuffer) Len() int {
	return len(f.file)
}

// Dirty returns whether the file has changed
func (f *FileBuffer) Dirty() bool {
	return f.dirty
}

// GetAddr gets the current file addr
func (f *FileBuffer) GetAddr() int {
	return f.addr
}

// SetAddr sets the current addr, errors if OOB
func (f *FileBuffer) SetAddr(i int) (e error) {
	if f.OOB(i) {
		return ErrOOB
	}
	f.addr = i
	return
}

// Clean resets the dirty flag
func (f *FileBuffer) Clean() {
	f.dirty = false
	f.lastDirty = false
	f.lastFile = []int{}
	f.lastAddr = 0
}

// FileToBuffer reads a file and creates a new FileBuffer from it
func FileToBuffer(file string) (fb *FileBuffer, e error) {
	fb = NewFileBuffer(nil)
	e = fb.ReadFile(0, file)
	if e == nil {
		fb.dirty = false
	}
	return
}

// SetMark sets a mark (by byte name) in the FileBuffer for later use
func (f *FileBuffer) SetMark(c byte, l int) (e error) {
	if f.OOB(l) {
		e = ErrOOB
		return
	}
	f.marks[c] = f.file[l]
	return
}

// GetMark gets a mark from the FileBuffer (by byte name)
func (f *FileBuffer) GetMark(c byte) (l int, e error) {
	bl, ok := f.marks[c]
	if !ok {
		return -1, fmt.Errorf("no such mark: %c", c)
	}
	for i := 0; i < f.Len(); i++ {
		if f.file[i] == bl {
			l = i
			return
		}
	}
	return -1, fmt.Errorf("mark was cleared: %c", c)
}

// Size return the size (in bytes) of the current file buffer
func (f *FileBuffer) Size() (s int) {
	for _, i := range f.file {
		s += len(f.buffer[i])
	}
	return
}

// Read reads in from an io.Reader interface and inserts at the current line address
func (f *FileBuffer) Read(line int, r io.Reader) (e error) {
	b := []string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		b = append(b, s.Text())
	}
	e = f.Insert(line, b)
	return
}

// ReadFile reads in a file and inserts it at the current line address
func (f *FileBuffer) ReadFile(line int, file string) (e error) {
	var fh *os.File
	if fh, e = os.Open(file); e != nil {
		e = fmt.Errorf("could not read file: %w", e)
		return
	}
	defer fh.Close()

	e = f.Read(line, fh)
	return
}

// Start a transaction
func (f *FileBuffer) Start() {
	f.mod = false
	f.tmpFile = make([]int, len(f.file))
	copy(f.tmpFile, f.file)
	f.tmpAddr = f.addr
	f.tmpDirty = f.dirty
}

// End a transaction
func (f *FileBuffer) End() {
	if f.mod {
		f.lastFile = f.tmpFile
		f.lastAddr = f.tmpAddr
		f.lastDirty = f.tmpDirty
	}
}

// Rewind restores the previous file
func (f *FileBuffer) Rewind() {
	if f.Dirty() || f.lastDirty {
		f.addr = f.lastAddr
		f.file = f.lastFile
		f.dirty = f.lastDirty
		f.mod = true
	}
}

// Touch is the correct way (even internally) to set the dirty & modified bits
func (f *FileBuffer) Touch() {
	f.dirty = true
	f.mod = true
}
