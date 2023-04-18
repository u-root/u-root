// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
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

func (m *mockReader) ReadAt([]byte, int64) (int, error) {
	m.called = true
	return 0, m.err
}

func TestLazyOpenerRead(t *testing.T) {
	for i, tt := range []struct {
		openErr    error
		reader     *mockReader
		wantCalled bool
	}{
		{
			openErr:    nil,
			reader:     &mockReader{},
			wantCalled: true,
		},
		{
			openErr:    io.EOF,
			reader:     nil,
			wantCalled: false,
		},
		{
			openErr: nil,
			reader: &mockReader{
				err: io.ErrUnexpectedEOF,
			},
			wantCalled: true,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			var opened bool
			lr := NewLazyOpener("testname", func() (io.Reader, error) {
				opened = true
				return tt.reader, tt.openErr
			})
			_, err := lr.Read([]byte{})
			if !opened {
				t.Fatalf("Read(): Reader was not opened")
			}
			if tt.openErr != nil && err != tt.openErr {
				t.Errorf("Read() = %v, want %v", err, tt.openErr)
			}
			if tt.reader != nil {
				if got, want := tt.reader.called, tt.wantCalled; got != want {
					t.Errorf("mockReader.Read() called is %v, want %v", got, want)
				}
				if tt.reader.err != nil && err != tt.reader.err {
					t.Errorf("Read() = %v, want %v", err, tt.reader.err)
				}
			}
		})
	}
}

func TestLazyOpenerReadAt(t *testing.T) {
	for i, tt := range []struct {
		limit   int64
		bufSize int
		openErr error
		reader  io.ReaderAt
		off     int64
		want    error
		wantB   []byte
	}{
		{
			limit:   -1,
			bufSize: 10,
			openErr: nil,
			reader:  &mockReader{},
		},
		{
			limit:   -1,
			bufSize: 10,
			openErr: io.EOF,
			reader:  nil,
			want:    io.EOF,
		},
		{
			limit:   -1,
			bufSize: 10,
			openErr: nil,
			reader: &mockReader{
				err: io.ErrUnexpectedEOF,
			},
			want: io.ErrUnexpectedEOF,
		},
		{
			limit:   -1,
			bufSize: 6,
			reader:  strings.NewReader("foobar"),
			wantB:   []byte("foobar"),
		},
		{
			limit:   -1,
			off:     3,
			bufSize: 3,
			reader:  strings.NewReader("foobar"),
			wantB:   []byte("bar"),
		},
		{
			limit:   5,
			off:     3,
			bufSize: 3,
			reader:  strings.NewReader("foobar"),
			wantB:   []byte("ba"),
		},
		{
			limit:   2,
			bufSize: 2,
			reader:  strings.NewReader("foobar"),
			wantB:   []byte("fo"),
		},
		{
			limit:  2,
			off:    2,
			reader: strings.NewReader("foobar"),
			want:   io.EOF,
		},
	} {
		t.Run(fmt.Sprintf("Test #%02d", i), func(t *testing.T) {
			var opened bool
			lr := NewLazyLimitOpenerAt("", tt.limit, func() (io.ReaderAt, error) {
				opened = true
				return tt.reader, tt.openErr
			})

			b := make([]byte, tt.bufSize)
			n, err := lr.ReadAt(b, tt.off)
			if !opened {
				t.Fatalf("Read(): Reader was not opened")
			}
			if !errors.Is(tt.want, err) {
				t.Errorf("Read() = %v, want %v", err, tt.want)
			}

			if err == nil {
				if !bytes.Equal(b[:n], tt.wantB) {
					t.Errorf("Read() = %s, want %s", b[:n], tt.wantB)
				}
			}
		})
	}
}
