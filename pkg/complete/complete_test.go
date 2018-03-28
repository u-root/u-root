package complete

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
