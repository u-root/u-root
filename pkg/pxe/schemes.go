// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pxe

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/u-root/u-root/pkg/uio"
	"pack.ag/tftp"
)

var (
	// ErrNoSuchScheme is returned by Schemes.GetFile and
	// Schemes.LazyGetFile if there is no registered FileScheme
	// implementation for the given URL scheme.
	ErrNoSuchScheme = errors.New("no such scheme")
)

// FileScheme represents the implementation of a URL scheme and gives access to
// downloading files of that scheme.
//
// For example, an http FileScheme implementation would download files using
// the HTTP protocol.
type FileScheme interface {
	// GetFile returns a reader that gives the contents of `u`.
	//
	// It may do so by downloading `u` and placing it in a buffer, or by
	// returning an io.ReaderAt that downloads the file.
	GetFile(u *url.URL) (io.ReaderAt, error)
}

var (
	// DefaultHTTPClient is the default HTTP FileScheme.
	//
	// It is not recommended to use this for HTTPS. We recommend creating an
	// http.Client that accepts only a private pool of certificates.
	DefaultHTTPClient = NewHTTPClient(http.DefaultClient)

	// DefaultTFTPClient is the default TFTP FileScheme.
	DefaultTFTPClient = NewTFTPClient()

	// DefaultSchemes are the schemes supported by PXE by default.
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
// download a file for that scheme.
type Schemes map[string]FileScheme

// RegisterScheme calls DefaultSchemes.Register.
func RegisterScheme(scheme string, fs FileScheme) {
	DefaultSchemes.Register(scheme, fs)
}

// Register registers a scheme identified by `scheme` to be `fs`.
func (s Schemes) Register(scheme string, fs FileScheme) {
	s[scheme] = fs
}

// GetFile downloads a file via DefaultSchemes. See Schemes.GetFile for
// details.
func GetFile(u *url.URL) (io.ReaderAt, error) {
	return DefaultSchemes.GetFile(u)
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

// GetFile downloads the file with the given `u`. `u.Scheme` is used to
// select the FileScheme via `s`.
//
// If `s` does not contain a FileScheme for `u.Scheme`, ErrNoSuchScheme is
// returned.
func (s Schemes) GetFile(u *url.URL) (io.ReaderAt, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}
	r, err := fg.GetFile(u)
	if err != nil {
		return nil, &URLError{URL: u, Err: err}
	}
	return &file{ReaderAt: r, url: u}, nil
}

// LazyGetFile calls LazyGetFile on DefaultSchemes. See Schemes.LazyGetFile.
func LazyGetFile(u *url.URL) (io.ReaderAt, error) {
	return DefaultSchemes.LazyGetFile(u)
}

// LazyGetFile returns a reader that will download the file given by `u` when
// Read is called, based on `u`s scheme. See Schemes.GetFile for more
// details.
func (s Schemes) LazyGetFile(u *url.URL) (io.ReaderAt, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}

	return &file{
		url: u,
		ReaderAt: uio.NewLazyOpenerAt(func() (io.ReaderAt, error) {
			r, err := fg.GetFile(u)
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

// GetFile implements FileScheme.GetFile.
func (t *TFTPClient) GetFile(u *url.URL) (io.ReaderAt, error) {
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

// GetFile implements FileScheme.GetFile.
func (h HTTPClient) GetFile(u *url.URL) (io.ReaderAt, error) {
	resp, err := h.c.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP server responded with code %d, want 200: response %v", resp.StatusCode, resp)
	}
	return uio.NewCachingReader(resp.Body), nil
}

// HTTPClientWithRetries implements FileScheme for HTTP files. It retries
// until a response is received or until the timeout.
type HTTPClientWithRetries struct {
	Client  *http.Client
	Timeout time.Duration
}

// GetFile implements FileScheme.GetFile.
func (h HTTPClientWithRetries) GetFile(u *url.URL) (io.ReaderAt, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	counter := 0
	delay := 500 * time.Millisecond
	maxDelay := 10 * time.Second

	timer := time.NewTimer(h.Timeout)
	defer timer.Stop()
	for {
		resp, err := h.Client.Do(req)
		if err == nil {
			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("HTTP server responded with code %d, want 200: response %v",
					resp.StatusCode, resp)
			}
			return uio.NewCachingReader(resp.Body), nil
		}
		counter++

		// Backoff
		delay += delay / 2
		if delay > maxDelay {
			delay = maxDelay
		}
		select {
		case <-timer.C:
			return nil, fmt.Errorf("http timeout after %d tries: %v", counter, err)
		case <-time.After(delay):
		}
	}
}

// LocalFileClient implements FileScheme for files on disk.
type LocalFileClient struct{}

// GetFile implements FileScheme.GetFile.
func (lfs LocalFileClient) GetFile(u *url.URL) (io.ReaderAt, error) {
	return os.Open(filepath.Clean(u.Path))
}
