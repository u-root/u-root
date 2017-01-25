// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	r, err := NewcReader(bytes.NewReader(testCPIO))

	if err != nil {
		t.Error(err)
	}
	var f *File
	var i int
	for f, err = r.RecRead(); err == nil; f, err = r.RecRead() {
		if f.String() != testResult[i] {
			t.Errorf("Value %d: got \n%s, want \n%s", i, f.String(), testResult[i])
		}
		t.Logf("Value %d: got \n%s, want \n%s", i, f.String(), testResult[i])
		i++
	}
}

func TestReadWrite(t *testing.T) {
	r, err := NewcReader(bytes.NewReader(testCPIO))
	if err != nil {
		t.Fatalf("Creating testCPIO reader: %v", err)
	}
	f, err := r.RecReadAll()
	if err != nil {
		t.Fatalf("Reading testCPIO reader: %v", err)
	}

	var b = &bytes.Buffer{}
	w, err := NewcWriter(b)
	if err != nil {
		t.Fatalf("TestReadWrite: new writer: %v", err)
	}

	if err := w.RecWriteAll(f); err != nil {
		t.Fatalf("%v", err)
	}

	if r, err = NewcReader(bytes.NewReader(b.Bytes())); err != nil {
		t.Errorf("%v", err)
	}
	var tf []*File
	if tf, err = r.RecReadAll(); err != nil {
		t.Fatalf("TestReadWrite: reading generated data: %v", err)
	}

	// We have to reread the original since the Data elements in the struct
	// have been consumed to write the second []byte
	if r, err = NewcReader(bytes.NewReader(testCPIO)); err != nil {
		t.Error(err)

	}
	if f, err = r.RecReadAll(); err != nil {
		t.Fatalf("Reading testCPIO reader second time: %v", err)
	}

	// Now check a few things: arrays should be same length, Headers should match,
	// names should be the same, and data should be the same. If this all works,
	// it means we read in serialized data, wrote it out, read it in, and the
	// structs all matched.
	if len(f) != len(tf) {
		t.Fatalf("[]file len from testCPIO %v and generated %v are not the same and should be", len(f), len(tf))
	}
	for i := range f {
		if f[i].Header != tf[i].Header {
			t.Errorf("index %d: testCPIO Header\n%v\ngenerated Header\n%v\n", i, f[i].Header, tf[i].Header)
		}
		if f[i].Name != tf[i].Name {
			t.Errorf("index %d: testCPIO Name\n%v\ngenerated Name\n%v\n", i, f[i].Name, tf[i].Name)
		}
		s, err := ioutil.ReadAll(f[i].Data)
		if err != nil {
			t.Errorf("index %d(%s): can't read from the source: %v", i, f[i].Name, err)
		}
		d, err := ioutil.ReadAll(tf[i].Data)
		if err != nil {
			t.Errorf("index %d(%s): can't read from the dest: %v", i, tf[i].Name, err)
		}
		if !reflect.DeepEqual(s, d) {
			t.Errorf("index %d: d(%s) is %v, s(%s) wanted %v", i, tf[i].Name, d, f[i].Name, s)
		}
	}
}

func TestBad(t *testing.T) {
	_, err := NewcReader(bytes.NewReader(badCPIO))
	t.Logf("NewcReader with badCPIO error is %v", err)

	if err == nil {
		t.Errorf("Wanted EOF err, got nil")
	}

	_, err = NewcReader(bytes.NewReader(badMagicCPIO))
	t.Logf("NewcReader with badMagicCPIO error is %v", err)

	if err == nil {
		t.Errorf("Wanted bad magic err, got nil")
	}
}
