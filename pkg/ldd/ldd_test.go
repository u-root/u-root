package uroot

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestLdd tests Ldd against /bin/date.
// This is just about guaranteed to have
// some output on most linux systems.
func TestLdd(t *testing.T) {
	n, err := Ldd([]string{"/bin/date"})
	if err != nil {
		t.Fatalf("Ldd on /bin/date: want nil, got %v", err)
	}
	t.Logf("TestLdd: /bin/date has deps of")
	for i := range n {
		t.Logf("\t%v", n[i])
	}
}

// TestLddList tests that the LddList is the
// same as the info returned by Ldd
func TestLddList(t *testing.T) {
	var libMap = make(map[string]bool)
	n, err := Ldd([]string{"/bin/date"})
	if err != nil {
		t.Fatalf("Ldd on /bin/date: want nil, got %v", err)
	}
	l, err := LddList([]string{"/bin/date"})
	if err != nil {
		t.Fatalf("LddList on /bin/date: want nil, got %v", err)
	}
	if len(n) != len(l) {
		t.Fatalf("Len of Ldd(%v) and LddList(%v): want same, got different", len(n), len(l))
	}
	for i := range n {
		libMap[n[i].FullName] = true
	}
	for i := range n {
		if ! libMap[l[i]] {
			t.Errorf("%v was in LddList but not in Ldd", l[i])
		}
	}
}

// This could have been a great test, if ld.so actually followed ITS OWN DOCS
// and used LD_LIBRARY_PATH. It doesn't.
func testLddBadSo(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "ldd")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	if err := os.Setenv("LD_LIBRARY_PATH", tempDir); err != nil {
		t.Fatalf("Setting LDD_LIBRARY_PATH to %v: want nil, got %v", tempDir, err)
	}
	if _, err := Ldd([]string{"/bin/date"}); err == nil {
		t.Fatalf("Ldd on /bin/date: want err, got nil")
	}
	t.Logf("Err on bad dir is %v", err)

}
