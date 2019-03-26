package acpi

import (
	"os"
	"testing"
)

func TestRSDP(t *testing.T) {
	if os.Getuid() != 0 {
		t.Logf("NOT root, skipping")
		t.Skip()
	}
	r, err := GetRSDP()
	if err != nil {
		t.Fatalf("GetRSDP: got %v, want nil", err)
	}
	t.Logf("%v", r)
	s, err := UnMarshalSDT(r)
	if err != nil {
		t.Fatalf("UnMarshalSDT: got %v, want nil", err)
	}
	t.Logf("SDT %v", s)
	tab, err := UnMarshalAll(s)
	if err != nil {
		t.Fatalf("UnMarshalAll: got %v, want nil", err)
	}
	t.Logf("%d entries", len(tab))
	t.Logf("%v", tab)
}
