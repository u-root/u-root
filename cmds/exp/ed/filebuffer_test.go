// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test all functions with func() int signatur
var testTableInt = []struct {
	name   string
	in     *FileBuffer
	exp    any
	method func(f *FileBuffer) int
}{
	{
		name:   "Test Len: Len()=0",
		in:     &FileBuffer{file: []int{}},
		exp:    0,
		method: (*FileBuffer).Len,
	},
	{
		name:   "Test Len: Len()=4",
		in:     &FileBuffer{file: []int{0, 1, 2, 3}},
		exp:    4,
		method: (*FileBuffer).Len,
	},
	{
		name:   "Test GetAddr: GetAddr()=0",
		in:     &FileBuffer{addr: 0},
		exp:    0,
		method: (*FileBuffer).GetAddr,
	},
	{
		name:   "Test GetAddr: GetAddr()=4",
		in:     &FileBuffer{addr: 4},
		exp:    4,
		method: (*FileBuffer).GetAddr,
	},
	{
		name:   "Test Size: Size()=4",
		in:     &FileBuffer{file: []int{0, 1, 2, 3}, buffer: []string{"0", "1", "2", "3"}},
		exp:    4,
		method: (*FileBuffer).Size,
	},
}

func TestBasicFuncsInt(t *testing.T) {
	for _, tt := range testTableInt {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method(tt.in)
			if got != tt.exp {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
		})
	}
}

// Test Dirty
var testTableDirty = []struct {
	name   string
	in     *FileBuffer
	exp    any
	method func(f *FileBuffer) bool
}{
	{
		name:   "Test Dirty: dirty = false",
		in:     &FileBuffer{dirty: false},
		exp:    false,
		method: (*FileBuffer).Dirty,
	},
	{
		name:   "Test Dirty: dirty = true",
		in:     &FileBuffer{dirty: true},
		exp:    true,
		method: (*FileBuffer).Dirty,
	},
}

func TestDirty(t *testing.T) {
	for _, tt := range testTableDirty {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method(tt.in)
			if got != tt.exp {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
		})
	}
}

// Test all functions with func() signatur
var testTable = []struct {
	name   string
	in     *FileBuffer
	exp    any
	method func(f *FileBuffer)
}{
	{
		name:   "Test Touch",
		in:     &FileBuffer{dirty: false, mod: false},
		exp:    &FileBuffer{dirty: true, mod: true},
		method: (*FileBuffer).Touch,
	},
	{
		name:   "Test Start",
		in:     &FileBuffer{mod: true, tmpFile: []int{}, file: []int{0, 1, 2, 3}, tmpAddr: 0, addr: 10, tmpDirty: false, dirty: true},
		exp:    &FileBuffer{mod: false, tmpFile: []int{0, 1, 2, 3}, file: []int{0, 1, 2, 3}, tmpAddr: 10, addr: 10, tmpDirty: true, dirty: true},
		method: (*FileBuffer).Start,
	},
	{
		name:   "Test End",
		in:     &FileBuffer{mod: true, lastFile: []int{}, tmpFile: []int{0, 1, 2, 3}, lastAddr: 0, tmpAddr: 10, lastDirty: false, tmpDirty: true},
		exp:    &FileBuffer{mod: true, lastFile: []int{0, 1, 2, 3}, tmpFile: []int{0, 1, 2, 3}, lastAddr: 10, tmpAddr: 10, lastDirty: true, tmpDirty: true},
		method: (*FileBuffer).End,
	},
	{
		name:   "Test Rewind",
		in:     &FileBuffer{mod: false, file: []int{}, lastFile: []int{0, 1, 2, 3}, addr: 0, lastAddr: 10, dirty: true, lastDirty: false},
		exp:    &FileBuffer{mod: true, file: []int{0, 1, 2, 3}, lastFile: []int{0, 1, 2, 3}, addr: 10, lastAddr: 10, dirty: false, lastDirty: false},
		method: (*FileBuffer).Rewind,
	},
	{
		name:   "Test Clean",
		in:     &FileBuffer{dirty: true, lastDirty: true, lastFile: []int{0, 1, 2, 3}, lastAddr: 10},
		exp:    &FileBuffer{dirty: false, lastDirty: false, lastFile: []int{}, lastAddr: 0},
		method: (*FileBuffer).Clean,
	},
}

