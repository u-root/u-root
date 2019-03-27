package acpi

import (
	"os"
	"reflect"
	"testing"
)

func TestSDT(t *testing.T) {
	if os.Getuid() != 0 {
		t.Logf("NOT root, skipping")
		t.Skip()
	}
	r, err := GetRSDP()
	if err != nil {
		t.Fatalf("TestSDT GetRSDP: got %v, want nil", err)
	}
	t.Logf("%v", r)
	s, err := UnMarshalSDT(r)
	if err != nil {
		t.Fatalf("TestSDT: got %v, want nil", err)
	}

	sraw, err := ReadRaw(r.Base())
	if err != nil {
		t.Fatalf("TestSDT: readraw got %v, want nil", err)
	}

	b, err := s.Marshal()
	if err != nil {
		t.Fatalf("Marshaling SDT: got %v, want nil", err)
	}

	if !reflect.DeepEqual(sraw, b) {
		t.Fatalf("TestSDT: input and output []byte differ: in %v, out %v: want same", sraw, b)
	}
}
