// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"testing"
)

func TestNormalize(t *testing.T) {
	for _, tt := range []struct {
		path string
		want string
	}{
		{
			path: "/foo/bar",
			want: "foo/bar",
		},
		{
			path: "foo////bar",
			want: "foo/bar",
		},
		{
			path: "/foo/bar/../baz",
			want: "foo/baz",
		},
		{
			path: "foo/bar/../baz",
			want: "foo/baz",
		},
		{
			path: "./foo/bar",
			want: "foo/bar",
		},
		{
			path: "foo/../../bar",
			want: "../bar",
		},
		{
			path: "",
			want: ".",
		},
		{
			path: ".",
			want: ".",
		},
	} {
		if got := Normalize(tt.path); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

type bad struct {
	err error
}

func (b *bad) WriteRecord(_ Record) error {
	return b.err
}

var _ RecordWriter = &bad{}

func TestWriteRecordsAndDirs(t *testing.T) {
	// Make sure it fails for the non DedupWriters
	if err := WriteRecordsAndDirs(&bad{}, nil); !errors.Is(err, os.ErrInvalid) {
		t.Errorf("WriteRecordsAndDirs(&bad{}, nil): got %v, want %v", err, os.ErrInvalid)
	}
	var paths = []struct {
		name string
		err  error
	}{
		{name: "a/b/c/d", err: nil},
		{name: "a/b/c/e", err: nil},
		{name: "a/b", err: nil},
	}

	recs := make([]Record, 0)
	for _, p := range paths {
		recs = append(recs, Directory(p.name, 0777))
	}
	var b bytes.Buffer
	w := Newc.Writer(&b)
	if err := WriteRecordsAndDirs(w, recs[:2]); err != nil {
		t.Fatalf("Writing %d records: got %v, want nil", len(recs), err)
	}

	out := "07070100000000000041FF0000000000000000000000000000000000000000000000000000000000000000000000000000000200000000a\x0007070100000000000041FF0000000000000000000000000000000000000000000000000000000000000000000000000000000400000000a/b\x00\x00\x0007070100000000000041FF0000000000000000000000000000000000000000000000000000000000000000000000000000000600000000a/b/c\x0007070100000000000041FF0000000000000000000000000000000000000000000000000000000000000000000000000000000800000000a/b/c/d\x00\x00\x0007070100000000000041FF0000000000000000000000000000000000000000000000000000000000000000000000000000000800000000a/b/c/e\x00\x00\x00"
	if b.String() != out {
		t.Fatalf("%q != %q", b.String(), out)
	}
	if err := WriteRecordsAndDirs(w, recs); !errors.Is(err, os.ErrExist) {
		t.Fatalf("Writing %d records: got %v, want %v", len(recs), err, os.ErrExist)
	}
	// Test a bad write.
	if err := WriteRecordsAndDirs(&bad{err: fs.ErrInvalid}, recs); !errors.Is(err, fs.ErrInvalid) {
		t.Fatalf("Writing %d records: got %v, want %v", len(recs), err, fs.ErrInvalid)
	}

}
