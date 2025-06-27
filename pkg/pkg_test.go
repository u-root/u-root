// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	"os"
	"testing"
)

func TestOsContext(t *testing.T) {
	p := OsContext(context.Background())

	if wd, err := p.Getwd(); wd == "" {
		t.Error("got empty working directory")
	} else if err != nil {
		t.Errorf("error from OsContext.Getwd: %v", err)
	} else if oswd, oserr := os.Getwd(); oserr != nil {
		t.Errorf("error from os.GetWd: %v", oserr)
	} else if oswd != wd {
		t.Errorf("got different working directory from OsContext.Getwd %#v and os.Getwd %#v", wd, oswd)
	}
}
