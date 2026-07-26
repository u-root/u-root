// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pez

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

func newTestPEImage(t *testing.T, sectionName string, payload []byte, sectionOffset uint32) []byte {
	t.Helper()

	buf := make([]byte, sectionOffset+uint32(len(payload)))

	// 1. MZ magic at offset 0
	buf[0] = 'M'
	buf[1] = 'Z'

	// 2. PE header offset at 0x3c
	peOffset := uint32(64)
	binary.LittleEndian.PutUint32(buf[0x3c:0x3c+4], peOffset)

	// 3. PE signature at peOffset
	copy(buf[peOffset:peOffset+4], []byte("PE\x00\x00"))

	// 4. COFF header
	// Number of sections is 1
	binary.LittleEndian.PutUint16(buf[peOffset+6:peOffset+8], 1)
	// Size of optional header is 0
	binary.LittleEndian.PutUint16(buf[peOffset+20:peOffset+22], 0)

	// 5. Section Table starts at peOffset + 24
	sectionTableOffset := peOffset + 24

	// Write name (8 bytes)
	copy(buf[sectionTableOffset:sectionTableOffset+8], []byte(sectionName))

	// Write Size of Raw Data at offset 16 of section header
	binary.LittleEndian.PutUint32(buf[sectionTableOffset+16:sectionTableOffset+20], uint32(len(payload)))

	// Write Pointer to Raw Data (Raw Offset) at offset 20 of section header
	binary.LittleEndian.PutUint32(buf[sectionTableOffset+20:sectionTableOffset+24], sectionOffset)

	// 6. Write payload at sectionOffset
	copy(buf[sectionOffset:], payload)

	return buf
}

type errReaderAt struct {
	err error
}

func (e errReaderAt) ReadAt(p []byte, off int64) (int, error) {
	return 0, e.err
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

	peImageWithBadOffset := newTestPEImage(t, ".linux", zstdImage, 128)
	binary.LittleEndian.PutUint32(peImageWithBadOffset[0x3c:0x3c+4], 1000)

	peImageWithBadSig := newTestPEImage(t, ".linux", zstdImage, 128)
	copy(peImageWithBadSig[64:68], []byte("PD\x00\x00"))

	tests := []struct {
		name    string
		input   io.ReaderAt
		want    []byte
		wantErr error
	}{
		{
			name:    "empty input",
			input:   bytes.NewReader([]byte{}),
			wantErr: io.EOF,
		},
		{
			name:    "magic mismatch",
			input:   bytes.NewReader(newTestImage(t, Header{Magic: 0x1234}, nil)),
			wantErr: ErrMagicMismatch,
		},
		{
			name:    "not zimage",
			input:   bytes.NewReader(newTestImage(t, Header{Magic: magic, Type: 0x1234}, nil)),
			wantErr: ErrNotZImage,
		},
		{
			name:    "unsupported compression",
			input:   bytes.NewReader(newTestImage(t, Header{Magic: magic, Type: typeZImage, Compression: [4]byte{1, 2, 3, 4}}, nil)),
			wantErr: ErrUnsupportedCompression,
		},
		{
			name:  "zstd compression",
			input: bytes.NewReader(zstdImage),
			want:  testPayload,
		},
		{
			name:    "zstd compression broken",
			input:   bytes.NewReader(zstdImageBroken),
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:  "PE image with .linux section",
			input: bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)),
			want:  testPayload,
		},
		{
			name:    "PE image with wrong section name",
			input:   bytes.NewReader(newTestPEImage(t, ".other", zstdImage, 128)),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image with bad PE offset",
			input:   bytes.NewReader(peImageWithBadOffset),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image with bad PE signature",
			input:   bytes.NewReader(peImageWithBadSig),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image truncated during section reading",
			input:   bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)[:100]),
			wantErr: ErrNotZImage,
		},
		{
			name:    "simulated read error in readAt",
			input:   errReaderAt{err: fmt.Errorf("read error")},
			wantErr: fmt.Errorf("read error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			payload, err := Extract(tt.input)
			if !errors.Is(err, tt.wantErr) && err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Extract(), unexpected error: %v, want: %v", err, tt.wantErr)
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
