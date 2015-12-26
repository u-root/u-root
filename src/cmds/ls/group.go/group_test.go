package group

import (
	"runtime"
	"testing"
)

func skip(t *testing.T) bool {
	if !implemented {
		t.Logf("user: not implemented; skipping tests")
		return true
	}

	switch runtime.GOOS {
	case "linux", "freebsd", "darwin":
		return false
	}

	t.Logf("user: Lookup not implemented on %s; skipping test", runtime.GOOS)
	return true
}

func TestCurrent(t *testing.T) {
	if skip(t) {
		return
	}
	g, err := Current()
	if err != nil {
		t.Fatalf("Current: %v", err)
	}
	if g.Name == "" {
		t.Fatalf("didn't get a group name")
	}
}

func compare(t *testing.T, want, got *Group) {
	if want.Gid != got.Gid {
		t.Errorf("got Gid=%q; want %q", got.Gid, want.Gid)
	}
	if want.Name != got.Name {
		t.Errorf("got Name=%q; want %q", got.Name, want.Name)
	}
	if len(want.Members) != len(got.Members) {
		t.Errorf("got %d Members; want %d", len(got.Members), len(want.Members))
	}
	for i := range want.Members {
		if want.Members[i] != got.Members[i] {
			t.Errorf(
				"got Members[%d]=%q; want %q", got.Members[i], want.Members[i])
		}
	}
}

func TestLookup(t *testing.T) {
	if skip(t) {
		return
	}
	want, err := Current()
	if err != nil {
		t.Fatalf("Current: %v", err)
	}
	got, err := Lookup(want.Name)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	compare(t, want, got)
}

func TestLookupId(t *testing.T) {
	if skip(t) {
		return
	}

	want, err := Current()
	if err != nil {
		t.Fatalf("Current: %v", err)
	}
	got, err := LookupId(want.Gid)
	if err != nil {
		t.Fatalf("LookupId: %v", err)
	}
	compare(t, want, got)
}
