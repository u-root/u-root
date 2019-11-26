// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package urlfetch

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/u-root/u-root/pkg/uio"
)

func TestFetch(t *testing.T) {
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
				return NewMockScheme("fooftp")
			},
			url: &url.URL{
				Scheme: "nosuch",
			},
			err: ErrNoSuchScheme,
		},
		{
			scheme: func() *MockScheme {
				return NewMockScheme("fooftp")
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

			// Test both Fetch and LazyFetch.
			for _, f := range []func(url *url.URL) (io.ReaderAt, error){
				s.Fetch,
				s.LazyFetch,
			} {
				r, err := f(tt.url)
				if uErr, ok := err.(*URLError); ok && uErr.Err != tt.err {
					t.Errorf("Fetch() = %v, want %v", uErr.Err, tt.err)
				} else if !ok && err != tt.err {
					t.Errorf("Fetch() = %v, want %v", err, tt.err)
				}
				if err != nil {
					return
				}
				content, err := ioutil.ReadAll(uio.Reader(r))
				if err != nil {
					t.Errorf("bytes.Buffer read returned an error? %v", err)
				}
				if got, want := string(content), tt.want; got != want {
					t.Errorf("Fetch() = %v, want %v", got, want)
				}
			}
		})
	}
}
