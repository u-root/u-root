// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cenkalti/backoff/v4"
	"github.com/u-root/uio/uio"
)

var (
	errTest = errors.New("test error")
	testURL = &url.URL{
		Scheme: "fooftp",
		Host:   "192.168.0.1",
		Path:   "/foo/pxelinux.cfg/default",
	}
)

var tests = []struct {
	name string
	// scheme returns a scheme for testing and a MockScheme to
	// confirm number of calls to Fetch. The distinction is useful
	// when MockScheme is decorated by a SchemeWithRetries. In many
	// cases, the same value is returned twice.
	scheme         func() (FileScheme, *MockScheme)
	url            *url.URL
	err            error
	want           string
	wantFetchCount uint
}{
	{
		name: "successful fetch",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			return s, s
		},
		url:            testURL,
		want:           "haha",
		wantFetchCount: 1,
	},
	{
		name: "scheme does not exist",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			return s, s
		},
		url: &url.URL{
			Scheme: "nosuch",
		},
		err:            ErrNoSuchScheme,
		wantFetchCount: 0,
	},
	{
		name: "host does not exist",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			return s, s
		},
		url: &url.URL{
			Scheme: "fooftp",
			Host:   "someotherplace",
		},
		err:            ErrNoSuchHost,
		wantFetchCount: 1,
	},
	{
		name: "file does not exist",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("somehost", "somefile", "somecontent")
			return s, s
		},
		url: &url.URL{
			Scheme: "fooftp",
			Host:   "somehost",
			Path:   "/someotherfile",
		},
		err:            ErrNoSuchFile,
		wantFetchCount: 1,
	},
	{
		name: "always err",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			s.SetErr(errTest, 9999)
			return s, s
		},
		url:            testURL,
		err:            errTest,
		wantFetchCount: 1,
	},
	{
		name: "retries but not necessary",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			r := &SchemeWithRetries{
				Scheme: s,
				// backoff.ZeroBackOff so unit tests run fast.
				BackOff: backoff.WithMaxRetries(&backoff.ZeroBackOff{}, 10),
			}
			return r, s
		},
		url:            testURL,
		want:           "haha",
		wantFetchCount: 1,
	},
	{
		name: "not enough retries",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			s.SetErr(errTest, 9999)
			r := &SchemeWithRetries{
				Scheme: s,
				// backoff.ZeroBackOff so unit tests run fast.
				BackOff: backoff.WithMaxRetries(&backoff.ZeroBackOff{}, 10),
			}
			return r, s
		},
		url:            testURL,
		err:            errTest,
		wantFetchCount: 11,
	},
	{
		name: "sufficient retries",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			s.SetErr(errTest, 5)
			r := &SchemeWithRetries{
				Scheme: s,
				// backoff.ZeroBackOff so unit tests run fast.
				BackOff: backoff.WithMaxRetries(&backoff.ZeroBackOff{}, 10),
			}
			return r, s
		},
		url:            testURL,
		want:           "haha",
		wantFetchCount: 6,
	},
	{
		name: "retry filter",
		scheme: func() (FileScheme, *MockScheme) {
			s := NewMockScheme("fooftp")
			s.Add("192.168.0.1", "/foo/pxelinux.cfg/default", "haha")
			s.SetErr(errTest, 5)
			r := &SchemeWithRetries{
				DoRetry: func(u *url.URL, err error) bool {
					return err != errTest
				},
				Scheme: s,
				// backoff.ZeroBackOff so unit tests run fast.
				BackOff: backoff.WithMaxRetries(&backoff.ZeroBackOff{}, 10),
			}
			return r, s
		},
		url:            testURL,
		err:            errTest,
		wantFetchCount: 1,
	},
}

func TestFetchWithoutCache(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%02d: %s", i, tt.name), func(t *testing.T) {
			var r io.Reader
			var err error

			fs, ms := tt.scheme()
			s := make(Schemes)
			s.Register(ms.Scheme, fs)

			r, err = s.FetchWithoutCache(context.TODO(), tt.url)
			if uErr, ok := err.(*URLError); ok && uErr.Err != tt.err {
				t.Errorf("FetchWithoutCache() = %v, want %v", uErr.Err, tt.err)
			} else if !ok && err != tt.err {
				t.Errorf("FetchWithoutCache() = %v, want %v", err, tt.err)
			}

			// Check number of calls before reading the file.
			numCalled := ms.NumCalled(tt.url)
			if numCalled != tt.wantFetchCount {
				t.Errorf("number times Fetch() called = %v, want %v",
					ms.NumCalled(tt.url), tt.wantFetchCount)
			}
			if err != nil {
				return
			}

			// Read the entire file.
			content, err := io.ReadAll(r)
			if err != nil {
				t.Errorf("bytes.Buffer read returned an error? %v", err)
			}
			if got, want := string(content), tt.want; got != want {
				t.Errorf("Fetch() = %v, want %v", got, want)
			}

			// Check number of calls after reading the file.
			numCalled = ms.NumCalled(tt.url)
			if numCalled != tt.wantFetchCount {
				t.Errorf("number times Fetch() called = %v, want %v",
					ms.NumCalled(tt.url), tt.wantFetchCount)
			}
		})
	}
}

