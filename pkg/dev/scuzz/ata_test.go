// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

import (
	"testing"
)

func TestAtaString(t *testing.T) {
	const want = "Copyright 2019 the u"
	var ts = ataString{
		'o'<<8 | 'C',
		'y'<<8 | 'p',
		'i'<<8 | 'r',
		'h'<<8 | 'g',
		' '<<8 | 't',
		'0'<<8 | '2',
		'9'<<8 | '1',
		't'<<8 | ' ',
		'e'<<8 | 'h',
		'u'<<8 | ' ',
	}

	got := ts.String()
	if got != want {
		t.Fatalf("Got %v, want %v", got, want)
	}
}

func TestMustLBA(t *testing.T) {
	var data dataBlock

	w, err := data.toWordBlock()
	if err != nil {
		t.Fatalf("toWordBlock: got %v, want nil", err)
	}

	if err := w.mustLBA(); err == nil {
		t.Errorf("bad mustLBA: got nil, want x")
	}

	data[0], data[49*2], data[83*2], data[86*2] = 0x80, 0x2, 0x40, 0x4

	w, err = data.toWordBlock()
	if err != nil {
		t.Fatalf("toWordBlock: got %v, want nil", err)
	}

	if err := w.mustLBA(); err != nil {
		t.Errorf("good mustLBA: got %v, want nil", err)
	}
}
