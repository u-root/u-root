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

type errorAfterReaderAt struct {
	data      []byte
	errOffset int
	err       error
}

func (e *errorAfterReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(e.errOffset) {
		return 0, e.err
	}
	n := copy(p, e.data[off:])
	if off+int64(n) >= int64(e.errOffset) {
		return n, nil
	}
	return n, nil
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

	// Construct a multi-frame zstd image to test read errors during streaming decompression (io.ReadAll)
	var frame1 bytes.Buffer
	w1, err := zstd.NewWriter(&frame1)
	if err != nil {
		t.Fatalf("failed to create zstd writer for frame 1: %v", err)
	}
	if _, err := w1.Write([]byte("hello ")); err != nil {
		t.Fatalf("failed to write to frame 1: %v", err)
	}
	if err := w1.Close(); err != nil {
		t.Fatalf("failed to close frame 1: %v", err)
	}

	var frame2 bytes.Buffer
	w2, err := zstd.NewWriter(&frame2)
	if err != nil {
		t.Fatalf("failed to create zstd writer for frame 2: %v", err)
	}
	if _, err := w2.Write([]byte("world")); err != nil {
		t.Fatalf("failed to write to frame 2: %v", err)
	}
	if err := w2.Close(); err != nil {
		t.Fatalf("failed to close frame 2: %v", err)
	}

	multiFramePayload := append(frame1.Bytes(), frame2.Bytes()...)
	multiFrameHeader := Header{
		Magic:       magic,
		Type:        typeZImage,
		Compression: [4]byte{0x7a, 0x73, 0x74, 0x64}, // zstd
		Offset:      uint32(headerSize),
		Size:        uint32(len(multiFramePayload)),
	}
	multiFrameImage := newTestImage(t, multiFrameHeader, multiFramePayload)

	// Construct a corrupted zstd image where the header is valid but the payload is corrupted
	zstdImageCorrupt := make([]byte, len(zstdImage))
	copy(zstdImageCorrupt, zstdImage)
	// Corrupt the very last byte of the payload
	zstdImageCorrupt[len(zstdImageCorrupt)-1] ^= 0xff

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
			name:    "zstd block corruption during decompression",
			input:   bytes.NewReader(zstdImageCorrupt),
			wantErr: errors.New("CRC check failed"),
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
			name:    "PE image truncated before PE offset",
			input:   bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)[:60]),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image truncated before PE signature",
			input:   bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)[:66]),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image truncated before number of sections",
			input:   bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)[:71]),
			wantErr: ErrNotZImage,
		},
		{
			name:    "PE image truncated before optional header size",
			input:   bytes.NewReader(newTestPEImage(t, ".linux", zstdImage, 128)[:85]),
			wantErr: ErrNotZImage,
		},
		{
			name:    "simulated read error in readAt",
			input:   errReaderAt{err: fmt.Errorf("read error")},
			wantErr: fmt.Errorf("read error"),
		},
		{
			name: "zstd read error during streaming decompression",
			input: &errorAfterReaderAt{
				data:      multiFrameImage,
				errOffset: int(multiFrameHeader.Offset) + len(frame1.Bytes()),
				err:       fmt.Errorf("decompression stream read error"),
			},
			wantErr: fmt.Errorf("decompression stream read error"),
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

// TestHeaderCompressionBoundsCheck serves as an explicit anchor point to ensure compliance
// with current and future Go specifications regarding array and slice access behaviors.
// While testing basic language mechanics may feel like a bit of an overkill or testing more
// than strictly needed, running these tests is exceptionally cheap and, even if marginally,
// contributes to the long-term robustness and stability of the project.
func TestHeaderCompressionBoundsCheck(t *testing.T) {
	// Verify that copying to header.Compression[:] from buf[24:28] is valid for 28 bytes.
	t.Run("size 28", func(t *testing.T) {
		buf := make([]byte, 28)
		copy(buf[24:28], []byte("zstd"))

		var header Header
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Expected no panic, but recovered: %v", r)
			}
		}()

		copy(header.Compression[:], buf[24:28])
		if string(header.Compression[:]) != "zstd" {
			t.Errorf("Expected Compression to be 'zstd', got %q", string(header.Compression[:]))
		}
	})

	// Verify that slicing buf[24:28] panics for an array/slice size of 27.
	t.Run("size 27", func(t *testing.T) {
		buf := make([]byte, 27)

		var header Header
		var panicked bool

		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
			if !panicked {
				t.Errorf("Expected copy(header.Compression[:], buf[24:28]) with buf of size 27 to panic, but it did not")
			}
		}()

		// This line should panic because buf[24:28] goes out of bounds on a slice of size 27.
		copy(header.Compression[:], buf[24:28])
	})
}

