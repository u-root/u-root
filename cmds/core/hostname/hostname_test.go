// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9 || linux

package main

import (
	"bytes"
	"testing"
)

func TestHostname(t *testing.T) {
	t.Run("wrong number of args", func(t *testing.T) {
		err := run(nil, []string{"hostname", "b", "c"})
		if err == nil {
			t.Error("expected error got nil")
		}
	})

	t.Run("should return hostname", func(t *testing.T) {
		b := &bytes.Buffer{}
		err := run(b, []string{"hostname"})
		if err != nil {
			t.Errorf("expected nil got %v", err)
		}

		if b.String() == "" {
			t.Error("expect non empty hostname")
		}
	})
}
