// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"os"
	"reflect"
	"testing"
)

// TestSDT tests basic functions, so that we can verify the marshal/unmarshal
// is idempotent.
func TestSDT(t *testing.T) {
	Debug = t.Logf
	if os.Getuid() != 0 {
		t.Logf("NOT root, skipping")
		t.Skip()
	}
	_, r, err := GetRSDP()
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
	// The sdt marshaling, because we need it to, also marshals the tables. Just check
	// the header bytes.
	b = b[:len(sraw.AllData())]
	if !reflect.DeepEqual(sraw.AllData(), b) {
		for i, c := range sraw.AllData() {
			t.Logf("%d: raw %#02x b %#02x", i, c, b[i])
		}
		t.Fatalf("TestSDT: input and output []byte differ: in %q, out %q: want same", sraw, b)
	}
}

// TestNewSDT tests to ensure that NewSDT returns a correct empty SDT and marshals
// to a HeaderLength byte array.
func TestNewSDT(t *testing.T) {
	s, err := NewSDT()
	if err != nil {
		t.Fatal(err)
	}
	if len(s.data) != HeaderLength {
		t.Fatalf("NewSDT: got size %d, want %d", len(s.data), HeaderLength)
	}
}
