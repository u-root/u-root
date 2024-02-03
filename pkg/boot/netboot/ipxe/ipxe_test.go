// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipxe

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
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
	"github.com/u-root/uio/uio"
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

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("parsing %q failed: %v", s, err))
	}
	return u
}

func TestParseURL(t *testing.T) {
	for _, tt := range []struct {
		filename string
		wd       *url.URL
		want     *url.URL
	}{
		{
			filename: "foobar",
			wd:       mustParseURL("http://[2001::1]:18282/files/more"),
			want:     mustParseURL("http://[2001::1]:18282/files/more/foobar"),
		},
		{
			filename: "/foobar",
			wd:       mustParseURL("http://[2001::1]:18282/files/more"),
			want:     mustParseURL("http://[2001::1]:18282/foobar"),
		},
		{
			filename: "http://[2002::2]/blabla",
			wd:       mustParseURL("http://[2001::1]:18282/files/more"),
			want:     mustParseURL("http://[2002::2]/blabla"),
		},
		{
			filename: "http://[2002::2]/blabla",
			wd:       nil,
			want:     mustParseURL("http://[2002::2]/blabla"),
		},
	} {
		got, err := parseURL(tt.filename, tt.wd)
		if err != nil {
			t.Errorf("parseURL(%q, %s) = %v, want %v", tt.filename, tt.wd, err, nil)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseURL(%q, %s) = %v, want %v", tt.filename, tt.wd, got, tt.want)
		}
	}
}

func TestIpxeConfig(t *testing.T) {
	content1 := "1111"
	content2 := "2222"
	content512_1 := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	content512_2 := "h123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"he23456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hell456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hello56789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellow6789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellowo789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellowor89abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"helloworl9abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	content1024 := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"h123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"he23456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hell456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hello56789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellow6789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellowo789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"hellowor89abcdef0123456789abcdef0123456789abcdef0123456789abcdef" +
		"helloworl9abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	for i, tt := range []struct {
		desc       string
		schemeFunc func() curl.Schemes
		curl       *url.URL
		want       *boot.LinuxImage
		err        error
	}{
		{
			desc: "all files exist, simple config with no cmdline, one relative file path",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				initrd initrd-file
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content2),
			},
		},
		{
			desc: "all files exist, simple config with no cmdline, one relative file path, premature end",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				initrd initrd-file`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content2),
			},
		},
		{
			desc: "all files exist, simple config with no cmdline, one relative file path, concatenate initrd",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				initrd initrd-file.001,initrd-file.002
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file.001", content512_1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file.002", content512_2)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content1024),
			},
		},
		{
			desc: "all files exist, simple config with no cmdline, one relative file path, multiline initrd",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				initrd initrd-file.001
				initrd initrd-file.002
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file.001", content512_1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd-file.002", content512_2)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content1024),
			},
		},
		{
			desc: "all files exist, simple config, no initrd",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
			},
		},
		{
			desc: "comments and blank lines",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				# the next line is blank

				kernel http://someplace.com/foobar/pxefiles/kernel
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
			},
		},
		{
			desc: "kernel does not exist, simple config",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: errorReader{&curl.URLError{
					URL: &url.URL{
						Scheme: "http",
						Host:   "someplace.com",
						Path:   "/foobar/pxefiles/kernel",
					},
					Err: curl.ErrNoSuchFile,
				}},
				Initrd:  nil,
				Cmdline: "",
			},
		},
		{
			desc: "config file does not exist",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			err: &curl.URLError{
				URL: &url.URL{
					Scheme: "http",
					Host:   "someplace.com",
					Path:   "/foobar/pxefiles/ipxeconfig",
				},
				Err: curl.ErrNoSuchHost,
			},
		},
		{
			desc: "invalid config",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", "")
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			err: ErrNotIpxeScript,
		},
		{
			desc: "empty config",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{},
		},
		{
			desc: "valid config with kernel cmdline args",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel earlyprintk=ttyS0 printk=ttyS0
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel:  strings.NewReader(content1),
				Cmdline: "earlyprintk=ttyS0 printk=ttyS0",
			},
		},
		{
			desc: "multi-scheme valid config",
			schemeFunc: func() curl.Schemes {
				conf := `#!ipxe
				kernel tftp://1.2.3.4/foobar/pxefiles/kernel
                                initrd http://someplace.com/someinitrd.gz
				boot`

				tftp := curl.NewMockScheme("tftp")
				tftp.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)

				http := curl.NewMockScheme("http")
				http.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				http.Add("someplace.com", "/someinitrd.gz", content2)

				s := make(curl.Schemes)
				s.Register(tftp.Scheme, tftp)
				s.Register(http.Scheme, http)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content2),
			},
		},
		{
			desc: "valid config with unsupported cmds",
			schemeFunc: func() curl.Schemes {
				s := make(curl.Schemes)
				fs := curl.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
                                initrd http://someplace.com/someinitrd.gz
                                set ip 0.0.0.0
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/someinitrd.gz", content2)
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: &boot.LinuxImage{
				Kernel: strings.NewReader(content1),
				Initrd: strings.NewReader(content2),
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			got, err := ParseConfig(context.Background(), ulogtest.Logger{t}, tt.curl, tt.schemeFunc())
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseConfig() got %v, want %v", err, tt.err)
				return
			} else if err != nil {
				return
			}
			want := tt.want

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
}
