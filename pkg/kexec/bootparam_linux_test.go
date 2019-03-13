// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package kexec

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	var nd = []byte("it's a nice day")
	bp := NewLinuxBootParams()
	copy(bp.CmdLine[:], nd)
	b, err := bp.Marshal()
	if len(b) != 0x1000 {
		t.Fatalf("Marshaling: got %d bytes, want 0x1000", len(b))
	}
	if err != nil {
		t.Fatalf("Marshaling: got %v, want nil", err)
	}
	ndt := string(b[0x800 : 0x800+len(nd)])
	if ndt != string(nd) {
		t.Fatalf("Marshaling: got %v, want %v", ndt, nd)
	}
}
