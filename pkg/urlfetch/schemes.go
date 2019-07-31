// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package urlfetch implements routines to fetch files given a URL.
//
// urlfetch currently supports HTTP, TFTP, local files, and a retrying HTTP
// client.
package urlfetch

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/u-root/u-root/pkg/uio"
	"pack.ag/tftp"
)

var (
	// ErrNoSuchScheme is returned by Schemes.Fetch and
	// Schemes.LazyFetch if there is no registered FileScheme
	// implementation for the given URL scheme.
	ErrNoSuchScheme = errors.New("no such scheme")
)

// FileScheme represents the implementation of a URL scheme and gives access to
// fetching files of that scheme.
//
// For example, an http FileScheme implementation would fetch files using
// the HTTP protocol.
type FileScheme interface {
	// Fetch returns a reader that gives the contents of `u`.
	//
	// It may do so by fetching `u` and placing it in a buffer, or by
	// returning an io.ReaderAt that fetchs the file.
	Fetch(u *url.URL) (io.ReaderAt, error)
}

var (
	// DefaultHTTPClient is the default HTTP FileScheme.
	//
	// It is not recommended to use this for HTTPS. We recommend creating an
	// http.Client that accepts only a private pool of certificates.
	DefaultHTTPClient = NewHTTPClient(http.DefaultClient)

	// DefaultTFTPClient is the default TFTP FileScheme.
	DefaultTFTPClient = NewTFTPClient()

	// DefaultSchemes are the schemes supported by default.
	DefaultSchemes = Schemes{
		"tftp": DefaultTFTPClient,
		"http": DefaultHTTPClient,
		"file": &LocalFileClient{},
	}
)

// URLError is an error involving URLs.
type URLError struct {
	URL *url.URL
	Err error
}

// Error implements error.Error.
func (s *URLError) Error() string {
	return fmt.Sprintf("encountered error %v with %q", s.Err, s.URL)
}

// IsURLError returns true iff err is a URLError.
func IsURLError(err error) bool {
	_, ok := err.(*URLError)
	return ok
}

// Schemes is a map of URL scheme identifier -> implementation that can
// fetch a file for that scheme.
type Schemes map[string]FileScheme

// RegisterScheme calls DefaultSchemes.Register.
func RegisterScheme(scheme string, fs FileScheme) {
	DefaultSchemes.Register(scheme, fs)
}

// Register registers a scheme identified by `scheme` to be `fs`.
func (s Schemes) Register(scheme string, fs FileScheme) {
	s[scheme] = fs
}

// Fetch fetchs a file via DefaultSchemes.
func Fetch(u *url.URL) (io.ReaderAt, error) {
	return DefaultSchemes.Fetch(u)
}

// file is an io.ReaderAt with a nice Stringer.
type file struct {
	io.ReaderAt

	url *url.URL
}

// String implements fmt.Stringer.
func (f file) String() string {
	return f.url.String()
}

// Fetch fetchs the file with the given `u`. `u.Scheme` is used to
// select the FileScheme via `s`.
//
// If `s` does not contain a FileScheme for `u.Scheme`, ErrNoSuchScheme is
// returned.
func (s Schemes) Fetch(u *url.URL) (io.ReaderAt, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}
	r, err := fg.Fetch(u)
	if err != nil {
		return nil, &URLError{URL: u, Err: err}
	}
	return &file{ReaderAt: r, url: u}, nil
}

// LazyFetch calls LazyFetch on DefaultSchemes.
func LazyFetch(u *url.URL) (io.ReaderAt, error) {
	return DefaultSchemes.LazyFetch(u)
}

// LazyFetch returns a reader that will Fetch the file given by `u` when
// Read is called, based on `u`s scheme. See Schemes.Fetch for more
// details.
func (s Schemes) LazyFetch(u *url.URL) (io.ReaderAt, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}

	return &file{
		url: u,
		ReaderAt: uio.NewLazyOpenerAt(func() (io.ReaderAt, error) {
			r, err := fg.Fetch(u)
			if err != nil {
				return nil, &URLError{URL: u, Err: err}
			}
			return r, nil
		}),
	}, nil
}

