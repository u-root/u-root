// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pez

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/klauspost/compress/zstd"
	"github.com/u-root/uio/uio"
)

func newTestImage(t *testing.T, hdr Header, payload []byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, hdr); err != nil {
		t.Fatalf("failed to encode header %v: %v", hdr, err)
	}
	if payload != nil {
		buf.Write(payload)
	}
	return buf.Bytes()
}

func TestExtract(t *testing.T) {
	t.Parallel()

	headerSize := binary.Size(Header{})

	testPayload := []byte("test payload")

	var zstdPayload bytes.Buffer
	zstdEncoder, err := zstd.NewWriter(&zstdPayload)
	if err != nil {
		t.Fatalf("failed to create zstd writer: %v", err)
	}
	if _, err := zstdEncoder.Write(testPayload); err != nil {
		t.Fatalf("failed to write to zstd writer: %v", err)
	}
	if err := zstdEncoder.Close(); err != nil {
		t.Fatalf("failed to close zstd writer: %v", err)
	}
	zstdImageHeader := Header{
		Magic:       magic,
		Type:        typeZImage,
		Compression: [4]byte{0x7a, 0x73, 0x74, 0x64}, // zstd
		Offset:      uint32(headerSize),
		Size:        uint32(zstdPayload.Len()),
	}
	zstdImage := newTestImage(t, zstdImageHeader, zstdPayload.Bytes())

	zstdImageBrokenHeader := Header{
		Magic:       magic,
		Type:        typeZImage,
		Compression: [4]byte{0x7a, 0x73, 0x74, 0x64}, // zstd
		Offset:      uint32(headerSize),
		Size:        uint32(zstdPayload.Len()) / 2,
	}
	zstdImageBroken := newTestImage(t, zstdImageBrokenHeader, zstdPayload.Bytes())

	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr error
	}{
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: io.EOF,
		},
		{
			name:    "magic mismatch",
			input:   newTestImage(t, Header{Magic: 0x1234}, nil),
			wantErr: ErrMagicMismatch,
		},
		{
			name:    "not zimage",
			input:   newTestImage(t, Header{Magic: magic, Type: 0x1234}, nil),
			wantErr: ErrNotZImage,
		},
		{
			name:    "unsupported compression",
			input:   newTestImage(t, Header{Magic: magic, Type: typeZImage, Compression: [4]byte{1, 2, 3, 4}}, nil),
			wantErr: ErrUnsupportedCompression,
		},
		{
			name:  "zstd compression",
			input: zstdImage,
			want:  testPayload,
		},
		{
			name:    "zstd compression broken",
			input:   zstdImageBroken,
			wantErr: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payload, err := Extract(bytes.NewReader(tt.input))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Extract(%x), unexpected error: %v, want: %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			got, err := io.ReadAll(uio.Reader(payload))
			if err != nil {
				t.Fatalf("Could not read kernel from loaded image: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Unexpected diff (-want +got) of extracted payload: %s", diff)
			}
		})
	}
}
