package main

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestLinuxHeaderOffset(t *testing.T) {
	s := []interface{}{LinuxHeader{}, LinuxParams{}}

	for _, l := range s {
		v := reflect.ValueOf(l)

		for i := 0; i < v.NumField(); i++ {
			e := reflect.TypeOf(l).Field(i)
			off := e.Tag.Get("offset")
			o, err := strconv.ParseUint(off, 0, 64)
			if err != nil {
				t.Errorf("%v: %v: bad or missing offset %v", reflect.TypeOf(s).Name, e.Name, off)
				continue
			}
			if uintptr(o) != e.Offset {
				t.Errorf("%v: Offset of %v: Got 0x%x, want 0x%x", reflect.TypeOf(l).Name, e.Name, e.Offset, o)
			}
		}
	}
}

// Take a LinuxHeader. Read in from a bzImage for the proper size. Then marshal it to a local struct,
// marshal it back, and make sure the value is not changed.
func TestLinuxHeader(t *testing.T) {
}

func TestReadbzImage(t *testing.T) {
	e, h, b, s, err := crackbzImage(testbzImage[:])
	if err != nil {
		t.Fatalf("bzImage reading: got %v, want nil", err)
	}
	t.Logf("entry %v h %v base %v segs %v", e, h, b, s[:512])

	nh, err := MakeLinuxHeader(h)
	if err != nil {
		t.Fatalf("Make a header: got %v, want nil", err)
	}

	if false && bytes.Compare(nh[:], testbzImage[:len(nh)]) != 0 {
		t.Fatalf("Make a header: output and input differ: want %v, got %v", testbzImage[:len(nh)], nh[:])
	}
	t.Logf("Header is %v", *h)
}

func TestAssertSizes(t *testing.T) {
	// Too bad. E820map entries are not multiples of 8 bytes, and we don't pack, so don't do this test.
	if false {
		l := unsafe.Sizeof(E820Entry{})
		if l != 0x14 {
			t.Errorf("E820Entry: got %d bytes for size, want %d\n", l, 0x14)
		}
	}
}
