// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"context"
	"testing"
)

func TestFallback(t *testing.T) {
	r := &getrandomReader{backup: true}
	b := make([]byte, 5)
	n, err := r.ReadContext(context.Background(), b)
	if err != nil {
		t.Fatalf("got %v, expected nil err", err)
	}
	if n != 5 {
		t.Fatalf("got %d bytes, expected 5 bytes", n)
	}
}
