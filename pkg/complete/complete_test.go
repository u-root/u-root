// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestSimple tests a basic completer for completion with arrays of strings,
// as might be used for builtin commands.
func TestSimple(t *testing.T) {
	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		tests    = []struct {
			in   string
			outs []string
		}{
			{"hi", []string{"hi"}},
			{"h", hnames},
			{"t", []string{"there"}},
		}
	)

	f := NewStringCompleter(allnames)
	for _, tst := range tests {
		o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("Complete %v: got %v, want nil", tst.in, err)
			continue
		}
		if !reflect.DeepEqual(o, tst.outs) {
			t.Errorf("Complete %v: got %v, want %v", tst.in, o, tst.outs)
		}
	}
}

// TestFile tests the file completer
func TestFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestComplete")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		tests    = []struct {
			in   string
			outs []string
		}{
			{"hi", []string{"hi"}}, // not necessarily intuitive, but the rule is match ONLY one name
			// if that name completes one thing.
			{"h", hnames},
			{"t", []string{"there"}},
		}
	)

	for _, n := range allnames {
		if err := ioutil.WriteFile(filepath.Join(tempDir, n), []byte{}, 0600); err != nil {
			t.Fatal(err)
		}
		t.Logf("Wrote %v", filepath.Join(tempDir, n))
	}
	f := NewFileCompleter(tempDir)
	errCount := 0
	for _, tst := range tests {
		o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("%v: got %v, want nil", tst.in, err)
			errCount++
			continue
		}
		t.Logf("tst %v gets %v", tst, o)
		// potential issue here: we assume FileCompleter, which uses glob, returns
		// sorted order. We'll see if that's an issue later.
		// adjust outs for the path and then check it.
		if len(o) != len(tst.outs) {
			t.Errorf("%v: %v results, want %v", tst, o, tst.outs)
			errCount++
			continue
		}
		for i := range o {
			p := filepath.Join(tempDir, tst.outs[i])
			if o[i] != p {
				t.Errorf("%v: got %v, want %v", tst.in, o, p)
				errCount++
				continue
			}
		}
		t.Logf("tst %v ok", tst)
	}
	t.Logf("%d errors", errCount)
}

// TestMulti tests a multi completer. It creates a multi completer consisting
// of a simple completer and another multicompleter, which in turn has two
// file completers. It also tests the Path completer.
func TestMulti(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestComplete")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		tests    = []struct {
			in   string
			outs []string
		}{
			{"hi", []string{"hi"}},
			{"h", hnames},
			{"t", []string{"there"}},
			{"ahi", []string{"bin/ahi"}},
			{"bh", []string{"sbin/bhi", "sbin/bhil", "sbin/bhit", "sbin/bhow"}},
		}
	)
	for _, p := range []string{"bin", "sbin"} {
		if err := os.MkdirAll(filepath.Join(tempDir, p), 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, n := range allnames {
		if err := ioutil.WriteFile(filepath.Join(tempDir, "bin", "a"+n), []byte{}, 0600); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(filepath.Join(tempDir, "sbin", "b"+n), []byte{}, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Join(tempDir, "bin"), filepath.Join(tempDir, "sbin"))); err != nil {
		t.Fatal(err)
	}
	p, err := NewPathCompleter()
	if err != nil {
		t.Fatal(err)
	}
	// note that since p is a Multi, this also checks nested Multis
	f := NewMultiCompleter(NewStringCompleter(allnames), p)

	for _, tst := range tests {
		o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("Error Complete %v: got %v, want nil", tst.in, err)
			continue
		}
		t.Logf("Complete: tst %v gets %v", tst, o)
		// potential issue here: we assume FileCompleter, which uses glob, returns
		// sorted order. We'll see if that's an issue later.
		// adjust outs for the path and then check it.
		if len(o) != len(tst.outs) {
			t.Errorf("Error Complete %v, wrong len for return: %v results, want %v", tst, o, tst.outs)
			continue
		}
		for i := range o {
			p := tst.outs[i]
			if tst.in[0] == 'a' || tst.in[0] == 'b' {
				p = filepath.Join(tempDir, tst.outs[i])
			}
			t.Logf("\tcheck %v", p)
			if o[i] != p {
				t.Errorf("Error Complete %v, %d'th result mismatches: got %v, want %v", tst.in, i, o[i], p)
			}
		}
		t.Logf("Done check %v: found %v", tst, o)
	}
}

