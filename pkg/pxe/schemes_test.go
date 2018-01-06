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

	// numCalled is a map of URI string -> number of times GetFile has been
	// called on that URI.
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
	uri := u.String()
	if c, ok := m.numCalled[uri]; ok {
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
	uri := u.String()
	if _, ok := m.numCalled[uri]; ok {
		m.numCalled[uri]++
	} else {
		m.numCalled[uri] = 1
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
		uri  *url.URL
		err  error
		want string
	}{
		{
			fs: func() *MockScheme {
				s := NewMockScheme("fooftp")
				s.Add("192.168.0.1", "/default", "haha")
				return s
			},
			uri: &url.URL{
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
			uri: &url.URL{
				Scheme: "fooftp",
			},
			err: errNoSuchHost,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			ms := tt.fs()
			fs := NewCachedFileScheme(ms)
			r, err := fs.GetFile(tt.uri)
			if err != tt.err {
				t.Errorf("GetFile(%s) = %v, want %v", tt.uri, err, tt.err)
				return
			} else if err == nil {
				content, err := ioutil.ReadAll(r)
				if err != nil {
					t.Errorf("ReadAll = %v, want nil", err)
				}
				if got := string(content); got != tt.want {
					t.Errorf("Read(%s) got %v, want %v", tt.uri, got, tt.want)
				}
			}

			r2, err2 := fs.GetFile(tt.uri)
			if err2 != tt.err {
				t.Errorf("GetFile2(%s) = %v, want %v", tt.uri, err2, tt.err)
				return
			} else if err2 == nil {
				content2, err := ioutil.ReadAll(r2)
				if err != nil {
					t.Errorf("ReadAll2 = %v, want nil", err)
				}
				if got := string(content2); got != tt.want {
					t.Errorf("Read2(%s) got %v, want %v", tt.uri, got, tt.want)
				}
			}

			if got := ms.NumCalled(tt.uri); got != 1 {
				t.Errorf("num called(%s) = %d, want 1", tt.uri, got)
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	for i, tt := range []struct {
		scheme func() *MockScheme
		wd     *url.URL
		uri    string
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
			uri:  "default",
			wd: &url.URL{
				Scheme: "fooftp",
				Host:   "192.168.0.1",
				Path:   "/foo/pxelinux.cfg",
			},
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				return s
			},
			uri: "nosuch://scheme/foo",
			err: ErrNoSuchScheme,
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				return s
			},
			uri: "fooftp://someotherplace",
			err: errNoSuchHost,
		},
		{
			scheme: func() *MockScheme {
				s := NewMockScheme("fooftp")
				s.Add("somehost", "somefile", "somecontent")
				return s
			},
			uri: "fooftp://somehost/someotherfile",
			err: errNoSuchFile,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			fs := tt.scheme()
			s := make(Schemes)
			s.Register(fs.scheme, fs)

			// Test both GetFile and LazyGetFile.
			for _, f := range []func(uri string, wd *url.URL) (io.Reader, error){
				s.GetFile,
				s.LazyGetFile,
			} {
				r, err := f(tt.uri, tt.wd)
				if got, want := err, tt.err; got != want {
					t.Errorf("GetFile() = %v, want %v", got, want)
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

func TestParseURI(t *testing.T) {
	for i, tt := range []struct {
		uri  string
		wd   *url.URL
		err  bool
		want *url.URL
	}{
		{
			uri: "default",
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
			uri: "http://192.168.2.1/configs/your-machine.cfg",
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
			got, err := parseURI(tt.uri, tt.wd)
			if (err != nil) != tt.err {
				t.Errorf("Wanted error (%v), but got %v", tt.err, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseURI() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
