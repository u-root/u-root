// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"bytes"
	"testing"
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		name      string
		all       bool
		kernel    bool
		node      bool
		release   bool
		version   bool
		machine   bool
		processor bool
		expected  params
	}{
		{
			name: "flag: -a",
			all:  true,
			expected: params{
				kernel:  true,
				node:    true,
				release: true,
				version: true,
				machine: true,
			},
		},
		{
			name: "no flags",
			expected: params{
				kernel: true,
			},
		},
		{
			name:    "flags: -s -n -r -v -m",
			kernel:  true,
			node:    true,
			release: true,
			version: true,
			machine: true,
			expected: params{
				kernel:  true,
				node:    true,
				release: true,
				version: true,
				machine: true,
			},
		},
		{
			name:      "flags: -p",
			processor: true,
			expected: params{
				machine: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parseParams(tt.all, tt.kernel, tt.node, tt.release, tt.version, tt.machine, tt.processor)
			if p != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, p)
			}
		})
	}
}

func TestRun(t *testing.T) {
	stdout := &bytes.Buffer{}
	p := params{kernel: true, node: true, release: true, version: true, machine: true}
	// it's possbile to test handleFlags() directly, but unix.Utsname is
	// platform dependent
	err := run(stdout, p)
	if err != nil || stdout.String() == "" {
		t.Error("expected no error and non empty output")
	}
}
