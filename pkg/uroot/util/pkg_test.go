// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestTODO(t *testing.T) {
	b := &bytes.Buffer{}
	f := func() {
		fmt.Fprintf(b, "hi %s", os.Args[0])
	}

	f = Usage(f, "there")
	f()
	if b.String() != "hi there" {
		t.Errorf("f(): Got %q, want %q", b.String(), "hi there")
	}
}