func TestInOut(t *testing.T) {
	var tests = []struct {
		ins   []string
		stack string
	}{
		{[]string{"a", "b", "c", "d"}, "d"},
		{[]string{""}, ""},
		{[]string{}, ""},
	}
	for _, tst := range tests {
		l := NewLine()
		if len(tst.ins) > 0 {
			l.Push(tst.ins...)
		}

		stack := l.Pop()
		if stack != tst.stack {
			t.Errorf("tst %v: got %v, want %v", tst, stack, tst.stack)
		}
	}
}

// TestInOut tests the InOut structures, which we don't know we want.
func TestInOutRW(t *testing.T) {
	var els = []string{"ab", "bc", "de", "fgh"}
	var outs = []string{"ab", "abbc", "abbcde", "abbcdefgh"}

	l := NewLine()
	t.Logf("%v %v %v", els, outs, l)
	for i := range els {
		s := strings.Join(els[:i+1], "")
		l.Write([]byte(s))
		b, err := l.ReadAll()
		if err != nil {
			t.Errorf("ReadAll of %s: got %v, want nil", s, err)
		}
		if string(b) != outs[i] {
			t.Errorf("Read back %s: got %s, want %s", s, string(b), s)
		}
	}
}

// TestLineReader tests Line Readers, and looks for proper read and output behavior.
func TestLineReader(t *testing.T) {
	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		tests    = []struct {
			in    string
			names []string
			out   string
		}{
			{"ther ", []string{"there"}, "there"},
			{"ther", []string{"there"}, "there"},
			{"\n", []string{}, ""},
			{"", []string{}, ""},
			{" ", []string{}, ""},
		}
	)
	for _, tst := range tests {
		r := bytes.NewBufferString(tst.in)
		t.Logf("Test %v", tst)
		cr, cw := io.Pipe()
		f := NewStringCompleter(allnames)

		l := NewLineReader(f, r, cw)
		var out []byte
		go func(o string, r io.Reader) {
			var err error
			out, err = ioutil.ReadAll(r)
			if err != nil {
				t.Errorf("reading console io.Pipe: got %v, want nil", err)
			}
			if string(out) != o {
				t.Errorf("console out: got %v, want %v", o, string(out))
			}
		}(tst.out, cr)

		s, err := l.ReadOne()

		t.Logf("ReadOne returns %v %v", s, err)
		if err != nil && err != io.EOF && err != ErrEOL {
			t.Fatal(err)
		}
		if len(s) != len(tst.names) {
			t.Fatalf("Got %d choices, want 1", len(s))
		}
		if len(s) == 0 {
			continue
		}
		if s[0] != tst.names[0] {
			t.Errorf("Got %v, want %v", s[0], tst.names[0])
		}
	}
}

// TestEnv tests the the environment completer.
func TestEnv(t *testing.T) {
	var tests = []struct {
		pathVal string
		nels    int
		err     error
	}{
		{"", 0, ErrEmptyEnv},
		{"a", 1, nil},
		{"A:B", 2, nil},
	}
	for _, tst := range tests {
		t.Logf("tst %v", tst)
		if err := os.Setenv("PATH", tst.pathVal); err != nil {
			t.Fatal(err)
		}
		e, err := NewPathCompleter()
		if tst.err != err {
			t.Errorf("tst %v: got %v, want %v", tst, err, tst.err)
			continue
		}
		t.Logf("NewPathCompleter returns %v, %v", e, err)
		if tst.nels == 0 && e != nil {
			t.Errorf("tst %v: got %v, want nil", tst, e)
			continue
		}
		if tst.nels == 0 {
			continue
		}
		if e == nil {
			t.Errorf("tst %v: got nil, want MultiCompleter", tst)
			continue
		}
		nels := len(e.(*MultiCompleter).Completers)
		if nels != tst.nels {
			t.Errorf("tst %v: got %d els, want %d", tst, nels, tst.nels)
		}
	}
}
