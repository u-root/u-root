// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/boottest"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/curl"
)

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("parsing %q failed: %v", s, err))
	}
	return u
}

type errorReader struct {
	err error
}

func (e errorReader) ReadAt(p []byte, n int64) (int, error) {
	return 0, e.err
}

func TestParseGeneral(t *testing.T) {
	kernel1 := "kernel1"
	kernel2 := "kernel2"
	globalInitrd := "globalInitrd"
	initrd1 := "initrd1"
	initrd2 := "initrd2"
	xengz := "xengz"
	mboot := "mboot.c32"
	boardDTB := "board.dtb"

	newMockScheme := func() *curl.MockScheme {
		fs := curl.NewMockScheme("tftp")
		fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
		fs.Add("1.2.3.4", "/foobar/pxefiles/kernel1", kernel1)
		fs.Add("1.2.3.4", "/foobar/pxefiles/kernel2", kernel2)
		fs.Add("1.2.3.4", "/foobar/pxefiles/global_initrd", globalInitrd)
		fs.Add("1.2.3.4", "/foobar/pxefiles/initrd1", initrd1)
		fs.Add("1.2.3.4", "/foobar/pxefiles/initrd2", initrd2)
		fs.Add("1.2.3.4", "/foobar/pxefiles/board.dtb", boardDTB)
		fs.Add("1.2.3.4", "/foobar/xen.gz", xengz)
		fs.Add("1.2.3.4", "/foobar/mboot.c32", mboot)

		fs.Add("2.3.4.5", "/barfoo/pxefiles/kernel1", kernel1)
		return fs
	}
	http := curl.NewMockScheme("http")
	http.Add("someplace.com", "/initrd2", initrd2)

	for i, tt := range []struct {
		desc        string
		configFiles map[string]string
		want        []boot.OSImage
		err         error
	}{
		{
			desc: "all files exist, simple config with cmdline initrd",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo
					kernel ./pxefiles/kernel1
					append initrd=./pxefiles/global_initrd`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Initrd:  strings.NewReader(globalInitrd),
					Cmdline: "initrd=./pxefiles/global_initrd",
				},
			},
		},
		{
			desc: "empty label",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo`,
			},
			want: nil,
		},

		{
			desc: "all files exist, simple config with directive initrd",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo
					kernel ./pxefiles/kernel1
					initrd ./pxefiles/initrd1
					append foo=bar`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Initrd:  strings.NewReader(initrd1),
					Cmdline: "foo=bar",
				},
			},
		},
		{
			desc: "all files exist, simple config, no initrd",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo
					kernel ./pxefiles/kernel1`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Initrd:  nil,
					Cmdline: "",
				},
			},
		},
		{
			desc: "all files exist, simple config with two initrd files",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label multi-initrd
					kernel ./pxefiles/kernel1
					initrd ./pxefiles/initrd1,./pxefiles/initrd2
					append foo=bar`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "multi-initrd",
					Kernel:  strings.NewReader(kernel1),
					Initrd:  boot.CatInitrds(strings.NewReader(initrd1), strings.NewReader(initrd2)),
					Cmdline: "foo=bar",
				},
			},
		},
		{
			desc: "an initrd file missing, config with two initrd files",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label multi-initrd
					kernel ./pxefiles/kernel1
					initrd ./pxefiles/initrd1,./pxefiles/no-initrd-here
					append foo=bar`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:   "multi-initrd",
					Kernel: strings.NewReader(kernel1),
					Initrd: errorReader{&curl.URLError{
						URL: &url.URL{
							Scheme: "tftp",
							Host:   "1.2.3.4",
							Path:   "/foobar/pxefiles/no-initrd-here",
						},
						Err: curl.ErrNoSuchFile,
					}},
					Cmdline: "foo=bar",
				},
			},
		},
		{
			desc: "kernel does not exist, simple config",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo
					kernel ./pxefiles/does-not-exist`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name: "foo",
					Kernel: errorReader{&curl.URLError{
						URL: &url.URL{
							Scheme: "tftp",
							Host:   "1.2.3.4",
							Path:   "/foobar/pxefiles/does-not-exist",
						},
						Err: curl.ErrNoSuchFile,
					}},
					Initrd:  nil,
					Cmdline: "",
				},
			},
		},
		{
			desc: "config file does not exist",
			err: &curl.URLError{
				URL: &url.URL{
					Scheme: "tftp",
					Host:   "1.2.3.4",
					Path:   "/foobar/pxelinux.cfg/default",
				},
				Err: curl.ErrNoSuchFile,
			},
		},
		{
			desc: "empty config",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": "",
			},
			want: nil,
		},
		{
			desc: "valid config with two Entries",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo

					label bar
					menu label Bla Bla Bla
					kernel ./pxefiles/kernel2
					append console=ttyS0

					label foo
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:    "Bla Bla Bla",
					Kernel:  strings.NewReader(kernel2),
					Cmdline: "console=ttyS0",
				},
			},
		},
		{
			desc: "menu default, linux directives",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					label bar
					menu label Bla Bla Bla
					kernel ./pxefiles/kernel2
					append console=ttyS0

					label foo
					menu default
					linux ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:    "Bla Bla Bla",
					Kernel:  strings.NewReader(kernel2),
					Cmdline: "console=ttyS0",
				},
			},
		},
		{
			desc: "valid config with two Entries, and a nerfdefault override",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo

					nerfdefault bar

					label foo
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0

					label bar
					kernel ./pxefiles/kernel2
					append console=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "bar",
					Kernel:  strings.NewReader(kernel2),
					Cmdline: "console=ttyS0",
				},
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
			},
		},
		{
			desc: "valid config with two Entries, and a nerfdefault override, order agnostic",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					nerfdefault bar

					default foo

					label foo
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0

					label bar
					kernel ./pxefiles/kernel2
					append console=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "bar",
					Kernel:  strings.NewReader(kernel2),
					Cmdline: "console=ttyS0",
				},
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
			},
		},

		{
			desc: "valid config with global APPEND directive",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					append foo=bar

					label foo
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0

					label bar
					kernel ./pxefiles/kernel2

					label baz
					kernel ./pxefiles/kernel2
					append -`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:   "foo",
					Kernel: strings.NewReader(kernel1),
					// Does not contain global APPEND.
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:   "bar",
					Kernel: strings.NewReader(kernel2),
					// Contains only global APPEND.
					Cmdline: "foo=bar",
				},
				&boot.LinuxImage{
					Name:   "baz",
					Kernel: strings.NewReader(kernel2),
					// "APPEND -" means ignore global APPEND.
					Cmdline: "",
				},
			},
		},
		{
			desc: "valid config with global APPEND with initrd",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default mcnulty
					append initrd=./pxefiles/global_initrd

					label mcnulty
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0

					label lester
					kernel ./pxefiles/kernel1

					label omar
					kernel ./pxefiles/kernel2
					append -

					label stringer
					kernel ./pxefiles/kernel2
					initrd ./pxefiles/initrd2
				`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:   "mcnulty",
					Kernel: strings.NewReader(kernel1),
					// Does not contain global APPEND.
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:   "lester",
					Kernel: strings.NewReader(kernel1),
					Initrd: strings.NewReader(globalInitrd),
					// Contains only global APPEND.
					Cmdline: "initrd=./pxefiles/global_initrd",
				},
				&boot.LinuxImage{
					Name:   "omar",
					Kernel: strings.NewReader(kernel2),
					// "APPEND -" means ignore global APPEND.
					Cmdline: "",
				},
				&boot.LinuxImage{
					Name:   "stringer",
					Kernel: strings.NewReader(kernel2),
					Initrd: strings.NewReader(initrd2),

					// TODO: See syslinux initrd handling. This SHOULD be
					//
					// initrd=./pxefiles/global_initrd initrd=./pxefiles/initrd2
					//
					// https://wiki.syslinux.org/wiki/index.php?title=Directives/append
					Cmdline: "initrd=./pxefiles/global_initrd",
				},
			},
		},
		{
			desc: "default label does not exist",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `default not-exist`,
			},
			want: nil,
		},
		{
			desc: "multi-scheme valid config",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
				default sheeeit

				label sheeeit
				kernel ./pxefiles/kernel2
				initrd http://someplace.com/initrd2`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:   "sheeeit",
					Kernel: strings.NewReader(kernel2),
					Initrd: strings.NewReader(initrd2),
				},
			},
		},
		{
			desc: "valid config with three includes",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default mcnulty

					include installer/txt.cfg
					include installer/stdmenu.cfg

					menu begin advanced
					  menu title Advanced Options
					  include installer/stdmenu.cfg
					menu end
				`,

				"/foobar/installer/txt.cfg": `
					label mcnulty
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0
				`,

				"/foobar/installer/stdmenu.cfg": `
					label omar
					kernel ./pxefiles/kernel2
				`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "mcnulty",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:   "omar",
					Kernel: strings.NewReader(kernel2),
				},
			},
		},
		{
			desc: "multiboot images",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo

					label bar
					menu label Bla Bla Bla
					kernel mboot.c32
					append xen.gz console=none --- ./pxefiles/kernel1 foobar hahaha --- ./pxefiles/initrd1

					label mbootnomodules
					kernel mboot.c32
					append xen.gz

					label foo
					linux mboot.c32
					append earlyprintk=ttyS0 printk=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(mboot),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.MultibootImage{
					Name:    "Bla Bla Bla",
					Kernel:  strings.NewReader(xengz),
					Cmdline: "console=none",
					Modules: []multiboot.Module{
						{
							Module:  strings.NewReader(kernel1),
							Cmdline: "./pxefiles/kernel1 foobar hahaha",
						},
						{
							Module:  strings.NewReader(initrd1),
							Cmdline: "./pxefiles/initrd1",
						},
					},
				},
				&boot.MultibootImage{
					Name:   "mbootnomodules",
					Kernel: strings.NewReader(xengz),
				},
			},
		},
		{
			desc: "simple config with all required file and fdt",
			configFiles: map[string]string{
				"/foobar/pxelinux.cfg/default": `
					default foo
					label foo
					kernel ./pxefiles/kernel1
					initrd ./pxefiles/global_initrd
					append foo=bar
					fdt ./pxefiles/board.dtb`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Initrd:  strings.NewReader(globalInitrd),
					Cmdline: "foo=bar",
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			fs := newMockScheme()
			for filename, content := range tt.configFiles {
				fs.Add("1.2.3.4", filename, content)
			}
			s := make(curl.Schemes)
			s.Register(fs.Scheme, fs)
			s.Register(http.Scheme, http)

			rootdir := &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/",
			}

			got, err := ParseConfigFile(context.Background(), s, "pxelinux.cfg/default", rootdir, "foobar")
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("AppendFile() got %v, want %v", err, tt.err)
			} else if err != nil {
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("ParseConfigFile yielded %d images, want %d images", len(got), len(tt.want))
			}

			for i, want := range tt.want {
				if err := boottest.SameBootImage(got[i], want); err != nil {
					t.Errorf("Boot image index %d not same: %v", i, err)
				}
			}
		})
	}
}

func TestParseCorner(t *testing.T) {
	for _, tt := range []struct {
		name       string
		s          curl.Schemes
		configFile string
		rootdir    *url.URL
		wd         string
		err        error
	}{
		{
			name:       "no schemes",
			s:          nil,
			configFile: "pxelinux.cfg/default",
			rootdir: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: &curl.URLError{
				URL: &url.URL{
					Scheme: "tftp",
					Host:   "1.2.3.4",
					Path:   "/foobar/pxelinux.cfg/default",
				},
				Err: curl.ErrNoSuchScheme,
			},
		},
		{
			name:       "no scheme and config file",
			s:          nil,
			configFile: "",
			rootdir: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: &curl.URLError{
				URL: &url.URL{
					Scheme: "tftp",
					Host:   "1.2.3.4",
					Path:   "/foobar",
				},
				Err: curl.ErrNoSuchScheme,
			},
		},
		{
			name:       "no scheme, config file, and working dir",
			s:          nil,
			configFile: "",
			rootdir:    nil,
			err: &curl.URLError{
				URL: &url.URL{},
				Err: curl.ErrNoSuchScheme,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfigFile(context.Background(), tt.s, tt.configFile, tt.rootdir, tt.wd)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseConfigFile() = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	for _, tt := range []struct {
		filename string
		rootdir  *url.URL
		wd       string
		want     *url.URL
	}{
		{
			filename: "foobar",
			rootdir:  mustParseURL("http://[2001::1]:18282/"),
			wd:       "files/more",
			want:     mustParseURL("http://[2001::1]:18282/files/more/foobar"),
		},
		{
			filename: "/foobar",
			rootdir:  mustParseURL("http://[2001::1]:18282"),
			wd:       "files/more",
			want:     mustParseURL("http://[2001::1]:18282/foobar"),
		},
		{
			filename: "http://[2002::2]/blabla",
			rootdir:  mustParseURL("http://[2001::1]:18282/files"),
			wd:       "more",
			want:     mustParseURL("http://[2002::2]/blabla"),
		},
		{
			filename: "http://[2002::2]/blabla",
			rootdir:  nil,
			want:     mustParseURL("http://[2002::2]/blabla"),
		},
	} {
		got, err := parseURL(tt.filename, tt.rootdir, tt.wd)
		if err != nil {
			t.Errorf("parseURL(%q, %s, %s) = %v, want %v", tt.filename, tt.rootdir, tt.wd, err, nil)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseURL(%q, %s, %s) = %v, want %v", tt.filename, tt.rootdir, tt.wd, got, tt.want)
		}
	}
}

func FuzzParseSyslinuxConfig(f *testing.F) {
	dirPath := f.TempDir()

	path := filepath.Join(dirPath, "isolinux.cfg")

	log.SetOutput(io.Discard)
	log.SetFlags(0)

	// get seed corpora from testdata_new files
	seeds, err := filepath.Glob("testdata/*/*/isolinux.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}

	f.Add([]byte("lABel 0\nAppend initrd"))
	f.Add([]byte("lABel 0\nkernel mboot.c32\nAppend ---"))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			return
		}

		// do not allow arbitrary files reads
		if bytes.Contains(data, []byte("include")) {
			return
		}

		err := os.WriteFile(path, data, 0o777)
		if err != nil {
			t.Fatalf("Failed to create configfile '%v':%v", path, err)
		}

		ParseLocalConfig(context.Background(), dirPath)
	})
}
