// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipxe

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/pxe"
	"github.com/u-root/u-root/pkg/uio"
)

func TestIpxeConfig(t *testing.T) {
	content1 := "1111"
	content2 := "2222"

	type config struct {
		kernel    string
		kernelErr error
		initrd    string
		initrdErr error
		cmdline   string
	}

	for i, tt := range []struct {
		desc       string
		schemeFunc func() pxe.Schemes
		curl       *url.URL
		config     *Config
		want       config
		err        error
	}{
		{
			desc: "all files exist, simple config with no cmdline",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{
				kernel: content1,
				initrd: content2,
			},
		},
		{
			desc: "all files exist, simple config, no initrd",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{
				kernel: content1,
			},
		},
		{
			desc: "kernel does not exist, simple config",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{
				kernelErr: &pxe.URLError{
					URL: &url.URL{
						Scheme: "http",
						Host:   "someplace.com",
						Path:   "/foobar/pxefiles/kernel",
					},
					Err: pxe.ErrNoSuchFile,
				},
				initrd:  "",
				cmdline: "",
			},
		},
		{
			desc: "config file does not exist",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
				s.Register(fs.Scheme, fs)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			err: &pxe.URLError{
				URL: &url.URL{
					Scheme: "http",
					Host:   "someplace.com",
					Path:   "/foobar/pxefiles/ipxeconfig",
				},
				Err: pxe.ErrNoSuchHost,
			},
		},
		{
			desc: "invalid config",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{},
		},
		{
			desc: "valid config with kernel cmdline args",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{
				kernel:  content1,
				cmdline: "earlyprintk=ttyS0 printk=ttyS0",
			},
		},
		{
			desc: "multi-scheme valid config",
			schemeFunc: func() pxe.Schemes {
				conf := `#!ipxe
				kernel tftp://1.2.3.4/foobar/pxefiles/kernel
                                initrd http://someplace.com/someinitrd.gz
				boot`

				tftp := pxe.NewMockScheme("tftp")
				tftp.Add("1.2.3.4", "/foobar/pxefiles/kernel", content1)

				http := pxe.NewMockScheme("http")
				http.Add("someplace.com", "/foobar/pxefiles/ipxeconfig", conf)
				http.Add("someplace.com", "/someinitrd.gz", content2)

				s := make(pxe.Schemes)
				s.Register(tftp.Scheme, tftp)
				s.Register(http.Scheme, http)
				return s
			},
			curl: &url.URL{
				Scheme: "http",
				Host:   "someplace.com",
				Path:   "/foobar/pxefiles/ipxeconfig",
			},
			want: config{
				kernel: content1,
				initrd: content2,
			},
		},
		{
			desc: "valid config with unsupported cmds",
			schemeFunc: func() pxe.Schemes {
				s := make(pxe.Schemes)
				fs := pxe.NewMockScheme("http")
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
			want: config{
				kernel: content1,
				initrd: content2,
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d] %s", i, tt.desc), func(t *testing.T) {
			c, err := NewConfigWithSchemes(tt.curl, tt.schemeFunc())
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewConfigWithSchemes() got %v, want %v", err, tt.err)
				return
			} else if err != nil {
				return
			}

			got := c.BootImage
			want := tt.want

			// Same kernel?
			if got.Kernel == nil && (len(want.kernel) > 0 || want.kernelErr != nil) {
				t.Errorf("want kernel, got none")
			}
			if got.Kernel != nil {
				k, err := uio.ReadAll(got.Kernel)
				if !reflect.DeepEqual(err, want.kernelErr) {
					t.Errorf("could not read kernel. got: %v, want %v", err, want.kernelErr)
				}
				if got, want := string(k), want.kernel; got != want {
					t.Errorf("got kernel %s, want %s", got, want)
				}
			}

			// Same initrd?
			if got.Initrd == nil && (len(want.initrd) > 0 || want.initrdErr != nil) {
				t.Errorf("want initrd, got none")
			}
			if got.Initrd != nil {
				i, err := uio.ReadAll(got.Initrd)
				if err != want.initrdErr {
					t.Errorf("could not read initrd. got: %v, want %v", err, want.initrdErr)
				}
				if got, want := string(i), want.initrd; got != want {
					t.Errorf("got initrd %s, want %s", got, want)
				}
			}

			// Same cmdline?
			if got, want := got.Cmdline, want.cmdline; got != want {
				t.Errorf("got cmdline %s, want %s", got, want)
			}
		})
	}
}
