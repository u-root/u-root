// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

func TestDmesg(t *testing.T) {
	guest.SkipIfNotInVM(t)

	for _, tt := range []struct {
		name      string
		buf       *bytes.Buffer
		bufIn     byte
		clear     bool
		readClear bool
		err       error
	}{
		{
			name:      "both flags set",
			buf:       &bytes.Buffer{},
			clear:     true,
			readClear: true,
			err:       os.ErrInvalid,
		},
		{
			name:      "both flags unset and buffer has content",
			buf:       &bytes.Buffer{},
			bufIn:     0xEE,
			clear:     false,
			readClear: false,
			err:       nil,
		},
		{
			name:      "clear log",
			buf:       &bytes.Buffer{},
			bufIn:     0x41,
			clear:     true,
			readClear: false,
			err:       nil,
		},
		{
			name:      "clear log after printing",
			buf:       &bytes.Buffer{},
			bufIn:     0x41,
			clear:     false,
			readClear: true,
			err:       nil,
		},
		{
			name:      "clear log after printing and buffer has content",
			buf:       &bytes.Buffer{},
			bufIn:     0xEE,
			clear:     false,
			readClear: true,
			err:       nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.buf.Write([]byte{tt.bufIn})
			buf.Write([]byte{tt.bufIn})
			if err := dmesg(tt.buf, tt.clear, tt.readClear); err != nil {
				// Some container environments return uid 0,
				// but they are lying. If the error is ErrPermission,
				// just return.
				if errors.Is(err, os.ErrPermission) {
					t.Skipf("Ignore test due to CI issue:%v", err)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("dmesg() = '%v', want: '%v'", err, tt.err)
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
