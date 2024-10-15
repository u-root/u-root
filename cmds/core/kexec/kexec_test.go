// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestParseCmdline(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected options
	}{
		{
			name: "Test with exec flag only",
			args: []string{"kexec", "-e"},
			expected: options{
				exec: true,
			},
		},
		{
			name: "Test with short load flag",
			args: []string{"kexec", "-l", "/path/to/kernel"},
			expected: options{
				load:       true,
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test with long load flag",
			args: []string{"kexec", "--load", "/path/to/kernel"},
			expected: options{
				load:       true,
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test with multiple flags 1",
			args: []string{"kexec", "-l", "-d", "/path/to/kernel"},
			expected: options{
				load:       true,
				debug:      true,
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test with multiple flags 2",
			args: []string{"kexec", "-l", "-d", "-L", "/path/to/kernel"},
			expected: options{
				load:        true,
				debug:       true,
				loadSyscall: true,
				kernelpath:  "/path/to/kernel",
			},
		},
		{
			name: "Test unix style flags 1",
			args: []string{"kexec", "-ld", "/path/to/kernel"},
			expected: options{
				load:       true,
				debug:      true,
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test unix style flags 2",
			args: []string{"kexec", "-ldL", "/path/to/kernel"},
			expected: options{
				load:        true,
				debug:       true,
				loadSyscall: true,
				kernelpath:  "/path/to/kernel",
			},
		},
		{
			name: "Test kernel arg special case 1",
			args: []string{"kexec", "-l", "/path/to/kernel", "-d"},
			expected: options{
				load:       true,
				debug:      true,
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test kernel arg special case 2",
			args: []string{"kexec", "-l", "/path/to/kernel", "-d", "-L"},
			expected: options{
				load:        true,
				debug:       true,
				loadSyscall: true,
				kernelpath:  "/path/to/kernel",
			},
		},
		{
			name: "Test kernel arg special case 3",
			args: []string{"kexec", "-l", "/path/to/kernel", "-d", "-c", "${CMDLINE}"},
			expected: options{
				load:       true,
				debug:      true,
				cmdline:    "${CMDLINE}",
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test command line",
			args: []string{"kexec", "-l", "-c", "${CMDLINE}", "/path/to/kernel"},
			expected: options{
				load:       true,
				cmdline:    "${CMDLINE}",
				kernelpath: "/path/to/kernel",
			},
		},
		{
			name: "Test append command line",
			args: []string{"kexec", "-l", "--append", "${CMDLINE}", "/path/to/kernel"},
			expected: options{
				load:          true,
				appendCmdline: "${CMDLINE}",
				kernelpath:    "/path/to/kernel",
			},
		},
		{
			name: "Test all set",
			args: []string{"kexec", "-c", "${CMDLINE}", "-d", "--dtb", "foo", "-e", "-x", "/some/file", "-i", "/path/to/initrd", "-l", "-L", "--module", "/mod1", "--reuse-cmdline", "/path/to/kernel"},
			expected: options{
				cmdline:      "${CMDLINE}",
				debug:        true,
				dtb:          "foo",
				exec:         true,
				extra:        "/some/file",
				initramfs:    "/path/to/initrd",
				load:         true,
				loadSyscall:  true,
				modules:      []string{"/mod1"},
				reuseCmdline: true,
				kernelpath:   "/path/to/kernel",
			},
		},
		{
			name: "Test all set special case",
			args: []string{"kexec", "-c", "${CMDLINE}", "-d", "--dtb", "foo", "-e", "-x", "/some/file", "-i", "/path/to/initrd", "-l", "/path/to/kernel", "-L", "--module", "/mod1", "--reuse-cmdline"},
			expected: options{
				cmdline:      "${CMDLINE}",
				debug:        true,
				dtb:          "foo",
				exec:         true,
				extra:        "/some/file",
				initramfs:    "/path/to/initrd",
				load:         true,
				loadSyscall:  true,
				modules:      []string{"/mod1"},
				reuseCmdline: true,
				kernelpath:   "/path/to/kernel",
			},
		},
		{
			name: "Test all set, multiple modules 1",
			args: []string{"kexec", "-c", "${CMDLINE}", "-d", "--dtb", "foo", "-e", "-x", "/some/file", "-i", "/path/to/initrd", "-l", "-L", "--module", "/mod1", "--module", "/mod2", "--reuse-cmdline", "/path/to/kernel"},
			expected: options{
				cmdline:      "${CMDLINE}",
				debug:        true,
				dtb:          "foo",
				exec:         true,
				extra:        "/some/file",
				initramfs:    "/path/to/initrd",
				load:         true,
				loadSyscall:  true,
				modules:      []string{"/mod1", "/mod2"},
				reuseCmdline: true,
				kernelpath:   "/path/to/kernel",
			},
		},
		{
			name: "Test all set, multiple modules 2",
			args: []string{"kexec", "-c", "${CMDLINE}", "-d", "--dtb", "foo", "-e", "-x", "/some/file", "-i", "/path/to/initrd", "-l", "-L", "--module", "/mod1", "--reuse-cmdline", "--module", "/mod2", "/path/to/kernel"},
			expected: options{
				cmdline:      "${CMDLINE}",
				debug:        true,
				dtb:          "foo",
				exec:         true,
				extra:        "/some/file",
				initramfs:    "/path/to/initrd",
				load:         true,
				loadSyscall:  true,
				modules:      []string{"/mod1", "/mod2"},
				reuseCmdline: true,
				kernelpath:   "/path/to/kernel",
			},
		},
		{
			name: "Test all set unix style flags",
			args: []string{"kexec", "-delL", "/path/to/kernel"},
			expected: options{
				debug:       true,
				exec:        true,
				load:        true,
				loadSyscall: true,
				kernelpath:  "/path/to/kernel",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := &options{}
			f := flag.NewFlagSet("test", flag.ExitOnError)
			o.parseCmdline(tc.args, f)

			tc.expected.purgatory = "default"
			if !reflect.DeepEqual(o, &tc.expected) {
				t.Errorf("\nexpected: %+v\ngot:     %+v", tc.expected, o)
			}
		})
	}
}

func TestHackLoadFlagValue(t *testing.T) {
	tests := []struct {
		input  []string
		output []string
	}{
		{
			input:  []string{"kexec", "-l", "/kernel", "-i", "/initramfs.cpio", "-c", "${CMDLINE}"},
			output: []string{"kexec", "-l", "/kernel", "-i", "/initramfs.cpio", "-c", "${CMDLINE}"},
		},
		{
			input:  []string{"kexec", "--load", "/kernel", "-i", "/initramfs.cpio", "-c", "${CMDLINE}"},
			output: []string{"kexec", "--load", "/kernel", "-i", "/initramfs.cpio", "-c", "${CMDLINE}"},
		},
		{
			input:  []string{"kexec", "-l", "-i", "/initramfs.cpio", "-c", "${CMDLINE}", "/kernel"},
			output: []string{"kexec", "-l", setButEmpty, "-i", "/initramfs.cpio", "-c", "${CMDLINE}", "/kernel"},
		},
		{
			input:  []string{"kexec", "--load", "-i", "/initramfs.cpio", "-c", "${CMDLINE}", "/kernel"},
			output: []string{"kexec", "--load", setButEmpty, "-i", "/initramfs.cpio", "-c", "${CMDLINE}", "/kernel"},
		},
	}

	for _, test := range tests {
		result := hackLoadFlagValue(test.input)
		if !reflect.DeepEqual(result, test.output) {
			t.Errorf("\nInput:    %v\nExpected: %v\nGot:      %v", test.input, test.output, result)
		}
	}
}
