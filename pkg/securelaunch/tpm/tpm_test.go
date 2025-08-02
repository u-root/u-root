// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tpm

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestMarshalPcrEvent(t *testing.T) {
	pcr := uint32(1)

	// 94023c06c5c6bb7b408d1b08b7e32371d9302436f31c8ee73eaae9f08ded7da7
	h := []byte{
		148, 2, 60, 6, 197, 198, 187, 123, 64, 141, 27, 8, 183, 227, 35,
		113, 217, 48, 36, 54, 243, 28, 142, 231, 62, 170, 233, 240, 141,
		237, 125, 167,
	}
	eventDesc := []byte("Test description")

	want := []byte{
		1, 0, 0, 0, 2, 5, 0, 0, 1, 0, 0, 0, 11, 0, 148, 2, 60, 6, 197,
		198, 187, 123, 64, 141, 27, 8, 183, 227, 35, 113, 217, 48, 36, 54, 243,
		28, 142, 231, 62, 170, 233, 240, 141, 237, 125, 167, 16, 0, 0, 0, 84,
		101, 115, 116, 32, 100, 101, 115, 99, 114, 105, 112, 116, 105, 111, 110,
	}

	got, err := marshalPcrEvent(pcr, h, eventDesc)
	if err != nil {
		log.Fatalf("marshalPcrEvent() = %v, not nil", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("marshalPcrEvent() = %v, want %v", got, want)
	}
}

func TestHashReader(t *testing.T) {
	testString := "test string"
	want := []byte{
		213, 87, 156, 70, 223, 204, 127, 24, 32, 112, 19, 230, 91,
		68, 228, 203, 78, 44, 34, 152, 244, 172, 69, 123, 168, 248, 39, 67, 243,
		30, 147, 11,
	}
	got := HashReader(strings.NewReader(testString))

	if !bytes.Equal(got, want) {
		t.Errorf("hashReader() = %v, want %v", got, want)
	}
}
