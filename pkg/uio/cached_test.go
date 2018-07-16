// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestCachingReaderRead(t *testing.T) {
	type read struct {
		// Buffer sizes to call Read with.
		size int

		// Buffer value corresponding Read(size) we want.
		want []byte

		// Error corresponding to Read(size) we want.
		err error
	}

	for i, tt := range []struct {
		// Content of the underlying io.Reader.
		content []byte

		// Read calls to make in order.
		reads []read
	}{
		{
			content: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
			reads: []read{
				{
					size: 0,
				},
				{
					size: 1,
					want: []byte{0x11},
				},
				{
					size: 2,
					want: []byte{0x22, 0x33},
				},
				{
					size: 0,
				},
				{
					size: 3,
					want: []byte{0x44, 0x55, 0x66},
				},
				{
					size: 4,
					want: []byte{0x77, 0x88, 0x99},
					err:  io.EOF,
				},
			},
		},
		{
			content: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
			reads: []read{
				{
					size: 11,
					want: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
					err:  io.EOF,
				},
			},
		},
		{
			content: nil,
			reads: []read{
				{
					size: 2,
					err:  io.EOF,
				},
				{
					size: 0,
				},
			},
		},
		{
			content: []byte{0x33, 0x22, 0x11},
			reads: []read{
				{
					size: 3,
					want: []byte{0x33, 0x22, 0x11},
					err:  nil,
				},
				{
					size: 0,
				},
				{
					size: 1,
					err:  io.EOF,
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			buf := NewCachingReader(bytes.NewBuffer(tt.content))
			for j, r := range tt.reads {
				p := make([]byte, r.size)
				m, err := buf.Read(p)
				if err != r.err {
					t.Errorf("Read#%d(%d) = %v, want %v", j, r.size, err, r.err)
				}
				if m != len(r.want) {
					t.Errorf("Read#%d(%d) = %d, want %d", j, r.size, m, len(r.want))
				}
				if !bytes.Equal(r.want, p[:m]) {
					t.Errorf("Read#%d(%d) = %v, want %v", j, r.size, p[:m], r.want)
				}
			}
		})
	}
}

func TestCachingReaderReadAt(t *testing.T) {
	type readAt struct {
		// Buffer sizes to call Read with.
		size int

		// Offset to read from.
		off int64

		// Buffer value corresponding Read(size) we want.
		want []byte

		// Error corresponding to Read(size) we want.
		err error
	}

	for i, tt := range []struct {
		// Content of the underlying io.Reader.
		content []byte

		// Read calls to make in order.
		reads []readAt
	}{
		{
			content: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
			reads: []readAt{
				{
					off:  0,
					size: 0,
				},
				{
					off:  0,
					size: 1,
					want: []byte{0x11},
				},
				{
					off:  1,
					size: 2,
					want: []byte{0x22, 0x33},
				},
				{
					off:  0,
					size: 0,
				},
				{
					off:  3,
					size: 3,
					want: []byte{0x44, 0x55, 0x66},
				},
				{
					off:  6,
					size: 4,
					want: []byte{0x77, 0x88, 0x99},
					err:  io.EOF,
				},
				{
					off:  0,
					size: 9,
					want: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
				},
				{
					off:  0,
					size: 10,
					want: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
					err:  io.EOF,
				},
				{
					off:  0,
					size: 8,
					want: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
				},
			},
		},
		{
			content: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
			reads: []readAt{
				{
					off:  10,
					size: 10,
					err:  io.EOF,
				},
				{
					off:  5,
					size: 4,
					want: []byte{0x66, 0x77, 0x88, 0x99},
				},
			},
		},
		{
			content: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
			reads: []readAt{
				{
					size: 9,
					want: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
				},
				{
					off:  5,
					size: 4,
					want: []byte{0x66, 0x77, 0x88, 0x99},
				},
				{
					off:  9,
					size: 1,
					err:  io.EOF,
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			buf := NewCachingReader(bytes.NewBuffer(tt.content))
			for j, r := range tt.reads {
				p := make([]byte, r.size)
				m, err := buf.ReadAt(p, r.off)
				if err != r.err {
					t.Errorf("Read#%d(%d) = %v, want %v", j, r.size, err, r.err)
				}
				if m != len(r.want) {
					t.Errorf("Read#%d(%d) = %d, want %d", j, r.size, m, len(r.want))
				}
				if !bytes.Equal(r.want, p[:m]) {
					t.Errorf("Read#%d(%d) = %v, want %v", j, r.size, p[:m], r.want)
				}
			}
		})
	}
}