func TestFetch(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%02d: %s", i, tt.name), func(t *testing.T) {
			var r io.ReaderAt
			var err error

			fs, ms := tt.scheme()
			s := make(Schemes)
			s.Register(ms.Scheme, fs)

			r, err = s.Fetch(context.TODO(), tt.url)
			if uErr, ok := err.(*URLError); ok && uErr.Err != tt.err {
				t.Errorf("Fetch() = %v, want %v", uErr.Err, tt.err)
			} else if !ok && err != tt.err {
				t.Errorf("Fetch() = %v, want %v", err, tt.err)
			}

			// Check number of calls before reading the file.
			numCalled := ms.NumCalled(tt.url)
			if numCalled != tt.wantFetchCount {
				t.Errorf("number times Fetch() called = %v, want %v",
					ms.NumCalled(tt.url), tt.wantFetchCount)
			}
			if err != nil {
				return
			}

			// Read the entire file.
			content, err := io.ReadAll(uio.Reader(r))
			if err != nil {
				t.Errorf("bytes.Buffer read returned an error? %v", err)
			}
			if got, want := string(content), tt.want; got != want {
				t.Errorf("Fetch() = %v, want %v", got, want)
			}

			// Check number of calls after reading the file.
			numCalled = ms.NumCalled(tt.url)
			if numCalled != tt.wantFetchCount {
				t.Errorf("number times Fetch() called = %v, want %v",
					ms.NumCalled(tt.url), tt.wantFetchCount)
			}
		})
	}
}

func TestLazyFetch(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test #%02d: %s", i, tt.name), func(t *testing.T) {
			var r io.ReaderAt
			var err error

			fs, ms := tt.scheme()
			s := make(Schemes)
			s.Register(ms.Scheme, fs)

			r, err = s.LazyFetch(tt.url)
			// Errors are deferred to when file is read except for ErrNoSuchScheme.
			if tt.err == ErrNoSuchScheme {
				if uErr, ok := err.(*URLError); ok && uErr.Err != ErrNoSuchScheme {
					t.Errorf("LazyFetch() = %v, want %v", uErr.Err, tt.err)
				}
			} else if err != nil {
				t.Errorf("LazyFetch() = %v, want nil", err)
			}

			// Check number of calls before reading the file.
			numCalled := ms.NumCalled(tt.url)
			if numCalled != 0 {
				t.Errorf("number times Fetch() called = %v, want 0", numCalled)
			}
			if err != nil {
				return
			}

			// Read the entire file.
			content, err := io.ReadAll(uio.Reader(r))
			if uErr, ok := err.(*URLError); ok && uErr.Err != tt.err {
				t.Errorf("ReadAll() = %v, want %v", uErr.Err, tt.err)
			} else if !ok && err != tt.err {
				t.Errorf("ReadAll() = %v, want %v", err, tt.err)
			}
			if got, want := string(content), tt.want; got != want {
				t.Errorf("ReadAll() = %v, want %v", got, want)
			}

			// Check number of calls after reading the file.
			numCalled = ms.NumCalled(tt.url)
			if numCalled != tt.wantFetchCount {
				t.Errorf("number times Fetch() called = %v, want %v",
					ms.NumCalled(tt.url), tt.wantFetchCount)
			}
		})
	}
}

func TestHttpFetches(t *testing.T) {
	c := "fetch content"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, c)
	}))
	defer ts.Close()

	fURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("url.Parse(%s) = %v, want no error", ts.URL, err)
	}

	// Fetch need to fetch the content as is.
	fetchFile, err := Fetch(context.Background(), fURL)
	if err != nil {
		t.Errorf("Fetch(context.Background(), %s) = %v, want no error", fURL, err)
	}
	got, err := io.ReadAll(io.NewSectionReader(fetchFile, 0, int64(len(c))))
	if err != nil {
		t.Errorf("io.ReadAll(%v) = %v, want no error", fetchFile, err)
	}
	if string(got) != c {
		t.Errorf("got %s, want %s", got, c)
	}

	// FetchWithoutCache need to fetch the content as is.
	fetchFileNoCache, err := FetchWithoutCache(context.Background(), fURL)
	if err != nil {
		t.Errorf("FetchWithoutCache(context.Background(), %s) = %v, want no error", fURL, err)
	}
	got, err = io.ReadAll(fetchFileNoCache)
	if err != nil {
		t.Errorf("io.ReadAll(%s) = %v, want no error", fetchFileNoCache, err)
	}
	if string(got) != c {
		t.Errorf("got %s, want %s", got, c)
	}
}
