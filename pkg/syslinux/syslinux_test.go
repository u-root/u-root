// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/urlfetch"
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

func TestAppendFile(t *testing.T) {
	content1 := "1111"
	content2 := "2222"
	content3 := "3333"
	content4 := "4444"

	for i, tt := range []struct {
		desc          string
		configFileURI string
		schemeFunc    func() urlfetch.Schemes
		wd            *url.URL
		want          *Config
		err           error
	}{
		{
			desc:          "all files exist, simple config with cmdline initrd",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel
				append initrd=./pxefiles/initrd`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/initrd", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel:  strings.NewReader(content1),
						Initrd:  strings.NewReader(content2),
						Cmdline: "initrd=./pxefiles/initrd",
					},
				},
			},
		},
		{
			desc:          "all files exist, simple config with directive initrd",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel
				initrd ./pxefiles/initrd
				append foo=bar`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/initrd", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel:  strings.NewReader(content1),
						Initrd:  strings.NewReader(content2),
						Cmdline: "foo=bar",
					},
				},
			},
		},
		{
			desc:          "all files exist, simple config, no initrd",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel:  strings.NewReader(content1),
						Initrd:  nil,
						Cmdline: "",
					},
				},
			},
		},
		{
			desc:          "kernel does not exist, simple config",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				label foo
				kernel ./pxefiles/kernel`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel: errorReader{&urlfetch.URLError{
							URL: &url.URL{
								Scheme: "tftp",
								Host:   "1.2.3.4",
								Path:   "/foobar/pxefiles/kernel",
							},
							Err: urlfetch.ErrNoSuchFile,
						}},
						Initrd:  nil,
						Cmdline: "",
					},
				},
			},
		},
		{
			desc:          "config file does not exist",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: &urlfetch.URLError{
				URL: &url.URL{
					Scheme: "tftp",
					Host:   "1.2.3.4",
					Path:   "/foobar/pxelinux.cfg/default",
				},
				Err: urlfetch.ErrNoSuchHost,
			},
		},
		{
			desc:          "empty config",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", "")
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "",
			},
		},
		{
			desc:          "valid config with two Entries",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo

				label foo
				kernel ./pxefiles/fookernel
				append earlyprintk=ttyS0 printk=ttyS0

				label bar
				kernel ./pxefiles/barkernel
				append console=ttyS0`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/fookernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/barkernel", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel:  strings.NewReader(content1),
						Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"bar": {
						Kernel:  strings.NewReader(content2),
						Cmdline: "console=ttyS0",
					},
				},
			},
		},
		{
			desc:          "valid config with global APPEND directive",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default foo
				append foo=bar

				label foo
				kernel ./pxefiles/fookernel
				append earlyprintk=ttyS0 printk=ttyS0

				label bar
				kernel ./pxefiles/barkernel

				label baz
				kernel ./pxefiles/barkernel
				append -`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/fookernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/barkernel", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "foo",
				Entries: map[string]*boot.LinuxImage{
					"foo": {
						Kernel: strings.NewReader(content1),
						// Does not contain global APPEND.
						Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"bar": {
						Kernel: strings.NewReader(content2),
						// Contains only global APPEND.
						Cmdline: "foo=bar",
					},
					"baz": {
						Kernel: strings.NewReader(content2),
						// "APPEND -" means ignore global APPEND.
						Cmdline: "",
					},
				},
			},
		},
		{
			desc:          "valid config with global APPEND with initrd",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default mcnulty
				append initrd=./pxefiles/normal_person

				label mcnulty
				kernel ./pxefiles/copkernel
				append earlyprintk=ttyS0 printk=ttyS0

				label lester
				kernel ./pxefiles/copkernel

				label omar
				kernel ./pxefiles/drugkernel
				append -

				label stringer
				kernel ./pxefiles/drugkernel
				initrd ./pxefiles/criminal
				`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/pxefiles/copkernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/drugkernel", content2)
				fs.Add("1.2.3.4", "/foobar/pxefiles/normal_person", content3)
				fs.Add("1.2.3.4", "/foobar/pxefiles/criminal", content4)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "mcnulty",
				Entries: map[string]*boot.LinuxImage{
					"mcnulty": {
						Kernel: strings.NewReader(content1),
						// Does not contain global APPEND.
						Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"lester": {
						Kernel: strings.NewReader(content1),
						Initrd: strings.NewReader(content3),
						// Contains only global APPEND.
						Cmdline: "initrd=./pxefiles/normal_person",
					},
					"omar": {
						Kernel: strings.NewReader(content2),
						// "APPEND -" means ignore global APPEND.
						Cmdline: "",
					},
					"stringer": {
						Kernel: strings.NewReader(content2),
						// See TODO in pxe.go initrd handling.
						Initrd:  strings.NewReader(content4),
						Cmdline: "initrd=./pxefiles/normal_person",
					},
				},
			},
		},
		{
			desc:          "default label does not exist",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				conf := `default avon`

				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			err: ErrDefaultEntryNotFound,
			want: &Config{
				DefaultEntry: "avon",
			},
		},
		{
			desc:          "multi-scheme valid config",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				conf := `default sheeeit

				label sheeeit
				kernel ./pxefiles/kernel
				initrd http://someplace.com/someinitrd.gz`

				tftp := urlfetch.NewMockScheme("tftp")
				tftp.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				tftp.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				tftp.Add("1.2.3.4", "/foobar/pxefiles/kernel", content2)

				http := urlfetch.NewMockScheme("http")
				http.Add("someplace.com", "/someinitrd.gz", content3)

				s := make(urlfetch.Schemes)
				s.Register(tftp.Scheme, tftp)
				s.Register(http.Scheme, http)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "sheeeit",
				Entries: map[string]*boot.LinuxImage{
					"sheeeit": {
						Kernel: strings.NewReader(content2),
						Initrd: strings.NewReader(content3),
					},
				},
			},
		},
		{
			desc:          "valid config with three includes",
			configFileURI: "pxelinux.cfg/default",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("tftp")
				fs.Add("1.2.3.4", "/foobar/pxelinux.0", "")
				conf := `default mcnulty

				include installer/txt.cfg
				include installer/stdmenu.cfg

				menu begin advanced
				  menu title Advanced Options
				  include installer/stdmenu.cfg
				menu end
				`

				txt := `
				label mcnulty
				kernel ./pxefiles/copkernel
				append earlyprintk=ttyS0 printk=ttyS0
				`

				stdmenu := `
				label omar
				kernel ./pxefiles/drugkernel
				`
				fs.Add("1.2.3.4", "/foobar/pxelinux.cfg/default", conf)
				fs.Add("1.2.3.4", "/foobar/installer/stdmenu.cfg", stdmenu)
				fs.Add("1.2.3.4", "/foobar/installer/txt.cfg", txt)
				fs.Add("1.2.3.4", "/foobar/pxefiles/copkernel", content1)
				fs.Add("1.2.3.4", "/foobar/pxefiles/drugkernel", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			wd: &url.URL{
				Scheme: "tftp",
				Host:   "1.2.3.4",
				Path:   "/foobar",
			},
			want: &Config{
				DefaultEntry: "mcnulty",
				Entries: map[string]*boot.LinuxImage{
					"mcnulty": {
						Kernel:  strings.NewReader(content1),
						Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
					},
					"omar": {
						Kernel: strings.NewReader(content2),
					},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			s := tt.schemeFunc()

			par := newParser(tt.wd)
			par.schemes = s

			if err := par.appendFile(tt.configFileURI); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("AppendFile() got %v, want %v", err, tt.err)
			} else if err != nil {
				return
			}
			c := par.config

			if got, want := c.DefaultEntry, tt.want.DefaultEntry; got != want {
				t.Errorf("DefaultEntry got %v, want %v", got, want)
			}

			for labelName, want := range tt.want.Entries {
				t.Run(fmt.Sprintf("label %s", labelName), func(t *testing.T) {
					got, ok := c.Entries[labelName]
					if !ok {
						t.Errorf("Config label %v does not exist", labelName)
						return
					}

					// Same kernel?
					if !uio.ReaderAtEqual(got.Kernel, want.Kernel) {
						t.Errorf("got kernel %s, want %s", mustReadAll(got.Kernel), mustReadAll(want.Kernel))
					}

					// Same initrd?
					if !uio.ReaderAtEqual(got.Initrd, want.Initrd) {
						t.Errorf("got initrd %s, want %s", mustReadAll(got.Initrd), mustReadAll(want.Initrd))
					}

					// Same cmdline?
					if got.Cmdline != want.Cmdline {
						t.Errorf("got cmdline %s, want %s", got.Cmdline, want.Cmdline)
					}
				})
			}

			// Check that the parser didn't make up Entries.
			for labelName := range c.Entries {
				if _, ok := tt.want.Entries[labelName]; !ok {
					t.Errorf("config has extra label %s, but not wanted", labelName)
				}
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
