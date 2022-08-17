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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/vfile"
	"golang.org/x/sys/unix"
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
				Kernel: &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				Initrd: nil,
			},
			want: "Linux(kernel=/boot/foobar)",
		},
		{
			desc: "dtb file",
			img: &LinuxImage{
				Kernel: &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				Initrd: &vfile.File{Reader: nil, FileName: "/boot/initrd"},
				KexecOpts: linux.KexecOptions{
					DTB: &vfile.File{Reader: nil, FileName: "/boot/board.dtb"},
				},
			},
			want: "Linux(kernel=/boot/foobar initrd=/boot/initrd dtb=/boot/board.dtb)",
		},
		{
			desc: "dtb file, no initrd",
			img: &LinuxImage{
				Kernel: &vfile.File{Reader: nil, FileName: "/boot/foobar"},
				KexecOpts: linux.KexecOptions{
					DTB: &vfile.File{Reader: nil, FileName: "/boot/board.dtb"},
				},
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

func checkReadOnly(t *testing.T, f *os.File) {
	t.Helper()
	wr := unix.O_RDWR | unix.O_WRONLY
	if am, err := unix.FcntlInt(f.Fd(), unix.F_GETFL, 0); err == nil && am&wr != 0 {
		t.Errorf("file %v opened for write, want read only", f)
	}
}

// checkFilePath checks if src and dst file are same file of fsrc were actually a os.File.
func checkFilePath(t *testing.T, fsrc io.ReaderAt, fdst *os.File) {
	t.Helper()
	if f, ok := fsrc.(*os.File); ok {
		if r, _ := mount.IsTmpRamfs(f.Name()); r {
			// Src is a file on tmpfs.
			if f.Name() != fdst.Name() {
				t.Errorf("Got a copied file %s, want skipping copy and use original file %s", fdst.Name(), f.Name())
			}
		}
	}
}

func setupTestFile(t *testing.T, path, content string) *os.File {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	n, err := f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	if n != len([]byte(content)) {
		t.Fatalf("want %d bytes written, but got %d", len([]byte(content)), n)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("could not close test file: %v", err)
	}
	nf, err := os.Open(path)
	if err != nil {
		t.Fatalf("could not open test file: %v", err)
	}
	return nf
}

// GenerateCatDummyInitrd return padded string from the given list of strings following the same padding format of CatInitrds.
func GenerateCatDummyInitrd(t *testing.T, initrds ...string) string {
	var ins []io.ReaderAt
	for _, c := range initrds {
		ins = append(ins, strings.NewReader(c))
	}
	final := CatInitrds(ins...)
	d, err := io.ReadAll(uio.Reader(final))
	if err != nil {
		t.Fatalf("failed to generate concatenated initrd : %v", err)
	}
	return string(d)
}

type wantData struct {
	loadedImage *LoadedLinuxImage
	cleanup     func()
	err         error
}

func TestLoadLinuxImage(t *testing.T) {
	testDir := t.TempDir()

	for _, tt := range []struct {
		name string
		li   *LinuxImage
		want wantData
	}{
		{
			name: "kernel is nil",
			li:   &LinuxImage{Kernel: nil},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: nil,
				},
				err: errNilKernel,
			},
		},
		{
			name: "basic happy case w/o initrd",
			li: &LinuxImage{
				Kernel: strings.NewReader("testkernel"),
			},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: setupTestFile(t, filepath.Join(testDir, "basic_happy_case_wo_initrd_bzimage"), "testkernel"),
				},
				err: nil,
			},
		},
		{
			name: "basic happy case w/ initrd",
			li: &LinuxImage{
				Kernel: strings.NewReader("testkernel"),
				Initrd: strings.NewReader("testinitrd"),
			},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: setupTestFile(t, filepath.Join(testDir, "basic_happy_case_w_initrd_bzImage"), "testkernel"),
					Initrd: setupTestFile(t, filepath.Join(testDir, "basic_happy_case_w_initrd_initramfs"), "testinitrd"),
				},
				err: nil,
			},
		},
		{
			name: "empty initrd, with DTB present", // Expect DTB appended to loaded initrd.
			li: &LinuxImage{
				Kernel: strings.NewReader("testkernel"),
				Initrd: nil,
				KexecOpts: linux.KexecOptions{
					DTB: strings.NewReader("testdtb"),
				},
			},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: setupTestFile(t, filepath.Join(testDir, "empty_inird_w_dtb_present_bzImage"), "testkernel"),
					Initrd: setupTestFile(t, filepath.Join(testDir, "empty_inird_w_dtb_present_initramfs"), "testdtb"),
				},
				err: nil,
			},
		},
		{
			name: "non-empty initrd, with DTB present", // Expect DTB appended to loaded initrd.
			li: &LinuxImage{
				Kernel: strings.NewReader("testkernel"),
				Initrd: strings.NewReader("testinitrd"),
				KexecOpts: linux.KexecOptions{
					DTB: strings.NewReader("testdtb"),
				},
			},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: setupTestFile(t, filepath.Join(testDir, "non_empty_inird_w_dtb_present_bzImage"), "testkernel"),
					Initrd: setupTestFile(t, filepath.Join(testDir, "non_empty_inird_w_dtb_present_initramfs"), GenerateCatDummyInitrd(t, "testinitrd", "testdtb")),
				},
				err: nil,
			},
		},
		{
			name: "oringnal kernel and initrd are files, skip copying",
			li: &LinuxImage{
				Kernel: setupTestFile(t, filepath.Join(testDir, "original_kernel_and_initrd_are_files_skip_copying_bzImage"), "testkernel"),
				Initrd: setupTestFile(t, filepath.Join(testDir, "original_kernel_and_initrd_are_files_skip_copying_initramfs"), "testinitrd"),
			},
			want: wantData{
				loadedImage: &LoadedLinuxImage{
					Kernel: setupTestFile(t, filepath.Join(testDir, "original_kernel_and_initrd_are_files_skip_copying_2_bzImage"), "testkernel"),
					Initrd: setupTestFile(t, filepath.Join(testDir, "original_kernel_and_initrd_are_files_skip_copying_2_initramfs"), "testinitrd"),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotImage, _, gotErr := loadLinuxImage(tt.li, true)
			if gotErr != nil {
				if gotErr != tt.want.err {
					t.Errorf("got error %v, want %v", gotErr, tt.want.err)
				}
				return
			}
			// Kernel is opened as read only, and contents match that from original LinuxImage.
			checkReadOnly(t, gotImage.Kernel)
			// If src is a read-only *os.File on tmpfs, shoukd skip copying.
			checkFilePath(t, tt.li.Kernel, gotImage.Kernel)
			kernelBytes, err := io.ReadAll(gotImage.Kernel)
			if err != nil {
				t.Fatalf("could not read kernel from loaded image: %v", err)
			}
			wantBytes, err := io.ReadAll(tt.want.loadedImage.Kernel)
			if err != nil {
				t.Fatalf("could not read expected kernel: %v", err)
			}
			if string(kernelBytes) != string(wantBytes) {
				t.Errorf("got kernel %s, want %s", string(kernelBytes), string(wantBytes))
			}
			// Initrd, if present, is opened as read only, and contents match that from original LinuxImage.
			// OR original initrd, with DTB appended.
			if tt.li.Initrd != nil {
				checkReadOnly(t, gotImage.Initrd)
				// If src is a read-only *os.File on tmpfs, should skip copying.
				checkFilePath(t, tt.li.Initrd, gotImage.Initrd)
				initrdBytes, err := io.ReadAll(gotImage.Initrd)
				if err != nil {
					t.Fatalf("could not read initrd from loaded image: %v", err)
				}
				wantInitrdBytes, err := io.ReadAll(tt.want.loadedImage.Initrd)
				if err != nil {
					t.Fatalf("could not read expected initrd: %v", err)
				}
				// Initrd involves appending, use cmp.Diff for catching the diff, easier to debug.
				if diff := cmp.Diff(string(initrdBytes), string(wantInitrdBytes)); diff != "" {
					t.Errorf("got initrd %s, want %s, diff (+got, -want): %s", string(initrdBytes), string(wantInitrdBytes), diff)
				}
			}
		})
	}
}
