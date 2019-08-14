// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipxe

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

func TestIpxeConfig(t *testing.T) {
	content1 := "1111"
	content2 := "2222"

	for i, tt := range []struct {
		desc       string
		schemeFunc func() urlfetch.Schemes
		curl       *url.URL
		want       *boot.LinuxImage
		err        error
	}{
		{
			desc: "all files exist, simple config with no cmdline",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
				conf := `#!ipxe
				kernel http://someplace.com/foobar/pxefiles/kernel
				initrd http://someplace.com/foobar/pxefiles/initrd
				boot`
				fs.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				fs.Add("someplace.com", "/foobar/pxefiles/kernel", content1)
				fs.Add("someplace.com", "/foobar/pxefiles/initrd", content2)
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
			desc: "all files exist, simple config, no initrd",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
				Kernel: errorReader{&urlfetch.URLError{
					URL: &url.URL{
						Scheme: "http",
						Host:   "someplace.com",
						Path:   "/foobar/pxefiles/kernel",
					},
					Err: urlfetch.ErrNoSuchFile,
				}},
				Initrd:  nil,
				Cmdline: "",
			},
		},
		{
			desc: "config file does not exist",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			err: &urlfetch.URLError{
				URL: &url.URL{
					Scheme: "http",
					Host:   "someplace.com",
					Path:   "/foobar/pxefiles/ipxeconfig",
				},
				Err: urlfetch.ErrNoSuchHost,
			},
		},
		{
			desc: "invalid config",
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			schemeFunc: func() urlfetch.Schemes {
				conf := `#!ipxe
				kernel tftp://1.2.3.4/foobar/pxefiles/kernel
                                initrd http://someplace.com/someinitrd.gz
				boot`

				tftp := urlfetch.NewMockScheme("tftp")
				tftp.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)

				http := urlfetch.NewMockScheme("http")
				http.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				http.Add("someplace.com", "/someinitrd.gz", content2)

				s := make(urlfetch.Schemes)
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
			schemeFunc: func() urlfetch.Schemes {
				s := make(urlfetch.Schemes)
				fs := urlfetch.NewMockScheme("http")
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
			got, err := ParseConfigWithSchemes(tt.curl, tt.schemeFunc())
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewConfigWithSchemes() got %v, want %v", err, tt.err)
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
