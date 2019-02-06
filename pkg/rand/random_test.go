// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"context"
	"testing"
	"time"
)

func TestRandomRead(t *testing.T) {
	b := make([]byte, 5)
	n, err := Read(b)
	if err != nil {
		t.Fatalf("got %v, expected nil err", err)
	}
	if n != 5 {
		t.Fatalf("got %d bytes, expected 5 bytes", n)
	}
}

func TestRandomReadContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b := make([]byte, 5)
	n, err := ReadContext(ctx, b)
	if err != nil {
		t.Fatalf("got %v, expected nil err", err)
	}
	if n != 5 {
		t.Fatalf("got %d bytes, expected 5 bytes", n)
	}
}
