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
	// implementation for the given URI scheme.
	ErrNoSuchScheme = errors.New("no such scheme")
)

// FileScheme represents the implementation of a URI scheme and gives access to
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
	DefaultTFTPClient FileScheme

	// DefaultSchemes are the schemes supported by PXE by default.
	DefaultSchemes Schemes
)

func init() {
	c, err := tftp.NewClient()
	if err != nil {
		panic(fmt.Sprintf("tftp.NewClient failed: %v", err))
	}
	DefaultTFTPClient = NewTFTPClient(c)

	DefaultSchemes = Schemes{
		"tftp": NewCachedFileScheme(DefaultTFTPClient),
		"http": NewCachedFileScheme(DefaultHTTPClient),
		"file": NewCachedFileScheme(&LocalFileClient{}),
	}
}

// Schemes is a map of URI scheme identifier -> implementation that can
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

func parseURI(uri string, wd *url.URL) (*url.URL, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if len(u.Scheme) == 0 {
		u.Scheme = wd.Scheme

		if len(u.Host) == 0 {
			// If this is not there, it was likely just a path.
			u.Host = wd.Host
			u.Path = filepath.Join(wd.Path, filepath.Clean(u.Path))
		}
	}
	return u, nil
}

// GetFile downloads a file via DefaultSchemes. See Schemes.GetFile for
// details.
func GetFile(uri string, wd *url.URL) (io.Reader, error) {
	return DefaultSchemes.GetFile(uri, wd)
}

// GetFile downloads the file with the given `uri`. `uri.Scheme` is used to
// select the FileScheme via `s`.
//
// If `s` does not contain a FileScheme for `uri.Scheme`, ErrNoSuchScheme is
// returned.
//
// If `uri` is just a relative path and not a full URI, `wd` is used as the
// "working directory" of that relative path; the resulting URI is roughly
// `path.Join(wd.String(), uri)`.
func (s Schemes) GetFile(uri string, wd *url.URL) (io.Reader, error) {
	u, err := parseURI(uri, wd)
	if err != nil {
		return nil, err
	}

	fg, ok := s[u.Scheme]
	if !ok {
		return nil, ErrNoSuchScheme
	}
	return fg.GetFile(u)
}

// LazyGetFile calls LazyGetFile on DefaultSchemes. See Schemes.LazyGetFile.
func LazyGetFile(uri string, wd *url.URL) (io.Reader, error) {
	return DefaultSchemes.LazyGetFile(uri, wd)
}

// LazyGetFile returns a reader that will download the file given by `uri` when
// Read is called, based on `uri`s scheme. See Schemes.GetFile for more
// details.
func (s Schemes) LazyGetFile(uri string, wd *url.URL) (io.Reader, error) {
	u, err := parseURI(uri, wd)
	if err != nil {
		return nil, err
	}

	fg, ok := s[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("could not get file based on scheme %q: no such scheme registered", u.Scheme)
	}

	return NewLazyOpener(func() (io.Reader, error) {
		return fg.GetFile(u)
	}), nil
}

// TFTPClient implements FileScheme for TFTP files.
type TFTPClient struct {
	c *tftp.Client
}

// NewTFTPClient returns a new TFTP client based on the given tftp.Client.
func NewTFTPClient(c *tftp.Client) FileScheme {
	return &TFTPClient{
		c: c,
	}
}

// GetFile implements FileScheme.GetFile.
func (t *TFTPClient) GetFile(u *url.URL) (io.Reader, error) {
	r, err := t.c.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %v", u, err)
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
		return nil, fmt.Errorf("Could not download file %s: %v", u, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not download file %s: response %v", u, resp)
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

	// cache is a map of URI string -> cached file or error object.
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
	uri := u.String()
	if cf, ok := cc.cache[uri]; ok {
		if cf.err != nil {
			return nil, cf.err
		}
		return cf.cr.NewReader(), nil
	}

	r, err := cc.fs.GetFile(u)
	if err != nil {
		cc.cache[uri] = cachedFile{err: err}
		return nil, err
	}
	cr := NewCachingReader(r)
	cc.cache[uri] = cachedFile{cr: cr}
	return cr.NewReader(), nil
}
