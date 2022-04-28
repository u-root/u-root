// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/vfile"
)

func TestLinuxLabel(t *testing.T) {
	dir := t.TempDir()

	osKernel, err := os.Create(filepath.Join(dir, "kernel"))
	if err != nil {
		t.Fatal(err)
	}

	osInitrd, err := os.Create(filepath.Join(dir, "initrd"))
	if err != nil {
		t.Fatal(err)
	}

	k, _ := url.Parse("http://127.0.0.1/kernel")
	i1, _ := url.Parse("http://127.0.0.1/initrd1")
	i2, _ := url.Parse("http://127.0.0.1/initrd2")
	httpKernel, _ := curl.LazyFetch(k)
	httpInitrd1, _ := curl.LazyFetch(i1)
	httpInitrd2, _ := curl.LazyFetch(i2)

	for _, tt := range []struct {
		desc string
		img  *LinuxImage
		want string
	}{
		{
			desc: "os.File",
			img: &LinuxImage{
				Kernel: osKernel,
				Initrd: osInitrd,
			},
			want: fmt.Sprintf("Linux(kernel=%s/kernel initrd=%s/initrd)", dir, dir),
		},
		{
			desc: "lazy file",
			img: &LinuxImage{
				Kernel: uio.NewLazyFile(filepath.Join(dir, "kernel")),
				Initrd: uio.NewLazyFile(filepath.Join(dir, "initrd")),
			},
			want: fmt.Sprintf("Linux(kernel=%s/kernel initrd=%s/initrd)", dir, dir),
		},
		{
			desc: "concat lazy file",
			img: &LinuxImage{
				Kernel: uio.NewLazyFile(filepath.Join(dir, "kernel")),
				Initrd: CatInitrds(
					uio.NewLazyFile(filepath.Join(dir, "initrd")),
					uio.NewLazyFile(filepath.Join(dir, "initrd")),
				),
			},
			want: fmt.Sprintf("Linux(kernel=%s/kernel initrd=%s/initrd,%s/initrd)", dir, dir, dir),
		},
		{
			desc: "curl lazy file",
			img: &LinuxImage{
				Kernel: httpKernel,
				Initrd: CatInitrds(
					httpInitrd1,
					httpInitrd2,
				),
			},
			want: "Linux(kernel=http://127.0.0.1/kernel initrd=http://127.0.0.1/initrd1,http://127.0.0.1/initrd2)",
		},
		{
			desc: "verified file",
			img: &LinuxImage{
				Kernel: &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				Initrd: CatInitrds(
					&vfile.File{Reader: nil, FileName: "/boot/initrd1"},
					&vfile.File{Reader: nil, FileName: "/boot/initrd2"},
				),
			},
			want: "Linux(kernel=/boot/foobar initrd=/boot/initrd1,/boot/initrd2)",
		},
		{
			desc: "no initrd",
			img: &LinuxImage{
				Kernel:     &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				Initrd:     nil,
				DeviceTree: nil,
			},
			want: "Linux(kernel=/boot/foobar)",
		},
		{
			desc: "dtb file",
			img: &LinuxImage{
				Kernel:     &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				Initrd:     &vfile.File{Reader: nil, FileName: "/boot/initrd"},
				DeviceTree: &vfile.File{Reader: nil, FileName: "/boot/board.dtb"},
			},
			want: "Linux(kernel=/boot/foobar initrd=/boot/initrd dtb=/boot/board.dtb)",
		},
		{
			desc: "dtb file, no initrd",
			img: &LinuxImage{
				Kernel:     &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				DeviceTree: &vfile.File{Reader: nil, FileName: "/boot/board.dtb"},
			},
			want: "Linux(kernel=/boot/foobar dtb=/boot/board.dtb)",
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

func TestCopyToFile(t *testing.T) {
	want := "abcdefg hijklmnop"
	buf := bytes.NewReader([]byte(want))

	f, err := copyToFileIfNotRegular(buf, true)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(f.Name())
	got, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got %s, expected %s", string(got), want)
	}
}

func TestLinuxRank(t *testing.T) {
	testRank := 2
	img := &LinuxImage{BootRank: testRank}
	l := img.Rank()
	if l != testRank {
		t.Fatalf("Expected Image rank %d, got %d", testRank, l)
	}
}