func TestBasicFuncs(t *testing.T) {
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.method(tt.in)
			if !reflect.DeepEqual(tt.in, tt.exp) {
				t.Errorf("Expected value: %v, got: %v", tt.exp, tt.in)
			}
		})
	}
}

// Test NewFileBuffer
var testTableNewFileBuffer = []struct {
	name string
	in   []string
	exp  *FileBuffer
}{
	{
		name: "Test NewFileBuffer",
		in:   []string{"0", "1", "2", "3"},
		exp:  &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, dirty: false, mod: false, addr: 0, marks: make(map[byte]int)},
	},
}

func TestNewFileBuffer(t *testing.T) {
	for _, tt := range testTableNewFileBuffer {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFileBuffer(tt.in)
			if !reflect.DeepEqual(got, tt.exp) {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
		})
	}
}

// Test OOB
var testTableOOB = []struct {
	name     string
	in       *FileBuffer
	exp      any
	methodin int
}{
	{
		name:     "Test OOB: return true",
		in:       &FileBuffer{},
		exp:      true,
		methodin: 4,
	},
	{
		name:     "Test OOB: return false",
		in:       &FileBuffer{file: []int{0, 1, 2, 3}},
		exp:      false,
		methodin: 2,
	},
}

func TestOOB(t *testing.T) {
	for _, tt := range testTableOOB {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.OOB(tt.methodin)
			if !reflect.DeepEqual(got, tt.exp) {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
		})
	}
}

// Test GetMust
var testTableGetMust = []struct {
	name      string
	in        *FileBuffer
	exp       string
	methodin1 int
	methodin2 bool
}{
	{
		name:      "Test GetMust",
		in:        &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:       "3",
		methodin1: 3,
		methodin2: true,
	},
}

func TestGetMust(t *testing.T) {
	for _, tt := range testTableGetMust {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.GetMust(tt.methodin1, tt.methodin2)
			if !reflect.DeepEqual(got, tt.exp) {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
		})
	}
}

// Test Get
var testTableGet = []struct {
	name     string
	in       *FileBuffer
	exp1     []string
	exp2     error
	methodin [2]int
}{
	{
		name:     "Test Get: with OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp1:     []string{},
		exp2:     ErrOOB,
		methodin: [2]int{0, 4},
	},
	{
		name:     "Test Get: without OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp1:     []string{"0", "1"},
		exp2:     fmt.Errorf(""),
		methodin: [2]int{0, 1},
	},
}

func TestGet(t *testing.T) {
	for _, tt := range testTableGet {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Get(tt.methodin)
			if err != nil {
				// Using cmp pkf with weird opt because reflect.DeepEqual does not work for empty string arrays
				alwaysEqual := cmp.Comparer(func(_, _ any) bool { return true })
				opt := cmp.FilterValues(func(x, y any) bool {
					vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
					return (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) &&
						(vx.Kind() == reflect.Slice || vx.Kind() == reflect.Map) &&
						(vx.Len() == 0 && vy.Len() == 0)
				}, alwaysEqual)
				if !cmp.Equal(got, tt.exp1, opt) || !reflect.DeepEqual(err.Error(), tt.exp2.Error()) {
					t.Errorf("Expected values: %v, %v, got: %v, %v", tt.exp1, tt.exp2, got, err)
				}
			}
		})
	}
}

