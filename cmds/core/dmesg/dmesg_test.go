// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestDmesg(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	for _, tt := range []struct {
		name      string
		buf       *bytes.Buffer
		bufIn     byte
		clear     bool
		readClear bool
		want      error
	}{
		{
			name:      "both flags set",
			buf:       &bytes.Buffer{},
			clear:     true,
			readClear: true,
			want:      fmt.Errorf("cannot specify both -clear and -read-clear"),
		},
		{
			name:      "both flags unset and buffer has content",
			buf:       &bytes.Buffer{},
			bufIn:     0xEE,
			clear:     false,
			readClear: false,
			want:      fmt.Errorf(""),
		},
		{
			name:      "clear log",
			buf:       &bytes.Buffer{},
			bufIn:     0x41,
			clear:     true,
			readClear: false,
			want:      fmt.Errorf(""),
		},
		{
			name:      "clear log after printing",
			buf:       &bytes.Buffer{},
			bufIn:     0x41,
			clear:     false,
			readClear: true,
			want:      fmt.Errorf(""),
		},
		{
			name:      "clear log after printing and buffer has content",
			buf:       &bytes.Buffer{},
			bufIn:     0xEE,
			clear:     false,
			readClear: true,
			want:      fmt.Errorf(""),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.buf.Write([]byte{tt.bufIn})
			buf.Write([]byte{tt.bufIn})
			if got := dmesg(tt.buf, tt.clear, tt.readClear); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("dmesg() = '%v', want: '%v'", got, tt.want)
				}
			} else {
				if tt.buf.String() != "A" && *clear {
					t.Errorf("System log should be cleared")
				} else if !strings.Contains(tt.buf.String(), buf.String()) && *readClear {
					t.Errorf("System log should contain %s", buf.String())
				} else if tt.buf.String() == "" && (!*clear && !*readClear) {
					t.Errorf("System log should not be cleared")
				}
			}
		})
	}
}
