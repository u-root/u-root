package pxe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"reflect"
	"testing"
)

type MockScheme struct {
	// scheme is the scheme name.
	scheme string

	// hosts is a map of host -> relative filename to host -> file contents.
	hosts map[string]map[string]string

	// numCalled is a map of URL string -> number of times GetFile has been
	// called on that URL.
	numCalled map[string]uint
}

func NewMockScheme(scheme string) *MockScheme {
	return &MockScheme{
		scheme:    scheme,
		hosts:     make(map[string]map[string]string),
		numCalled: make(map[string]uint),
	}
}

func (m *MockScheme) Add(host string, p string, content string) {
	_, ok := m.hosts[host]
	if !ok {
		m.hosts[host] = make(map[string]string)
	}

	m.hosts[host][path.Clean(p)] = content
}

func (m *MockScheme) NumCalled(u *url.URL) uint {
	url := u.String()
	if c, ok := m.numCalled[url]; ok {
		return c
	}
	return 0
}

var (
	errWrongScheme = errors.New("wrong scheme")
	errNoSuchHost  = errors.New("no such host exists")
	errNoSuchFile  = errors.New("no such file exists on this host")
)

func (m *MockScheme) GetFile(u *url.URL) (io.Reader, error) {
	url := u.String()
	if _, ok := m.numCalled[url]; ok {
		m.numCalled[url]++
	} else {
		m.numCalled[url] = 1
	}

	if u.Scheme != m.scheme {
		return nil, errWrongScheme
	}

	files, ok := m.hosts[u.Host]
	if !ok {
		return nil, errNoSuchHost
	}

	content, ok := files[path.Clean(u.Path)]
	if !ok {
		return nil, errNoSuchFile
	}
	return bytes.NewBufferString(content), nil
}

func TestCachedFileSchemeGetFile(t *testing.T) {
	for i, tt := range []struct {
		fs   func() *MockScheme
		url  *url.URL
		err  error
		want string
	}{
		{
			fs: func() *MockScheme {
				s := NewMockScheme("fooftp")
				s.Add("192.168.0.1", "/default", "haha")
				return s
			},
			url: &url.URL{
				Scheme: "fooftp",
				Host:   "192.168.0.1",
				Path:   "/default",
			},
			want: "haha",
		},
		{
			fs: func() *MockScheme {
				return NewMockScheme("fooftp")
			},
			url: &url.URL{
				Scheme: "fooftp",
			},
			err: errNoSuchHost,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			ms := tt.fs()
			fs := NewCachedFileScheme(ms)
			r, err := fs.GetFile(tt.url)
			if err != tt.err {
				t.Errorf("GetFile(%s) = %v, want %v", tt.url, err, tt.err)
				return
			} else if err == nil {
				content, err := ioutil.ReadAll(r)
				if err != nil {
					t.Errorf("ReadAll = %v, want nil", err)
				}
				if got := string(content); got != tt.want {
					t.Errorf("Read(%s) got %v, want %v", tt.url, got, tt.want)
				}
			}

			r2, err2 := fs.GetFile(tt.url)
			if err2 != tt.err {
				t.Errorf("GetFile2(%s) = %v, want %v", tt.url, err2, tt.err)
				return
			} else if err2 == nil {
				content2, err := ioutil.ReadAll(r2)
				if err != nil {
					t.Errorf("ReadAll2 = %v, want nil", err)
				}
				if got := string(content2); got != tt.want {
					t.Errorf("Read2(%s) got %v, want %v", tt.url, got, tt.want)
				}
			}

			if got := ms.NumCalled(tt.url); got != 1 {
				t.Errorf("num called(%s) = %d, want 1", tt.url, got)
			}
		})
	}
}

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
			err: errNoSuchHost,
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
			err: errNoSuchFile,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			fs := tt.scheme()
			s := make(Schemes)
			s.Register(fs.scheme, fs)

			// Test both GetFile and LazyGetFile.
			for _, f := range []func(url *url.URL) (io.Reader, error){
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
				content, err := ioutil.ReadAll(r)
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
