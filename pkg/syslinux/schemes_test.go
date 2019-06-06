// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syslinux

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/uio"
)

func TestGetFile(t *testing.T) {
	for i, tt := range []struct {
		scheme func() *MockScheme
		url    *url.URL
		err    error
		want   string
	}{
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
				return s
			},
			want: "haha",
			url: &url.URL{
				Scheme: "fooftp",
				Host:   "192.168.0.1",
				Path:   "/foo/pxelinux.cfg/default",
			},
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				return s
			},
			url: &url.URL{
				Scheme: "nosuch",
			},
			err: ErrNoSuchScheme,
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				return s
			},
			url: &url.URL{
				Scheme: "fooftp",
				Host:   "someotherplace",
			},
			err: ErrNoSuchHost,
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				s.Add("somehost", "somefile", "somecontent")
				return s
			},
			url: &url.URL{
				Scheme: "fooftp",
				Host:   "somehost",
				Path:   "/someotherfile",
			},
			err: ErrNoSuchFile,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			fs := tt.scheme()
			s := make(Schemes)
			s.Register(fs.Scheme, fs)

			// Test both GetFile and LazyGetFile.
			for _, f := range []func(url *url.URL) (io.ReaderAt, error){
				s.GetFile,
				s.LazyGetFile,
			} {
				r, err := f(tt.url)
				if uErr, ok := err.(*URLError); ok && uErr.Err != tt.err {
					t.Errorf("GetFile() = %v, want %v", uErr.Err, tt.err)
				} else if !ok && err != tt.err {
					t.Errorf("GetFile() = %v, want %v", err, tt.err)
				}
				if err != nil {
					return
				}
				content, err := ioutil.ReadAll(uio.Reader(r))
				if err != nil {
					t.Errorf("bytes.Buffer read returned an error? %v", err)
				}
				if got, want := string(content), tt.want; got != want {
					t.Errorf("GetFile() = %v, want %v", got, want)
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
