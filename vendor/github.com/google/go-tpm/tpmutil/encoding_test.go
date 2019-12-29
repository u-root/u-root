// Copyright (c) 2018, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpmutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type invalidPacked struct {
	A []int
	B uint32
}

func testEncodingInvalidSlices(t *testing.T, f func(io.Writer, interface{}) error) {
	d := ioutil.Discard

	// The packedSize function doesn't handle slices to anything other than bytes.
	var invalid []int
	if err := f(d, invalid); err == nil {
		t.Fatal("The packing function incorrectly succeeds for a slice of integers")
	}
	if err := f(d, &invalid); err == nil {
		t.Fatal("The packing function incorrectly succeeds for a pointer to a slice of integers")
	}

	invalid2 := invalidPacked{
		A: make([]int, 10),
		B: 137,
	}
	if err := f(d, invalid2); err == nil {
		t.Fatal("The packing function incorrectly succeeds for a struct that contains an integer slice")
	}
	if err := f(d, &invalid2); err == nil {
		t.Fatal("The packing function incorrectly succeeds for a pointer to a struct that contains an integer slice")
	}

	if err := f(d, d); err == nil {
		t.Fatal("The packing function incorrectly succeeds for a non-packable value")
	}
}

func TestEncodingPackTypeInvalid(t *testing.T) {
	f := func(w io.Writer, i interface{}) error {
		return packType(w, i)
	}

	testEncodingInvalidSlices(t, f)
}

type simplePacked struct {
	A uint32
	B uint32
}

type nestedPacked struct {
	SP simplePacked
	C  uint32
}

type nestedSlice struct {
	A uint32
	S U32Bytes
}

func TestEncodingPackType(t *testing.T) {
	buf := make([]byte, 10)
	inputs := []interface{}{
		uint32(3),
		buf,
		&buf,
		simplePacked{137, 138},
		nestedPacked{simplePacked{137, 138}, 139},
		nestedSlice{137, buf},
		[]byte(nil),
		RawBytes(buf),
	}
	for _, i := range inputs {
		if err := packType(ioutil.Discard, i); err != nil {
			t.Errorf("packType(%#v): %v", i, err)
		}
	}
}

func TestEncodingPackTypeWriteFail(t *testing.T) {
	u32WithOneByte := U32Bytes([]byte{1})
	u32Empty := U32Bytes([]byte(nil))

	tests := []struct {
		limit int
		in    interface{}
	}{
		{4, &u32WithOneByte},
		{3, &u32Empty},
	}
	for _, tt := range tests {
		if err := packType(&limitedDiscard{tt.limit}, tt.in); err == nil {
			t.Errorf("packType(%#v) with write size limit %d returned nil, want error", tt.in, tt.limit)
		}
	}
}

// limitedDiscard is an implementation of io.Writer that accepts a given number
// of bytes before returning errors.
type limitedDiscard struct {
	remaining int
}

// Write writes p to the limitedDiscard instance.
func (l *limitedDiscard) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > l.remaining {
		n = l.remaining
		err = io.EOF
	}

	l.remaining -= n
	return
}

func TestEncodingCommandHeaderInvalidBody(t *testing.T) {
	var invalid []int
	ch := commandHeader{1, 0, 2}
	_, err := packWithHeader(ch, invalid)
	if err == nil {
		t.Fatal("packWithHeader incorrectly packed a body that with an invalid int slice member")
	}
}

func TestEncodingInvalidPack(t *testing.T) {
	var invalid []int
	ch := commandHeader{1, 0, 2}
	_, err := packWithHeader(ch, invalid)
	if err == nil {
		t.Fatal("packWithHeader incorrectly packed a body that with an invalid int slice member")
	}

	_, err = Pack(invalid)
	if err == nil {
		t.Fatal("pack incorrectly packed a slice of int")
	}
}

