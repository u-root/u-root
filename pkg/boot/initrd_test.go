// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/uio"
)

func TestConcatenationSize(t *testing.T) {
	a := strings.NewReader("yay")
	b := bytes.NewReader(make([]byte, 777))
	reader := CatInitrds(a, b)
	res, _ := uio.ReadAll(reader)
	size := len(res)
	if size != 1536 {
		t.Errorf("want 1536 bytes, got %v", size)
	}
}

func TestConcatenationBytes(t *testing.T) {
	a := strings.NewReader("foo")
	b := strings.NewReader("bar")
	reader := CatInitrds(a, b)
	res, _ := uio.ReadAll(reader)
	if res[0] != 'f' {
		t.Errorf("byte 0 is not f")
	}
	if res[513] != 'a' {
		t.Errorf("byte 513 is not a")
	}
}
