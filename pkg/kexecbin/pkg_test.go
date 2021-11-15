// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexecbin

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Equal(a, b []string) bool {
	fmt.Printf("a: %v, b: %v\n", a, b)
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestCommandLine(t *testing.T) {
	tests := []struct {
		name              string
		kernelFilePath    string
		kernelCommandline string
		initrdFilePath    string
		dtFilePath        string
		wantCommandline   []string
	}{
		{
			name:              "Empty Input",
			kernelFilePath:    "",
			kernelCommandline: "",
			initrdFilePath:    "",
			dtFilePath:        "",
			wantCommandline:   []string{"--reuse-cmdline"},
		},
		{
			name:              "Kernel only",
			kernelFilePath:    "/test/kernel.vmlinuz",
			kernelCommandline: "",
			initrdFilePath:    "",
			dtFilePath:        "",
			wantCommandline:   []string{"-l", "/test/kernel.vmlinuz", "--reuse-cmdline"},
		},
		{
			name:              "Initrd only",
			kernelFilePath:    "",
			kernelCommandline: "",
			initrdFilePath:    "/test/initrd.img",
			dtFilePath:        "",
			wantCommandline:   []string{"--reuse-cmdline", "--initrd=/test/initrd.img"},
		},
		{
			name:              "All options",
			kernelFilePath:    "/test/kernel.vmlinuz",
			kernelCommandline: "console=ttyS0,115200",
			initrdFilePath:    "/test/initrd.img",
			dtFilePath:        "/test/dt.dts",
			wantCommandline: []string{"-l", "/test/kernel.vmlinuz", "--command-line=console=ttyS0,115200",
				"--initrd=/test/initrd.img", "--dtb=/test/dt.dts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandLine := buildCommandline(tt.kernelFilePath, tt.kernelCommandline, tt.initrdFilePath, tt.dtFilePath)

			DeviceTreePaths := []string{"/sys/firmware/fdt", "/proc/device-tree"}

			for _, dtFilePath := range DeviceTreePaths {
				if _, err := os.Stat(dtFilePath); err == nil {
					tt.wantCommandline = append(tt.wantCommandline, "--dtb="+dtFilePath)
					break
				}
			}

			if !reflect.DeepEqual(commandLine, tt.wantCommandline) {
				t.Errorf("buildCommandLine fails. Want %v but have %v", tt.wantCommandline, commandLine)
			}
		})
	}
}
