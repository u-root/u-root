// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

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
	tmode   TransferMode
}

func (r *readRequestMock) Addr() *net.UDPAddr          { return r.addr }
func (r *readRequestMock) Name() string                { return r.name }
func (r *readRequestMock) Write(p []byte) (int, error) { return r.writer.Write(p) }
func (r *readRequestMock) WriteSize(i int64)           { r.size = &i }
func (r *readRequestMock) WriteError(c ErrorCode, m string) {
	r.errCode = c
	r.errMsg = m
}
func (r *readRequestMock) TransferMode() TransferMode { return r.tmode }

func TestFileServer_ServeTFTP(t *testing.T) {
	text := getTestData(t, "text")

	cases := []struct {
		name    string
		reqName string

		expectedData      []byte
		expectedSize      *int64
		expectedErrorCode ErrorCode
		expectedErrorMsg  string
	}{
		{
			name:    "file exists",
			reqName: "text",

			expectedData: text,
			expectedSize: ptrInt64(int64(len(text))),
		},
		{
			name:    "file does not exist",
			reqName: "other",

			expectedErrorCode: ErrCodeFileNotFound,
			expectedErrorMsg:  `File "other" does not exist`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs := FileServer("testdata")

			req := readRequestMock{name: c.reqName}

			fs.ServeTFTP(&req)

			// Data
			if !reflect.DeepEqual(c.expectedData, req.writer.Bytes()) {
				t.Errorf("expected data to be %s, but it was %s", c.expectedData, req.writer.String())
			}

			// Size
			if !reflect.DeepEqual(c.expectedSize, req.size) {
				if c.expectedSize == nil || req.size == nil {
					t.Errorf("expected size to be %v, but it was %v", c.expectedSize, req.size)
				} else {
					t.Errorf("expected size to be %v, but it was %v", *c.expectedSize, *req.size)
				}
			}

			// Error Code
			if c.expectedErrorCode != req.errCode {
				t.Errorf("expected error code to be %s, but it was %s", c.expectedErrorCode, req.errCode)
			}

			// Error Message
			if c.expectedErrorMsg != req.errMsg {
				t.Errorf("expected error msg to be %q, but it was %q", c.expectedErrorMsg, req.errMsg)
			}
		})
	}
}

type writeRequestMock struct {
	addr    *net.UDPAddr
	name    string
	reader  bytes.Buffer
	errCode ErrorCode
	errMsg  string
	size    *int64
	tmode   TransferMode
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
func (r *writeRequestMock) TransferMode() TransferMode { return r.tmode }

func TestFileServer_ReceiveTFTP(t *testing.T) {
	text := getTestData(t, "text")

	cases := []struct {
		name    string
		reqName string
		data    []byte

		expectedFilename  string
		expectedData      []byte
		expectedErrorCode ErrorCode
		expectedErrorMsg  string
	}{
		{
			name:    "success",
			reqName: "text",
			data:    text,

			expectedData: text,
		},
		{
			name:    "fail",
			reqName: "",

			expectedData:      []byte{},
			expectedErrorCode: ErrCodeAccessViolation,
			expectedErrorMsg:  `Cannot create file "."`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatal(err)
			}
			fs := FileServer(dir)

			req := writeRequestMock{name: c.reqName}
			req.reader.Write(c.data)

			fs.ReceiveTFTP(&req)

			// Data
			data, _ := ioutil.ReadFile(filepath.Join(dir, c.reqName))
			if !reflect.DeepEqual(c.expectedData, data) {
				t.Errorf("expected data to be %s, but it was %s", c.expectedData, data)
			}

			// Error Code
			if c.expectedErrorCode != req.errCode {
				t.Errorf("expected error code to be %s, but it was %s", c.expectedErrorCode, req.errCode)
			}

			// Error Message
			if c.expectedErrorMsg != req.errMsg {
				t.Errorf("expected error msg to be %q, but it was %q", c.expectedErrorMsg, req.errMsg)
			}
		})
	}
}
