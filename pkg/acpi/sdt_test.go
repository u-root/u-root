package acpi

import (
	"os"
	"reflect"
	"testing"
)

func TestSDT(t *testing.T) {
	Debug = t.Logf
	if os.Getuid() != 0 {
		t.Logf("NOT root, skipping")
		t.Skip()
	}
	r, err := GetRSDP()
	if err != nil {
		t.Fatalf("TestSDT GetRSDP: got %v, want nil", err)
	}
	t.Logf("%q", r)
	s, err := UnMarshalSDT(r)
	if err != nil {
		t.Fatalf("TestSDT: got %q, want nil", err)
	}
	t.Logf("%q::%s", s, ShowTable(s))
	sraw, err := ReadRaw(r.Base())
	if err != nil {
		t.Fatalf("TestSDT: readraw got %q, want nil", err)
	}
	t.Logf("%q", sraw)
	b, err := s.Marshal()
	if err != nil {
		t.Fatalf("Marshaling SDT: got %q, want nil", err)
	}
	t.Logf("%q", b)
	if !reflect.DeepEqual(sraw.AllData(), b) {
		for i, c := range sraw.AllData() {
			t.Logf("%d: raw %#02x b %#02x", i, c, b[i])
		}
		t.Fatalf("TestSDT: input and output []byte differ: in %q, out %q: want same", sraw, b)
	}
}
