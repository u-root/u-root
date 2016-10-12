// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import (
	"bytes"
	"io/ioutil"
	"net"
	"path/filepath"
	"reflect"
	"testing"
)

type readRequestMock struct {
	addr    *net.UDPAddr
	name    string
	writer  bytes.Buffer
	errCode ErrorCode
	errMsg  string
	size    *int64
}

func (r *readRequestMock) Addr() *net.UDPAddr          { return r.addr }
func (r *readRequestMock) Name() string                { return r.name }
func (r *readRequestMock) Write(p []byte) (int, error) { return r.writer.Write(p) }
func (r *readRequestMock) WriteSize(i int64)           { r.size = &i }
func (r *readRequestMock) WriteError(c ErrorCode, m string) {
	r.errCode = c
	r.errMsg = m
}

func TestFileServer_ServeTFTP(t *testing.T) {
	text := getTestData(t, "text")

	cases := map[string]struct {
		name string

		expectedData      []byte
		expectedSize      *int64
		expectedErrorCode ErrorCode
		expectedErrorMsg  string
	}{
		"file exists": {
			name: "text",

			expectedData: text,
			expectedSize: ptrInt64(int64(len(text))),
		},
		"file does not exist": {
			name: "other",

			expectedErrorCode: ErrCodeFileNotFound,
			expectedErrorMsg:  `File "other" does not exist`,
		},
	}

	for label, c := range cases {
		fs := FileServer("testdata")

		req := readRequestMock{name: c.name}

		fs.ServeTFTP(&req)

		// Data
		if !reflect.DeepEqual(c.expectedData, req.writer.Bytes()) {
			t.Errorf("%s: Expected data to be %s, but it was %s", label, c.expectedData, req.writer.String())
		}

		// Size
		if !reflect.DeepEqual(c.expectedSize, req.size) {
			if c.expectedSize == nil || req.size == nil {
				t.Errorf("%s: Expected size to be %v, but it was %v", label, c.expectedSize, req.size)
			} else {
				t.Errorf("%s: Expected size to be %v, but it was %v", label, *c.expectedSize, *req.size)
			}
		}

		// Error Code
		if c.expectedErrorCode != req.errCode {
			t.Errorf("%s: Expected error code to be %s, but it was %s", label, c.expectedErrorCode, req.errCode)
		}

		// Error Message
		if c.expectedErrorMsg != req.errMsg {
			t.Errorf("%s: Expected error msg to be %q, but it was %q", label, c.expectedErrorMsg, req.errMsg)
		}
	}
}

type writeRequestMock struct {
	addr    *net.UDPAddr
	name    string
	reader  bytes.Buffer
	errCode ErrorCode
	errMsg  string
	size    *int64
}

func (r *writeRequestMock) Addr() *net.UDPAddr         { return r.addr }
func (r *writeRequestMock) Name() string               { return r.name }
func (r *writeRequestMock) Read(p []byte) (int, error) { return r.reader.Read(p) }
func (r *writeRequestMock) Size() (int64, error) {
	if r.size != nil {
		return *r.size, nil
	}
	return 0, ErrSizeNotReceived
}
func (r *writeRequestMock) WriteError(c ErrorCode, m string) {
	r.errCode = c
	r.errMsg = m
}

func TestFileServer_ReceiveTFTP(t *testing.T) {
	text := getTestData(t, "text")

	cases := map[string]struct {
		name string
		data []byte

		expectedFilename  string
		expectedData      []byte
		expectedErrorCode ErrorCode
		expectedErrorMsg  string
	}{
		"success": {
			name: "text",
			data: text,

			expectedData: text,
		},
		"fail": {
			name: "",

			expectedData:      []byte{},
			expectedErrorCode: ErrCodeAccessViolation,
			expectedErrorMsg:  `Cannot create file "."`,
		},
	}

	for label, c := range cases {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatal(err)
		}
		fs := FileServer(dir)

		req := writeRequestMock{name: c.name}
		req.reader.Write(c.data)

		fs.ReceiveTFTP(&req)

		// Data
		data, _ := ioutil.ReadFile(filepath.Join(dir, c.name))
		if !reflect.DeepEqual(c.expectedData, data) {
			t.Errorf("%s: Expected data to be %s, but it was %s", label, c.expectedData, data)
		}

		// Error Code
		if c.expectedErrorCode != req.errCode {
			t.Errorf("%s: Expected error code to be %s, but it was %s", label, c.expectedErrorCode, req.errCode)
		}

		// Error Message
		if c.expectedErrorMsg != req.errMsg {
			t.Errorf("%s: Expected error msg to be %q, but it was %q", label, c.expectedErrorMsg, req.errMsg)
		}
	}
}
