// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curl

import (
	"context"
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

	// nextErr is the error to return for the next nextErrCount calls to
	// Fetch. Note this introduces state into the MockScheme which is only
	// okay in this scenario because MockScheme is only used for testing.
	nextErr      error
	nextErrCount int
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

// SetErr sets the error which is returned on the next count calls to Fetch.
func (m *MockScheme) SetErr(err error, count int) {
	m.nextErr = err
	m.nextErrCount = count
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

func mockFetch(m *MockScheme, u *url.URL) (*strings.Reader, error) {
	url := u.String()
	m.numCalled[url]++

	if u.Scheme != m.Scheme {
		return nil, ErrWrongScheme
	}

	if m.nextErrCount > 0 {
		m.nextErrCount--
		return nil, m.nextErr
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

// Fetch implements FileScheme.Fetch.
func (m *MockScheme) Fetch(ctx context.Context, u *url.URL) (io.ReaderAt, error) {
	return mockFetch(m, u)
}

// FetchWithoutCache implements FileScheme.FetchWithoutCache.
func (m *MockScheme) FetchWithoutCache(ctx context.Context, u *url.URL) (io.Reader, error) {
	return mockFetch(m, u)
}
