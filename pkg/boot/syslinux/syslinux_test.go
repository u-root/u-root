// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/curl"
	"github.com/u-root/u-root/pkg/uio"
)

func mustReadAll(r io.ReaderAt) string {
	if r == nil {
		return ""
	}
	b, err := uio.ReadAll(r)
	if err != nil {
		return fmt.Sprintf("read error: %s", err)
	}
	return string(b)
}

type errorReader struct {
	err error
}

func (e errorReader) ReadAt(p []byte, n int64) (int, error) {
	return 0, e.err
}

func sameBootImage(got, want boot.OSImage) error {
	if got.Label() != want.Label() {
		return fmt.Errorf("got image label %s, want %s", got.Label(), want.Label())
	}

	if gotLinux, ok := got.(*boot.LinuxImage); ok {
		wantLinux, ok := want.(*boot.LinuxImage)
		if !ok {
			return fmt.Errorf("got image %s is Linux image, but %s is not", got, want)
		}

		// Same kernel?
		if !uio.ReaderAtEqual(gotLinux.Kernel, wantLinux.Kernel) {
			return fmt.Errorf("got kernel %s, want %s", mustReadAll(gotLinux.Kernel), mustReadAll(wantLinux.Kernel))
		}

		// Same initrd?
		if !uio.ReaderAtEqual(gotLinux.Initrd, wantLinux.Initrd) {
			return fmt.Errorf("got initrd %s, want %s", mustReadAll(gotLinux.Initrd), mustReadAll(wantLinux.Initrd))
		}

		// Same cmdline?
		if gotLinux.Cmdline != wantLinux.Cmdline {
			return fmt.Errorf("got cmdline %s, want %s", gotLinux.Cmdline, wantLinux.Cmdline)
		}
		return nil
	}

	return fmt.Errorf("non-Linux images not supported yet")
}

func TestParseGeneral(t *testing.T) {
	kernel1 := "kernel1"
	kernel2 := "kernel2"
	globalInitrd := "globalInitrd"
	initrd1 := "initrd1"
	initrd2 := "initrd2"

	newMockScheme := func() *curl.MockScheme {
		fs := curl.NewMockScheme("tftp")
		fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
		fs.Add("1.2.3.4", "/foobar/pxefiles/kernel1", kernel1)
		fs.Add("1.2.3.4", "/foobar/pxefiles/kernel2", kernel2)
		fs.Add("1.2.3.4", "/foobar/pxefiles/global_initrd", globalInitrd)
		fs.Add("1.2.3.4", "/foobar/pxefiles/initrd1", initrd1)
		fs.Add("1.2.3.4", "/foobar/pxefiles/initrd2", initrd2)

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

					label foo
					kernel ./pxefiles/kernel1
					append earlyprintk=ttyS0 printk=ttyS0

					label bar
					kernel ./pxefiles/kernel2
					append console=ttyS0`,
			},
			want: []boot.OSImage{
				&boot.LinuxImage{
					Name:    "foo",
					Kernel:  strings.NewReader(kernel1),
					Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
				},
				&boot.LinuxImage{
					Name:    "bar",
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
					// See TODO in pxe.go initrd handling.
					Initrd:  strings.NewReader(initrd2),
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
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			fs := newMockScheme()
			for filename, content := range tt.configFiles {
				fs.Add("1.2.3.4", filename, content)
			}
			s := make(curl.Schemes)
			s.Register(fs.Scheme, fs)
			s.Register(http.Scheme, http)

			wd := &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			}

			got, err := ParseConfigFile(context.Background(), s, "pxelinux.cfg/default", wd)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("AppendFile() got %v, want %v", err, tt.err)
			} else if err != nil {
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("ParseConfigFile yielded %d images, want %d images", len(got), len(tt.want))
			}

			for i, want := range tt.want {
				if err := sameBootImage(got[i], want); err != nil {
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
		wd         *url.URL
		err        error
	}{
		{
			name:       "no schemes",
			s:          nil,
			configFile: "pxelinux.cfg/default",
			wd: &url.URL{
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
			wd: &url.URL{
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
			wd:         nil,
			err: &curl.URLError{
				URL: &url.URL{},
				Err: curl.ErrNoSuchScheme,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfigFile(context.Background(), tt.s, tt.configFile, tt.wd)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseConfigFile() = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	for i, tt := range []struct {
		url  string
		wd   *url.URL
		err  bool
		want *url.URL
	}{
		{
			url: "default",
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "192.168.1.1",
				Path:   "/foobar/pxelinux.cfg",
			},
			want: &url.URL{
				Scheme: "tftp",
				Host:   "192.168.1.1",
				Path:   "/foobar/pxelinux.cfg/default",
			},
		},
		{
			url: "http://192.168.2.1/configs/your-machine.cfg",
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "192.168.1.1",
				Path:   "/foobar/pxelinux.cfg",
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "192.168.2.1",
				Path:   "/configs/your-machine.cfg",
			},
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			got, err := parseURL(tt.url, tt.wd)
			if (err != nil) != tt.err {
				t.Errorf("Wanted error (%v), but got %v", tt.err, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseURL() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