// Test Copy
var testTableCopy = []struct {
	name     string
	in       *FileBuffer
	exp      error
	methodin [2]int
}{
	{
		name:     "Test Copy: with OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      ErrOOB,
		methodin: [2]int{0, 4},
	},
	{
		name:     "Test Copy: without OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      nil,
		methodin: [2]int{0, 1},
	},
}

func TestCopy(t *testing.T) {
	for _, tt := range testTableCopy {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Copy(tt.methodin)
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.exp.Error()) {
					t.Errorf("Expected value: %v, got: %v", tt.exp, err)
				}
			}
		})
	}
}

// Test Paste & Insert
var testTablePasteInsert = []struct {
	name     string
	in       *FileBuffer
	exp      error
	methodin int
}{
	{
		name:     "Test Paste & Insert: with OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      ErrOOB,
		methodin: 5,
	},
	{
		name:     "Test Paste & Insert: without OOB error and nlines == 0",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      nil,
		methodin: 1,
	},
	{
		name:     "Test Paste & Insert: without OOB error and nlines != 0",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, cbuf: []string{"0", "1", "2", "3"}},
		exp:      nil,
		methodin: 1,
	},
}

func TestPasteInsert(t *testing.T) {
	for _, tt := range testTablePasteInsert {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Paste(tt.methodin)
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.exp.Error()) {
					t.Errorf("Expected value: %v, got: %v", tt.exp, err)
				}
			}
		})
	}
}

// Test Delete
var testTableDelete = []struct {
	name     string
	in       *FileBuffer
	exp      error
	methodin [2]int
}{
	{
		name:     "Test Delete: with OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      ErrOOB,
		methodin: [2]int{0, 4},
	},
	{
		name:     "Test Delete: without OOB error and addr OOB edge case",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      nil,
		methodin: [2]int{0, 3},
	},
}

func TestDelete(t *testing.T) {
	for _, tt := range testTableDelete {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Delete(tt.methodin)
			if err != nil {
				if !reflect.DeepEqual(err.Error(), tt.exp.Error()) {
					t.Errorf("Expected value: %v, got: %v", tt.exp, err)
				}
			}
		})
	}
}

// Test SetAddr
var testTableSetAddr = []struct {
	name     string
	in       *FileBuffer
	exp      error
	methodin int
}{
	{
		name:     "Test SetAddr: with OOB error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      ErrOOB,
		methodin: 4,
	},
	{
		name:     "Test SetAddr: without OOB error and addr OOB edge case",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}},
		exp:      nil,
		methodin: 0,
	},
}

func TestSetAddr(t *testing.T) {
	for _, tt := range testTableSetAddr {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.SetAddr(tt.methodin)
			if err != nil {
				if err.Error() != tt.exp.Error() {
					t.Errorf("Expected value: %v, got: %v", tt.exp, err)
				}
			} else {
				if tt.methodin != tt.in.addr {
					t.Errorf("Expected value: %v, got: %v", tt.methodin, tt.in.addr)
				}
			}
		})
	}
}

// Test SetMark
var testTableSetMark = []struct {
	name      string
	in        *FileBuffer
	exp       int
	err       error
	methodin1 byte
	methodin2 int
}{
	{
		name:      "Test SetMark: with OOB error",
		in:        &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, marks: make(map[byte]int)},
		exp:       0,
		err:       ErrOOB,
		methodin1: 0,
		methodin2: 5,
	},
	{
		name:      "Test SetMark: without OOB error",
		in:        &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, marks: make(map[byte]int)},
		exp:       3,
		err:       ErrOOB,
		methodin1: 0,
		methodin2: 3,
	},
}

func TestSetMark(t *testing.T) {
	for _, tt := range testTableSetMark {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.SetMark(tt.methodin1, tt.methodin2)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("Expected value: %v, got: %v", tt.exp, err)
				}
			} else {
				if tt.exp != tt.in.marks[0] {
					t.Errorf("Expected value: %v, got: %v", tt.exp, tt.in.marks[0])
				}
			}
		})
	}
}

