package pxe

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
	// returning an io.Reader that downloads the file.
	GetFile(u *url.URL) (io.Reader, error)
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
		"tftp": NewCachedFileScheme(DefaultTFTPClient),
		"http": NewCachedFileScheme(DefaultHTTPClient),
		"file": NewCachedFileScheme(&LocalFileClient{}),
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
func GetFile(u *url.URL) (io.Reader, error) {
	return DefaultSchemes.GetFile(u)
}

// GetFile downloads the file with the given `u`. `u.Scheme` is used to
// select the FileScheme via `s`.
//
// If `s` does not contain a FileScheme for `u.Scheme`, ErrNoSuchScheme is
// returned.
func (s Schemes) GetFile(u *url.URL) (io.Reader, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}
	r, err := fg.GetFile(u)
	if err != nil {
		return nil, &URLError{URL: u, Err: err}
	}
	return r, nil
}

// LazyGetFile calls LazyGetFile on DefaultSchemes. See Schemes.LazyGetFile.
func LazyGetFile(u *url.URL) (io.Reader, error) {
	return DefaultSchemes.LazyGetFile(u)
}

// LazyGetFile returns a reader that will download the file given by `u` when
// Read is called, based on `u`s scheme. See Schemes.GetFile for more
// details.
func (s Schemes) LazyGetFile(u *url.URL) (io.Reader, error) {
	fg, ok := s[u.Scheme]
	if !ok {
		return nil, &URLError{URL: u, Err: ErrNoSuchScheme}
	}

	return NewLazyOpener(func() (io.Reader, error) {
		r, err := fg.GetFile(u)
		if err != nil {
			return nil, &URLError{URL: u, Err: err}
		}
		return r, nil
	}), nil
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
func (t *TFTPClient) GetFile(u *url.URL) (io.Reader, error) {
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
	return r, nil
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
func (h HTTPClient) GetFile(u *url.URL) (io.Reader, error) {
	resp, err := h.c.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP server responded with code %d, want 200: response %v", resp.StatusCode, resp)
	}
	return resp.Body, nil
}

// LocalFileClient implements FileScheme for files on disk.
type LocalFileClient struct{}

// GetFile implements FileScheme.GetFile.
func (lfs LocalFileClient) GetFile(u *url.URL) (io.Reader, error) {
	return os.Open(filepath.Clean(u.Path))
}

type cachedFile struct {
	cr  *CachingReader
	err error
}

// CachedFileScheme implements FileScheme and caches files downloaded from a
// FileScheme.
type CachedFileScheme struct {
	fs FileScheme

	// cache is a map of URL string -> cached file or error object.
	cache map[string]cachedFile
}

// NewCachedFileScheme returns a caching wrapper for the given FileScheme `fs`.
func NewCachedFileScheme(fs FileScheme) FileScheme {
	return &CachedFileScheme{
		fs:    fs,
		cache: make(map[string]cachedFile),
	}
}

// GetFile implements FileScheme.GetFile.
func (cc *CachedFileScheme) GetFile(u *url.URL) (io.Reader, error) {
	url := u.String()
	if cf, ok := cc.cache[url]; ok {
		// File is in cache.
		if cf.err != nil {
			return nil, cf.err
		}
		return cf.cr.NewReader(), nil
	}

	r, err := cc.fs.GetFile(u)
	if err != nil {
		cc.cache[url] = cachedFile{err: err}
		return nil, err
	}
	cr := NewCachingReader(r)
	cc.cache[url] = cachedFile{cr: cr}
	return cr.NewReader(), nil
}
