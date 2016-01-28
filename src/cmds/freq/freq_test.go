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

func TestFreq(t *testing.T) {
	var tab = []struct {
		i string
		o string
		s int
		a []string
	}{
		{
			"u-root",
			"-        1\no        2\nr        1\nt        1\nu        1\n",
			0,
			[]string{"-c"},
		},
		{
			"ron",
			"6e        1\n6f        1\n72        1\n",
			0,
			[]string{"-x"},
		},
	}

	tmpDir, err := ioutil.TempDir("", "TestFreq")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	if false {
		defer os.RemoveAll(tmpDir)
	}

	testfreqpath := filepath.Join(tmpDir, "testfreq.exe")
	out, err := exec.Command("go", "build", "-o", testfreqpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/freq: %v\n%s", testfreqpath, err, string(out))
	}

	t.Logf("Built %v for test", testfreqpath)
	for _, v := range tab {
		c := exec.Command(testfreqpath, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.s {
			t.Errorf("Freq %v < %v > %v: want (exit: %v), got (exit %v)", v.a, v.i, v.o, v.s, s)
			continue
		}

		if err != nil && s != v.s {
			t.Errorf("Freq %v < %v > %v: want nil, got %v", v.a, v.i, v.o, err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Freq %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
		t.Logf("[ok] Freq %v < %v: %v", v.a, v.i, v.o)
	}
}