// TFTPClient implements FileScheme for TFTP files.
type TFTPClient struct {
	opts []tftp.ClientOpt
}

// NewTFTPClient returns a new TFTP client based on the given tftp.ClientOpt.
func NewTFTPClient(opts ...tftp.ClientOpt) FileScheme {
	return &TFTPClient{
		opts: opts,
	}
}

// Fetch implements FileScheme.Fetch.
func (t *TFTPClient) Fetch(u *url.URL) (io.ReaderAt, error) {
	// TODO(hugelgupf): These clients are basically stateless, except for
	// the options. Figure out whether you actually have to re-establish
	// this connection every time. Audit the TFTP library.
	c, err := tftp.NewClient(t.opts...)
	if err != nil {
		return nil, err
	}

	r, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}
	return uio.NewCachingReader(r), nil
}

// SchemeWithRetries wraps a FileScheme and automatically retries (with
// backoff) when Fetch returns a non-nil err.
type SchemeWithRetries struct {
	Scheme  FileScheme
	BackOff backoff.BackOff
}

// Fetch implements FileScheme.Fetch.
func (s *SchemeWithRetries) Fetch(u *url.URL) (io.ReaderAt, error) {
	var err error
	s.BackOff.Reset()
	for d := time.Duration(0); d != backoff.Stop; d = s.BackOff.NextBackOff() {
		if d > 0 {
			time.Sleep(d)
		}

		var r io.ReaderAt
		r, err = s.Scheme.Fetch(u)
		if err != nil {
			log.Printf("Error: Getting %v: %v", u, err)
			continue
		}
		return r, nil
	}

	log.Printf("Error: Too many retries to get file %v", u)
	return nil, err
}

// HTTPClient implements FileScheme for HTTP files.
type HTTPClient struct {
	c *http.Client
}

// NewHTTPClient returns a new HTTP FileScheme based on the given http.Client.
func NewHTTPClient(c *http.Client) *HTTPClient {
	return &HTTPClient{
		c: c,
	}
}

// Fetch implements FileScheme.Fetch.
func (h HTTPClient) Fetch(u *url.URL) (io.ReaderAt, error) {
	resp, err := h.c.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP server responded with code %d, want 200: response %v", resp.StatusCode, resp)
	}
	return uio.NewCachingReader(resp.Body), nil
}

// HTTPClientWithRetries implements FileScheme for HTTP files and automatically
// retries (with backoff) upon an error.
type HTTPClientWithRetries struct {
	Client  *http.Client
	BackOff backoff.BackOff
}

// Fetch implements FileScheme.Fetch.
func (h HTTPClientWithRetries) Fetch(u *url.URL) (io.ReaderAt, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	h.BackOff.Reset()
	for d := time.Duration(0); d != backoff.Stop; d = h.BackOff.NextBackOff() {
		if d > 0 {
			time.Sleep(d)
		}

		var resp *http.Response
		// Note: err uses the scope outside the for loop.
		resp, err = h.Client.Do(req)
		if err != nil {
			log.Printf("Error: HTTP client: %v", err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("Error: HTTP server responded with code %d, want 200: response %v", resp.StatusCode, resp)
			continue
		}
		return uio.NewCachingReader(resp.Body), nil
	}
	log.Printf("Error: Too many retries to download %v", u)
	return nil, fmt.Errorf("too many HTTP retries: %v", err)
}

// LocalFileClient implements FileScheme for files on disk.
type LocalFileClient struct{}

// Fetch implements FileScheme.Fetch.
func (lfs LocalFileClient) Fetch(u *url.URL) (io.ReaderAt, error) {
	return os.Open(filepath.Clean(u.Path))
}
