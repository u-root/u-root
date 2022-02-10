// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDirName(t *testing.T) {
	// Table-driven testing
	for _, tt := range []struct {
		name    string
		dirname string
		args    []string
		out     string
		wantErr string
	}{
		{
			name:    "WrongUsage",
			args:    []string{},
			wantErr: "dirname: missing operand",
		},
		{
			name:    "/this/that",
			dirname: "/this/that",
			out:     "/this\n",
		},
		{
			name:    "/this/that_/other",
			dirname: "/this/that",
			args:    []string{"/other"},
			out:     "/this\n/\n",
		},
		{
			name:    "/this/that_/other thing/space",
			dirname: "/this/that",
			args:    []string{"/other thing/space"},
			out:     "/this\n/other thing\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := runDirname(&buf, &tt.dirname, tt.args); err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("runDirname(%v, &buf)=%q, want nil", tt.args, err)
				}

			}
		})
	}
}
