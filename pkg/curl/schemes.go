// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package curl implements routines to fetch files given a URL.
//
// curl currently supports HTTP, TFTP, and local files.
package curl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

// FileSchemeRetryFilter contains extra RetryFilter method for a FileScheme
// wrapped by SchemeWithRetries.
type FileSchemeRetryFilter interface {
	// RetryFilter lets a FileScheme filter for errors returned by Fetch
	// which are worth retrying. If this interface is not implemented, the
	// default for SchemeWithRetries is to always retry. RetryFilter
	// returns true to indicate a request should be retried.
	RetryFilter(u *url.URL, err error) bool
}

var (
	// DefaultHTTPClient is the default HTTP FileScheme.
	//
	// It is not recommended to use this for HTTPS. We recommend creating an
	// http.Client that accepts only a private pool of certificates.
	DefaultHTTPClient = NewHTTPClient(http.DefaultClient)

	// DefaultTFTPClient is the default TFTP FileScheme.
	DefaultTFTPClient = NewTFTPClient(tftp.ClientMode(tftp.ModeOctet), tftp.ClientBlocksize(1450), tftp.ClientWindowsize(65535))

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

// RetryFilter implements FileSchemeRetryFilter.
func (s Schemes) RetryFilter(u *url.URL, err error) bool {
	fg, ok := s[u.Scheme]
	if !ok {
		return false
	}
	if fg, ok := fg.(FileSchemeRetryFilter); ok {
		return fg.RetryFilter(u, err)
	}
	return true
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
		ReaderAt: uio.NewLazyOpenerAt(u.String(), func() (io.ReaderAt, error) {
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

// RetryFilter implements FileSchemeRetryFilter.
func (t *TFTPClient) RetryFilter(u *url.URL, err error) bool {
	// The tftp does not export the necessary structs to get the
	// code out of the error message cleanly.
	return !strings.Contains(err.Error(), "FILE_NOT_FOUND")
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
		// Note: err uses the scope outside the for loop.
		r, err = s.Scheme.Fetch(u)
		if err == nil {
			return r, nil
		}

		log.Printf("Error: Getting %v: %v", u, err)
		if s, ok := s.Scheme.(FileSchemeRetryFilter); ok && !s.RetryFilter(u, err) {
			return r, err
		}
		log.Printf("Retrying %v", u)
	}

	log.Printf("Error: Too many retries to get file %v", u)
	return nil, err
}

// HTTPClientCodeError is returned by HTTPClient.Fetch when the server replies
// with a non-200 code.
type HTTPClientCodeError struct {
	Err      error
	HTTPCode int
}

// Error implements error for HTTPClientCodeError.
func (h *HTTPClientCodeError) Error() string {
	return fmt.Sprintf("HTTP server responded with error code %d, want 200: response %v", h.HTTPCode, h.Err)
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
		return nil, &HTTPClientCodeError{err, resp.StatusCode}
	}
	return uio.NewCachingReader(resp.Body), nil
}

// RetryFilter implements FileSchemeRetryFilter.
func (h HTTPClient) RetryFilter(u *url.URL, err error) bool {
	if err, ok := err.(*HTTPClientCodeError); ok && err.HTTPCode == 200 {
		return false
	}
	return true
}

// LocalFileClient implements FileScheme for files on disk.
type LocalFileClient struct{}

// Fetch implements FileScheme.Fetch.
func (lfs LocalFileClient) Fetch(u *url.URL) (io.ReaderAt, error) {
	return os.Open(filepath.Clean(u.Path))
}
