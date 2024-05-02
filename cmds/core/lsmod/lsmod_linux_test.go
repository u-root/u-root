// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("file don't exists", func(t *testing.T) {
		err := run(nil, "filenotexists")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("lsmod", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		err := run(stdout, "./testdata/modules.txt")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !strings.Contains(stdout.String(), "nf_defrag_ipv6") {
			t.Errorf("expected to have nf_defrag_ipv6, got %v", stdout.String())
		}
	})
}
