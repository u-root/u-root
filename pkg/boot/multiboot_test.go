// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLabel(t *testing.T) {
	dir := t.TempDir()

	osKernel, err := os.Create(filepath.Join(dir, "kernel"))
	if err != nil {
		t.Fatal(err)
	}

	cmdLine := "console=ttyS0"
	for _, tt := range []struct {
		desc string
		img  *MultibootImage
		want string
	}{
		{
			desc: "include name",
			img: &MultibootImage{
				Name:   "multiboot_test",
				Kernel: osKernel,
			},
			want: "multiboot_test",
		},
		{
			desc: "wo name and ibft",
			img: &MultibootImage{
				Kernel:  osKernel,
				Cmdline: cmdLine,
			},
			want: fmt.Sprintf("Multiboot(kernel=%s/kernel cmdline=%s iBFT=<nil>)", dir, cmdLine),
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.img.Label()
			if got != tt.want {
				t.Errorf("Label() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRank(t *testing.T) {
	testRank := 2
	img := &MultibootImage{BootRank: testRank}
	l := img.Rank()
	if l != testRank {
		t.Fatalf("Expected Image rank %d, got %d", testRank, l)
	}
}
