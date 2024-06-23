// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race && !tinygo && linux && (amd64 || riscv64 || arm64)
// +build !race
// +build !tinygo
// +build linux
// +build amd64 riscv64 arm64

package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		err  error
		p    params
		args []string
	}{
		{
			args: []string{"echo", "hello", "u-root"},
			p:    params{},
		},
		{
			args: []string{"echo", "hello", "u-root"},
			p: params{
				output: filepath.Join(tmp, "file-test-one-1"),
			},
		},
		{
			p:   params{},
			err: errUsage,
		},
	}

	for _, test := range tests {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		err := run(nil, stdout, stderr, test.p, test.args...)
		if !errors.Is(err, test.err) {
			t.Fatalf("expected %v, got %v", test.err, err)
		}
	}

}
