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

func TestSimple(t *testing.T) {
	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		tests    = []struct {
			c string
			o []string
		}{
			{"hi", hinames},
			{"h", hnames},
			{"t", []string{"there"}},
		}
	)

	f := NewStringCompleter(allnames)
	debug = t.Logf
	for _, tst := range tests {
		o, err := f.Complete(tst.c)
		if err != nil {
			t.Errorf("Complete %v: got %v, want nil", tst.c, err)
			continue
		}
		if !reflect.DeepEqual(o, tst.o) {
			t.Errorf("Complete %v: got %v, want %v", tst.c, o, tst.o)
		}
	}
}

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
			c string
			o []string
		}{
			{"hi", hinames},
			{"h", hnames},
			{"t", []string{"there"}},
		}
	)

	for _, n := range allnames {
		if err := ioutil.WriteFile(filepath.Join(tempDir, n), []byte{}, 0600); err != nil {
			t.Fatal(err)
		}
	}
	f := NewFileCompleter(tempDir)
	debug = t.Logf
	for _, tst := range tests {
		o, err := f.Complete(tst.c)
		if err != nil {
			t.Errorf("Complete %v: got %v, want nil", tst.c, err)
			continue
		}
		if !reflect.DeepEqual(o, tst.o) {
			t.Errorf("Complete %v: got %v, want %v", tst.c, o, tst.o)
		}
	}
}

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
			c string
			o []string
		}{
			{"hi", hinames},
			{"h", hnames},
			{"t", []string{"there"}},
			{"ahi", []string{"ahi", "ahil", "ahit"}},
			{"bh", []string{"bhi", "bhil", "bhit", "bhow"}},
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
	debug = t.Logf
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
		o, err := f.Complete(tst.c)
		if err != nil {
			t.Errorf("Complete %v: got %v, want nil", tst.c, err)
			continue
		}
		if !reflect.DeepEqual(o, tst.o) {
			t.Errorf("Complete %v: got %v, want %v", tst.c, o, tst.o)
		}
		t.Logf("Check %v: found %v", tst, o)
	}
}

func TestInOut(t *testing.T) {
	var tests = []struct {
		i []string
		o string
	}{
		{[]string{"a", "b", "c", "d"}, "d"},
		{[]string{""}, ""},
		{[]string{}, ""},
	}
	for _, tst := range tests {
		l := NewLine()
		if len(tst.i) > 0 {
			l.Push(tst.i[0], tst.i[1:]...)
		}

		o := l.Pop()
		if o != tst.o {
			t.Errorf("tst %v: got %v, want %v", tst, o, tst.o)
		}
	}
}

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

func TestLineReader(t*testing.T) {
	var (
		hinames  = []string{"hi", "hil", "hit"}
		hnames   = append(hinames, "how")
		allnames = append(hnames, "there")
		r = bytes.NewBufferString("ther\t")
	)
	cr, cw := io.Pipe()
	f := NewStringCompleter(allnames)
	debug = t.Logf

	l := NewLineReader(f, r, cw)
	var out []byte
	go func() {
		var err error
		out, err = ioutil.ReadAll(cr)
		if err != nil {
			t.Errorf("reading console io.Pipe: got %v, want nil", err)
		}
		if string(out) != "there" {
			t.Errorf("console out: got %v, want ther", string(out))
		}
	}()

	s, err := l.ReadOne()
	
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 1 {
		t.Fatalf("Got %d choices, want 1", len(s))
	}
	if s[0] != "there" {
		t.Errorf("Got %v choices, want there", s[0])
	}
}
