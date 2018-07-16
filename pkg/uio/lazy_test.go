// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"fmt"
	"io"
	"testing"
)

type mockReader struct {
	// called is whether Read has been called.
	called bool

	// err is the error to return on Read.
	err error
}

func (m *mockReader) Read([]byte) (int, error) {
	m.called = true
	return 0, m.err
}

func TestLazyOpenerRead(t *testing.T) {
	for i, tt := range []struct {
		openErr    error
		openReader *mockReader
		wantCalled bool
	}{
		{
			openErr:    nil,
			openReader: &mockReader{},
			wantCalled: true,
		},
		{
			openErr:    io.EOF,
			openReader: nil,
			wantCalled: false,
		},
		{
			openErr: nil,
			openReader: &mockReader{
				err: io.ErrUnexpectedEOF,
			},
			wantCalled: true,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			var opened bool
			lr := NewLazyOpener(func() (io.Reader, error) {
				opened = true
				return tt.openReader, tt.openErr
			})
			_, err := lr.Read([]byte{})
			if !opened {
				t.Fatalf("Read(): Reader was not opened")
			}
			if tt.openErr != nil && err != tt.openErr {
				t.Errorf("Read() = %v, want %v", err, tt.openErr)
			}
			if tt.openReader != nil {
				if got, want := tt.openReader.called, tt.wantCalled; got != want {
					t.Errorf("mockReader.Read() called is %v, want %v", got, want)
				}
				if tt.openReader.err != nil && err != tt.openReader.err {
					t.Errorf("Read() = %v, want %v", err, tt.openReader.err)
				}
			}
		})
	}
}