func TestEncodingCommandHeaderEncoding(t *testing.T) {
	ch := commandHeader{1, 0, 2}
	var c uint32 = 137
	in := c

	b, err := packWithHeader(ch, in)
	if err != nil {
		t.Fatal("Couldn't pack the bytes:", err)
	}

	var hdr commandHeader
	var size uint32
	if _, err := Unpack(b, &hdr, &size); err != nil {
		t.Fatal("Couldn't unpack the packed bytes")
	}

	if size != 137 {
		t.Fatal("Got the wrong size back")
	}
}

func TestEncodingInvalidUnpack(t *testing.T) {
	var i *uint32
	i = nil
	// The value ui is a serialization of uint32(0).
	ui := []byte{0, 0, 0, 0}
	uiBuf := bytes.NewBuffer(ui)
	if err := UnpackBuf(uiBuf, i); err == nil {
		t.Fatal("UnpackBuf incorrectly deserialized into a nil pointer")
	}

	var ii uint32
	if err := UnpackBuf(uiBuf, ii); err == nil {
		t.Fatal("UnpackBuf incorrectly deserialized into a non pointer")
	}

	var b U32Bytes
	var empty []byte
	emptyBuf := bytes.NewBuffer(empty)
	if err := UnpackBuf(emptyBuf, &b); err == nil {
		t.Fatal("UnpackBuf incorrectly deserialized an empty byte array into U32Bytes")
	}

	// Try to deserialize a byte array that has a length but not enough bytes.
	// The slice ui represents uint32(1), which is the length of an empty byte array.
	ui2 := []byte{0, 0, 0, 1}
	uiBuf2 := bytes.NewBuffer(ui2)
	if err := UnpackBuf(uiBuf2, &b); err == nil {
		t.Fatal("UnpackBuf incorrectly deserialized a byte array that didn't have enough bytes available")
	}

	var iii []int
	ui3 := []byte{0, 0, 0, 1}
	uiBuf3 := bytes.NewBuffer(ui3)
	if err := UnpackBuf(uiBuf3, &iii); err == nil {
		t.Fatal("UnpackBuf incorrectly deserialized into a slice of ints (only byte slices are supported)")
	}

}

func TestSelfMarshaler(t *testing.T) {
	var empty16 U16Bytes
	var empty32 U32Bytes
	subTests := []struct {
		encoded []byte
		decoded interface{}
	}{
		{[]byte{0, 0}, &empty16},
		{[]byte{0, 1, 137}, &empty16},
		{[]byte{0, 0, 0, 0}, &empty32},
		{[]byte{0, 0, 0, 1, 137}, &empty32},
	}
	for _, st := range subTests {
		t.Logf("Attempting to Marshal/Unmarshal %#v into %T", st.encoded, st.decoded)
		buffer := bytes.NewBuffer(st.encoded)
		if err := UnpackBuf(buffer, st.decoded); err != nil {
			t.Fatalf("UnpackBuf failed: %v", err)
		}
		packed, err := Pack(st.decoded)
		if err != nil {
			t.Fatalf("Pack failed: %v", err)
		}
		if !bytes.Equal(packed, st.encoded) {
			t.Fatalf("Pack failed: got %#v, want: %#v", packed, st.encoded)
		}
	}
}

