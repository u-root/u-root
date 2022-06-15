// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"os"
	"reflect"
	"testing"
)

func TestParseFromBytes(t *testing.T) {
	imgBytes, err := os.ReadFile("testdata/Image")
	if err != nil {
		t.Fatal(err)
	}

	wantImage := Image{
		Header: Arm64Header{
			Code0:      0xfa405a4d,
			Code1:      0x141cbfff,
			TextOffset: 0x0,
			ImageSize:  0x940000,
			Flags:      0xa,
			Res2:       0x0,
			Res3:       0x0,
			Res4:       0x0,
			Magic:      0x644D5241,
			Res5:       0x40,
		},
		Data: imgBytes,
	}

	// 1. Success parsing.
	got, err := ParseFromBytes(imgBytes)
	if err != nil {
		t.Errorf("ParseFromBytes(imgBytes) returned error %v, want no error", err)
	}
	if !reflect.DeepEqual(wantImage.Header, got.Header) {
		t.Errorf("got %+v, want %+v", got.Header, wantImage.Header)
	}
}
