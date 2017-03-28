// iopyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ed is a simple line-oriented editor
//
// Synopsis:
//     dd
//
// Description:
//
// Options:
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewTextFile(t *testing.T) {
	var test = []struct {
		d          string
		start, end int
		f          file
	}{
		{d: "a\nb\nc\n", start: 1, end: 3, f: file{dot: 1, data: []byte("a\nb\nc\n"), lines: []int{0, 2, 4}}},
		// I'm not sure this is right.
		{d: "", start: 1, end: 1, f: file{dot: 1}},
	}

	for _, v := range test {
		f, err := NewTextEditor(readerio(bytes.NewBufferString(v.d)))
		if err != nil {
			t.Errorf("%v: want nil, got %v", v.d, err)
			continue
		}
		if err := v.f.Equal(f); err != nil {
			t.Errorf("%v vs %v: want nil, got %v", v.f, *f.(*file), err)
		}

	}

}

func TestReadTextFile(t *testing.T) {
	var test = []struct {
		d     string
		where int
		f     file
	}{
		//{"a\nb\nc\n", 1, file{dot: 1, data: []byte("a\na\nb\nc\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
		//{"a\nb\nc\n", 4, file{dot: 4, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
		//{"a\nb\nc\n", 2, file{dot: 2, data: []byte("a\nb\na\nb\nc\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{"a\nb\nc\n", 0, file{dot: 1, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
	}
	debug = t.Logf
	for _, v := range test {
		r := bytes.NewBufferString(v.d)
		f, err := NewTextEditor(readerio(r))
		// We are adding the file after a 0-length slice
		// We want dot to be at the point we added it.
		f.Move(v.where)
		r = bytes.NewBufferString(v.d)
		_, err = f.Read(r, v.where, v.where)
		if err != nil {
			t.Errorf("Error reading %v: %v", v.d, err)
			continue
		}
		if err := v.f.Equal(f); err != nil {
			t.Errorf("%v vs. %v: want nil, got %v", v.f, *f.(*file), err)
		}

	}

}

