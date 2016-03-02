package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

// TestWc do a sequential tests with n args and x input at tab
// of tests below
func TestWc(t *testing.T) {
	var tab = []struct {
		i string
		o string
		s int
		a []string
	}{
		{"simple test count words", "\t4\n", 0, []string{"-w"}}, // don't fail more
		{"lines\nlines\n", "\t2\n", 0, []string{"-l"}},
		{"count chars\n", "\t12\n", 0, []string{"-c"}},
		{"↓→←↑asdf", "\t8\n", 0, []string{"-m"}},
		{"↓€®", "\t3\n", 0, []string{"-r"}},
		{"↑↑↑↑", "\t0\t1\t12\n", 0, nil},
	}

	tmpDir, err := ioutil.TempDir("", "TestWc")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	if false {
		defer os.RemoveAll(tmpDir)
	}

	testwcpath := filepath.Join(tmpDir, "testwc.exe")
	out, err := exec.Command("go", "build", "-o", testwcpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/wc: %v\n%s", testwcpath, err, string(out))
	}

	t.Logf("Built %v for test", testwcpath)
	for _, v := range tab {
		c := exec.Command(testwcpath, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.s {
			t.Errorf("Wc %v < %v > %v: want (exit: %v), got (exit %v)", v.a, v.i, v.o, v.s, s)
			continue
		}

		if err != nil && s != v.s {
			t.Errorf("Wc %v < %v > %v: want nil, got %v", v.a, v.i, v.o, err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Wc %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
		t.Logf("[ok] Wc %v < %v: %v", v.a, v.i, v.o)
	}
}