func TestEncodingUnpack(t *testing.T) {
	// Deserialize the empty byte array.
	var b U32Bytes
	// The slice ui represents uint32(0), which is the length of an empty byte array.
	ui := []byte{0, 0, 0, 0}
	uiBuf := bytes.NewBuffer(ui)
	if err := UnpackBuf(uiBuf, &b); err != nil {
		t.Fatal("UnpackBuf failed to unpack the empty byte array")
	}

	// A byte slice of length 1 with a single entry: b[0] == 137
	ui2 := []byte{0, 0, 0, 1, 137}
	uiBuf2 := bytes.NewBuffer(ui2)
	if err := UnpackBuf(uiBuf2, &b); err != nil {
		t.Fatal("UnpackBuf failed to unpack a byte array with a single value in it")
	}

	if !bytes.Equal([]byte(b), []byte{137}) {
		t.Fatal("UnpackBuf unpacked a small byte array incorrectly")
	}

	sp := simplePacked{137, 138}
	bsp, err := Pack(sp)
	if err != nil {
		t.Fatal("Couldn't pack a simple struct:", err)
	}
	var sp2 simplePacked
	if _, err := Unpack(bsp, &sp2); err != nil {
		t.Fatal("Couldn't unpack a simple struct:", err)
	}

	if sp.A != sp2.A || sp.B != sp2.B {
		t.Fatal("Unpacked simple struct didn't match the original")
	}

	// Try unpacking a version that's missing a byte at the end.
	if _, err := Unpack(bsp[:len(bsp)-1], &sp2); err == nil {
		t.Fatal("unpack incorrectly unpacked from a byte array that didn't have enough values")
	}

	np := nestedPacked{sp, 139}
	bnp, err := Pack(np)
	if err != nil {
		t.Fatal("Couldn't pack a nested struct")
	}
	var np2 nestedPacked
	if _, err := Unpack(bnp, &np2); err != nil {
		t.Fatal("Couldn't unpack a nested struct:", err)
	}
	if np.SP.A != np2.SP.A || np.SP.B != np2.SP.B || np.C != np2.C {
		t.Fatal("Unpacked nested struct didn't match the original")
	}

	ns := nestedSlice{137, b}
	bns, err := Pack(&ns)
	if err != nil {
		t.Fatal("Couldn't pack a struct with a nested byte slice:", err)
	}
	var ns2 nestedSlice
	if _, err := Unpack(bns, &ns2); err != nil {
		t.Fatal("Couldn't unpacked a struct with a nested slice:", err)
	}
	if ns.A != ns2.A || !bytes.Equal(ns.S, ns2.S) {
		t.Logf("orginal = %+v", ns)
		t.Logf("decoded = %+v", ns2)
		t.Fatal("Unpacked struct with nested slice didn't match the original")
	}

	var hs []Handle
	if _, err := Unpack([]byte{0, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, &hs); err != nil {
		t.Fatal("Couldn't unpack a list of Handles:", err)
	}
	if want := []Handle{0x01020304, 0x05060708, 0x090a0b0c}; !reflect.DeepEqual(want, hs) {
		t.Fatalf("Unpacking []Handle: got %v, want %v", hs, want)
	}
}

func TestPartialUnpack(t *testing.T) {
	u1, u2 := uint32(1), uint32(2)
	buf, err := Pack(u1, u2)
	if err != nil {
		t.Fatalf("packing uint32 value: %v", err)
	}

	var gu1, gu2 uint32
	read1, err := Unpack(buf, &gu1)
	if err != nil {
		t.Fatalf("unpacking first uint32 value: %v", err)
	}
	if gu1 != u1 {
		t.Errorf("first unpacked value: got %d, want %d", gu1, u1)
	}
	read2, err := Unpack(buf[read1:], &gu2)
	if err != nil {
		t.Fatalf("unpacking second uint32 value: %v", err)
	}
	if gu2 != u2 {
		t.Errorf("second unpacked value: got %d, want %d", gu2, u2)
	}

	if read1+read2 != len(buf) {
		t.Errorf("sum of bytes read doesn't ad up to total packed size: got %d+%d=%d, want %d", read1, read2, read1+read2, len(buf))
	}
}

func TestUnpackHandlesArea(t *testing.T) {
	buf := []byte{
		0, 2,
		0, 0, 0, 1,
		0, 0, 5, 57,
	}
	var out []Handle

	if _, err := Unpack(buf, &out); err != nil {
		t.Fatalf("Unpack(%v, %T) failed: %v", buf, &out, err)
	}
	if want := []Handle{1, 1337}; !reflect.DeepEqual(out, want) {
		t.Errorf("Unpack(%v, %T): %T = %v, want %v", buf, &out, out, out, want)
	}
}
