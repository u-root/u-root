// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package urlfetch

import (
	"errors"
	"io"
	"net/url"
	"path"
	"strings"
)

// MockScheme is a Scheme mock for testing.
type MockScheme struct {
	// scheme is the scheme name.
	Scheme string

	// hosts is a map of host -> relative filename to host -> file contents.
	hosts map[string]map[string]string

	// numCalled is a map of URL string -> number of times Fetch has been
	// called on that URL.
	numCalled map[string]uint
}

// NewMockScheme creates a new MockScheme with the given scheme name.
func NewMockScheme(scheme string) *MockScheme {
	return &MockScheme{
		Scheme:    scheme,
		hosts:     make(map[string]map[string]string),
		numCalled: make(map[string]uint),
	}
}

// Add adds a file to the MockScheme
func (m *MockScheme) Add(host string, p string, content string) {
	_, ok := m.hosts[host]
	if !ok {
		m.hosts[host] = make(map[string]string)
	}

	m.hosts[host][path.Clean(p)] = content
}

// NumCalled returns how many times a url has been looked up.
func (m *MockScheme) NumCalled(u *url.URL) uint {
	url := u.String()
	if c, ok := m.numCalled[url]; ok {
		return c
	}
	return 0
}

var (
	// ErrWrongScheme means the wrong mocked scheme was used.
	ErrWrongScheme = errors.New("wrong scheme")
	// ErrNoSuchHost means there is no host record in the mock.
	ErrNoSuchHost = errors.New("no such host exists")
	// ErrNoSuchFile means there is no file record in the mock.
	ErrNoSuchFile = errors.New("no such file exists on this host")
)

// Fetch implements FileScheme.Fetch.
func (m *MockScheme) Fetch(u *url.URL) (io.ReaderAt, error) {
	url := u.String()
	if _, ok := m.numCalled[url]; ok {
		m.numCalled[url]++
	} else {
		m.numCalled[url] = 1
	}

	if u.Scheme != m.Scheme {
		return nil, ErrWrongScheme
	}

	files, ok := m.hosts[u.Host]
	if !ok {
		return nil, ErrNoSuchHost
	}

	content, ok := files[path.Clean(u.Path)]
	if !ok {
		return nil, ErrNoSuchFile
	}
	return strings.NewReader(content), nil
}