// Test GetMark
var testTableGetMark = []struct {
	name     string
	in       *FileBuffer
	exp      int
	err      int
	methodin byte
}{
	{
		name:     "Test GetMark: with error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, marks: make(map[byte]int)},
		exp:      0,
		err:      -1,
		methodin: 0x00,
	},
	{
		name:     "Test GetMark: without error",
		in:       &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, marks: make(map[byte]int)},
		exp:      1,
		err:      0,
		methodin: 0x01,
	},
	{
		name:     "Test GetMark: without second error",
		in:       &FileBuffer{marks: make(map[byte]int)},
		exp:      1,
		err:      -1,
		methodin: 0x01,
	},
}

func TestGetMark(t *testing.T) {
	for _, tt := range testTableGetMark {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.marks[1] = 1
			got, err := tt.in.GetMark(tt.methodin)
			if err != nil {
				if got != tt.err {
					t.Errorf("Expected value: %v, got: %v", tt.err, err)
				}
			} else {
				if reflect.DeepEqual(got, tt.in.marks[0]) {
					t.Errorf("Expected value: %v, got: %v", tt.exp, got)
				}
			}
		})
	}
}

// Test Read
var testTableRead = []struct {
	name      string
	in        *FileBuffer
	err       error
	methodin1 int
	methodin2 io.Reader
}{
	{
		name:      "Test Read",
		in:        &FileBuffer{},
		err:       nil,
		methodin1: 0,
		methodin2: &bytes.Buffer{},
	},
}

func TestRead(t *testing.T) {
	for _, tt := range testTableRead {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Read(tt.methodin1, tt.methodin2)
			if err != nil {
				if err != tt.err {
					t.Errorf("Expected value: %v, got: %v", tt.err, err)
				}
			}
		})
	}
}

// Test ReadFile
func CreateFile(name string) string {
	f, err := os.Create(name)
	if err != nil {
		return ""
	}
	defer f.Close()
	f.Write([]byte{0x01})
	return f.Name()
}

var testTableReadFile = []struct {
	name      string
	in        *FileBuffer
	err       error
	methodin1 int
	methodin2 string
}{
	{
		name:      "Test ReadFile: file exists",
		in:        &FileBuffer{},
		err:       nil,
		methodin1: 0,
		methodin2: CreateFile("readfile"),
	},
	{
		name:      "Test ReadFile: file does not exist",
		in:        &FileBuffer{},
		err:       fmt.Errorf("could not read file: open dir: no such file or directory"),
		methodin1: 0,
		methodin2: "dir",
	},
}

func TestReadFile(t *testing.T) {
	for _, tt := range testTableReadFile {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.ReadFile(tt.methodin1, tt.methodin2)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("Expected value: %v, got: %v", tt.err, err)
				}
			}
			os.Remove(tt.methodin2)
		})
	}
}

// Test FileToBuffer
var testTableFileToBuffer = []struct {
	name     string
	in       *FileBuffer
	exp      *FileBuffer
	err      error
	methodin string
}{
	{
		name:     "Test FileToBuffer: file exists",
		in:       &FileBuffer{},
		exp:      &FileBuffer{mod: true, file: []int{0}},
		err:      fmt.Errorf(""),
		methodin: CreateFile("filetobuffer"),
	},
	{
		name:     "Test FileToBuffer: file does not exist",
		in:       &FileBuffer{},
		exp:      NewFileBuffer(nil),
		err:      fmt.Errorf("could not read file: open %s: no such file or directory", "dir"),
		methodin: "dir",
	},
}

func TestFileToBuffer(t *testing.T) {
	for _, tt := range testTableFileToBuffer {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileToBuffer(tt.methodin)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("Expected value: %v, got: %v", tt.err, err)
				}
			} else if reflect.DeepEqual(tt.exp, got) {
				t.Errorf("Expected value: %v, got: %v", tt.exp, got)
			}
			os.Remove(tt.methodin)
		})
	}
}