func TestWriteTextFile(t *testing.T) {
	var test = []struct {
		d          string
		start, end int
		err        string
		f          file
	}{
		{d: "a\nb\n", start: 1, end: 3, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{d: "b\nc\n", start: 2, end: 4, f: file{dot: 3, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{d: "b\n", start: 2, end: 3, f: file{dot: 2, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{d: "a\n", start: 1, end: 2, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{err: "file is 6 lines and [start, end] is [40, 1]", start: 40, end: 1, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{err: "file is 6 lines and [start, end] is [1, 40]", start: 1, end: 40, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{err: "file is 6 lines and [start, end] is [40, 60]", start: 40, end: 60, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
	}

	debug = t.Logf
	for _, v := range test {
		n, err := ioutil.TempFile("", "ed")
		if err != nil {
			t.Fatalf("TempFile: want nil, got %v", err)
		}
		defer os.Remove(n.Name())

		_, err = v.f.WriteFile(n.Name(), v.start, v.end)
		if err != nil && v.err != err.Error() {
			t.Errorf("Error Writing %v, want err \n%v got err \n%v", v.f, v.err, err)
			continue
		}
		check, err := ioutil.ReadFile(n.Name())
		if err != nil {
			t.Errorf("Error reading back %v: %v", n.Name(), err)
			continue
		}
		if string(check) != v.d {
			t.Errorf("Error reading back: want %v, got %v", v.d, string(check))
		}

	}

}

func TestSimpleTextCommand(t *testing.T) {
	var test = []struct {
		c          string
		start, end int
		err        string
		f          *file
	}{
		{c: "Z", err: "Z: unknown command", start: 1, end: 2, f: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{c: ".Z", err: "Z: unknown command", start: 1, end: 2, f: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{c: "4Z", err: "Z: unknown command", start: 1, end: 2, f: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{c: "1,3Z", err: "Z: unknown command", start: 1, end: 2, f: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{c: "/a/Z", err: "Z: unknown command", start: 1, end: 2, f: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}}},
	}

	for _, v := range test {
		err := DoCommand(v.f, v.c)
		if err != nil && v.err != err.Error() {
			t.Error(err.Error())
		}
	}

}

func TestFileTextCommand(t *testing.T) {
	var test = []struct {
		c          string
		start, end int
		err        string
		res        string
		f          file
	}{
		{c: "e%s", start: 1, end: 7, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}}},
		{c: "0r%s", start: 1, end: 12, f: file{dot: 1, data: []byte("a\nb\nc\na\nb\nc\na\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22}}},
	}

	n, err := ioutil.TempFile("", "ed")
	if err != nil {
		t.Fatalf("TempFile: want nil, got %v", err)
	}
	defer os.Remove(n.Name())

	d := "a\nb\nc\na\nb\nc\n"
	if err := ioutil.WriteFile(n.Name(), []byte(d), 0666); err != nil {
		t.Fatalf("writing test data %s: %v", d, err)
	}

	f, err := NewTextEditor(readerio(bytes.NewBufferString(d)))
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, v := range test {
		cmd := fmt.Sprintf(v.c, n.Name())
		err := DoCommand(f, cmd)
		t.Logf("f after %v command is %v", v.c, f)
		if err != nil && v.err != err.Error() {
			t.Error(err.Error())
		}
		if err := v.f.Equal(f); err != nil {
			t.Errorf("%v: want nil, got %v", cmd, err)
		}
	}

}

func TestFileDelete(t *testing.T) {
	var test = []struct {
		start, end int
		fi         *file
		fo         *file
	}{
		{
			start: 1, end: 2,
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("b\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8}},
		},
		{
			start: 1, end: 4,
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("a\nb\nc"), lines: []int{0, 2, 4}},
		},
	}

	debug = t.Logf
	for _, v := range test {
		v.fi.Replace([]byte{}, v.start, v.end)
		if v.fi.dot != v.fo.dot {
			t.Errorf("Delete [%d, %d]: want dot %v, got %v", v.start, v.end, v.fo.dot, v.fi.dot)
		}
		if !reflect.DeepEqual(v.fi.lines, v.fo.lines) {
			t.Errorf("Delete [%d, %d]: want %v, got %v", v.start, v.end, v.fo.lines, v.fi.lines)
		}
		if !reflect.DeepEqual(v.fi.data, v.fo.data) {
			t.Errorf("Delete [%d, %d]: want %v, got %v", v.start, v.end, v.fo.data, v.fi.data)
		}
	}

}

func TestFileDCommand(t *testing.T) {
	var test = []struct {
		c   string
		fi  *file
		err string
		fo  *file
	}{
		{
			c:  "1d",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("b\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8}},
		},
		{
			c:  "1,3d",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("a\nb\nc"), lines: []int{0, 2, 4}},
		},
	}

	debug = t.Logf
	for _, v := range test {
		err := DoCommand(v.fi, v.c)
		if v.err == "" && err != nil {
			t.Errorf("%v: want nil, got %v", v.c, err)
			continue
		}
		if err != nil && err.Error() != v.err {
			t.Errorf("%v: want %v, got %v", v.c, v.err, err)
			continue
		}
		if v.fi.dot != v.fo.dot {
			t.Errorf("want dot %v, got %v", v.fo.dot, v.fi.dot)
		}
		if !reflect.DeepEqual(v.fi.lines, v.fo.lines) {
			t.Errorf("%v: want %v, got %v", v.c, v.fo.lines, v.fi.lines)
		}
		if !reflect.DeepEqual(v.fi.data, v.fo.data) {
			t.Errorf("%v: want %v, got %v", v.c, v.fo.data, v.fi.data)
		}
	}

}

func TestFileSCommand(t *testing.T) {
	var test = []struct {
		c   string
		fi  *file
		err string
		fo  *file
	}{
		{
			c:  "1s/a/b/",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("b\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}},
		},
		{
			c:  "1,2s/a/A/g",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("A\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
		},
		{
			c:  "1,$s/a/A/g",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("A\nb\nc\nA\nb\nc\n"), lines: []int{0, 2, 4, 6, 8, 10}},
		},
		{
			c:  "1,$s/a/A/g",
			fi: &file{dot: 1, data: []byte("a\nb\nc\na\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
			fo: &file{dot: 1, data: []byte("A\nb\nc\nA\nb\nc"), lines: []int{0, 2, 4, 6, 8, 10}},
		},
	}

	debug = t.Logf
	for _, v := range test {
		err := DoCommand(v.fi, v.c)
		if v.err == "" && err != nil {
			t.Errorf("%v: want nil, got %v", v.c, err)
			continue
		}
		if err != nil && err.Error() != v.err {
			t.Errorf("%v: want %v, got %v", v.c, v.err, err)
			continue
		}
		if v.fi.dot != v.fo.dot {
			t.Errorf("want dot %v, got %v", v.fo.dot, v.fi.dot)
		}
		if !reflect.DeepEqual(v.fi.lines, v.fo.lines) {
			t.Errorf("%v: want %v, got %v", v.c, v.fo.lines, v.fi.lines)
		}
		if !reflect.DeepEqual(v.fi.data, v.fo.data) {
			t.Errorf("%v: want %v, got %v", v.c, v.fo.data, v.fi.data)
		}
	}

}
