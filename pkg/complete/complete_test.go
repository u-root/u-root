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
			x    string
			outs []string
		}{
			{"hi", "hi", []string{}},
			{"h", "", hnames},
			{"t", "", []string{"there"}},
		}
	)

	f := NewStringCompleter(allnames)
	for _, tst := range tests {
		x, o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("Complete %v: got %v, want nil", tst.in, err)
			continue
		}
		if x != tst.x {
			t.Errorf("Complete %v: got %v, want %v", tst.in, x, tst.x)
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
			in string
			x  string
			g  []string
		}{
			{"hi", "hi", hinames[1:]},
			{"h", "", hnames},
			{"t", "", []string{"there"}},
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
		x, o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("%v: got %v, want nil", tst.in, err)
			errCount++
			continue
		}
		t.Logf("tst %v gets %v %v", tst, x, o)
		// potential issue here: we assume FileCompleter, which uses glob, returns
		// sorted order. We'll see if that's an issue later.
		// adjust outs for the path and then check it.
		if len(o) != len(tst.g) {
			t.Errorf("%v: %v results, want %v", tst, o, tst.g)
			errCount++
			continue
		}
		if tst.x != x && x != filepath.Join(tempDir, tst.x) {
			t.Errorf("%v: got %v, want %v", tst.in, x, filepath.Join(tempDir, tst.x))
			errCount++
		}
		for i := range o {
			p := filepath.Join(tempDir, tst.g[i])
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
			x    string
			outs []string
		}{
			{"hi", "hi", []string{}},
			{"h", "", hnames},
			{"t", "", []string{"there"}},
			{"ahi", "bin/ahi", []string{"bin/ahil", "bin/ahit"}},
			{"bh", "", []string{"sbin/bhi", "sbin/bhil", "sbin/bhit", "sbin/bhow"}},
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
		x, o, err := f.Complete(tst.in)
		if err != nil {
			t.Errorf("Error Complete %v: got %v, want nil", tst.in, err)
			continue
		}
		t.Logf("Complete: tst %v gets %v", tst, o)
		if tst.x != x && x != filepath.Join(tempDir, tst.x) {
			t.Errorf("ERROR %v: got %v, want %v", tst.in, x, tst.x)
		}

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
			in      string
			names   []string
			x       string
			choices []string
			out     string
		}{
			{"there\t", []string{"there"}, "there", []string{}, "there"},
			{"there", []string{"there"}, "", []string{}, "there"},
			{"\n", []string{}, "", []string{}, ""},
			{"", []string{}, "", []string{}, ""},
			{" ", []string{}, "", []string{}, ""},
		}
	)
	Debug = t.Logf
	for i, tst := range tests[:1] {
		r := bytes.NewBufferString(tst.in)
		t.Logf("%d: Test %v", i, tst)
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

		err := l.ReadLine()

		x := l.Exact
		s := l.Candidates
		t.Logf("ReadLine returns %v %v %v", x, s, err)
		if err != nil && err != io.EOF && err != ErrEOL {
			t.Errorf("Test %d: got %v, want nil", i, err)
			continue
		}
		if len(s) != len(tst.choices) {
			t.Errorf("Test %d: Got %d choices, want %d", i, len(s), len(tst.choices))
			continue
		}
		if len(s) == 0 {
			continue
		}
		if x != tst.x {
			t.Errorf("Test %d: Got %v, want %v", i, x, tst.x)
			continue
		}
		t.Logf("%d passes", i)
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
