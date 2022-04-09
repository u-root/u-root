// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

// Test major and minor numbers greater then 255.
//
// This is supported since Linux 2.6. The major/minor numbers used for this
// test are (1110, 74616). According to "kdev_t.h":
//
//       mkdev(1110, 74616)
//     = mkdev(0x456, 0x12378)
//     = (0x12378 & 0xff) | (0x456 << 8) | ((0x12378 & ~0xff) << 12)
//     = 0x12345678

func TestMknod(t *testing.T) {
	d := t.TempDir()

	// Creating testtable and run tests
	for _, tt := range []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "no flag",
			input: []string{filepath.Join(d, "testdev")},
			want:  "usage:",
		},
		{
			name:  "no path, no flag",
			input: []string{},
			want:  "usage:",
		},
		{
			name:  "p flag with arguments",
			input: []string{filepath.Join(d, "testdev"), "p", "254", "3"},
			want:  "device type p",
		},
		{
			name:  "p flag",
			input: []string{filepath.Join(d, "testdev"), "p"},
			want:  "device type p",
		},
		{
			name:  "b flag with only one argument",
			input: []string{filepath.Join(d, "testdev"), "b", "254"},
			want:  "usage:",
		},
		{
			name:  "b flag with both arguments",
			input: []string{filepath.Join(d, "testdev1"), "b", "1", "254"},
			want:  "mode 61b0:",
		},
		{
			name:  "b flag with large node",
			input: []string{filepath.Join(d, "large_node"), "b", "1110", "74616"},
			want:  "mode 61b0:",
		},
		{
			name:  "c flag without an argument",
			input: []string{filepath.Join(d, "testdev"), "c"},
			want:  "device type c",
		},
		{
			name:  "c flag with both arguments",
			input: []string{filepath.Join(d, "testdev2"), "c", "254", "1"},
			want:  "mode 21b0:",
		},
		{
			name:  "b flag without arguments",
			input: []string{filepath.Join(d, "testdev"), "b"},
			want:  "device type b",
		},
		{
			name:  "invalid flag",
			input: []string{filepath.Join(d, "testdev"), "k"},
			want:  "device type not recognized:",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := mknod(tt.input); got != nil {
				if !strings.Contains(got.Error(), tt.want) {
					t.Errorf("mknod() = '%v', want to contain: '%v'", got, tt.want)
				}
			} else if tt.name == "b flag with large node" {
				// Check the device number.
				var s unix.Stat_t
				if err := unix.Stat(filepath.Join(d, "large_node"), &s); err != nil {
					t.Fatal(err)
				}
				if s.Rdev != 0x12345678 {
					t.Fatalf("expected the device number to be 0x12345678, got %#x", s.Rdev)
				}
			}
		})
	}
}
