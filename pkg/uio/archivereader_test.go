// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/pierrec/lz4/v4"
)

const choices = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func TestArchiveReaderRegular(t *testing.T) {
	dataStr := strings.Repeat("This is an important data!@#$%^^&&*&**(()())", 1000)

	ar, err := NewArchiveReader(bytes.NewReader([]byte(dataStr)))
	if err != nil {
		t.Fatalf("newArchiveReader(bytes.NewReader(%v)) returned error: %v", []byte(dataStr), err)
	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, ar); err != nil {
		t.Errorf("io.Copy(%v, %v) returned error: %v, want nil.", buf, ar, err)
	}
	if buf.String() != dataStr {
		t.Errorf("got %s, want %s", buf.String(), dataStr)
	}
}

func TestArchiveReaderPreReadShort(t *testing.T) {
	dataStr := "short data"
	ar, err := NewArchiveReader(bytes.NewReader([]byte(dataStr)))
	if err != nil {
		t.Errorf("newArchiveReader(bytes.NewReader([]byte(%s))) returned err: %v, want nil", dataStr, err)
	}
	got, err := io.ReadAll(ar)
	if err != nil {
		t.Errorf("got error reading archive reader: %v, want nil", err)
	}
	if string(got) != dataStr {
		t.Errorf("got %s, want %s", string(got), dataStr)
	}
	// Pre-read nothing.
	dataStr = ""
	ar, err = NewArchiveReader(bytes.NewReader([]byte(dataStr)))
	if err != ErrPreReadError {
		t.Errorf("newArchiveReader(bytes.NewReader([]byte(%s))) returned err: %v, want %v", dataStr, err, ErrPreReadError)
	}
	got, err = io.ReadAll(ar)
	if err != nil {
		t.Errorf("got error reading archive reader: %v, want nil", err)
	}
	if string(got) != dataStr {
		t.Errorf("got %s, want %s", string(got), dataStr)
	}
}

// randomString generates random string of fixed length in a fast and simple way.
func randomString(l int) string {
	rand.Seed(time.Now().UnixNano())
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = byte(choices[rand.Intn(len(choices))])
	}
	return string(r)
}

func checkArchiveReaderLZ4(t *testing.T, tt archiveReaderLZ4Case) {
	t.Helper()

	srcR := bytes.NewReader([]byte(tt.dataStr))

	srcBuf := new(bytes.Buffer)
	lz4w := tt.setup(srcBuf)

	n, err := io.Copy(lz4w, srcR)
	if err != nil {
		t.Fatalf("io.Copy(%v, %v) returned error: %v, want nil", lz4w, srcR, err)
	}
	if n != int64(len([]byte(tt.dataStr))) {
		t.Fatalf("got %d bytes compressed, want %d", n, len([]byte(tt.dataStr)))
	}
	if err = lz4w.Close(); err != nil {
		t.Fatalf("Failed to close lz4 writer: %v", err)
	}

	// Test ArchiveReader reading it.
	ar, err := NewArchiveReader(bytes.NewReader(srcBuf.Bytes()))
	if err != nil {
		t.Fatalf("newArchiveReader(bytes.NewReader(%v)) returned error: %v", srcBuf.Bytes(), err)
	}
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, ar); err != nil {
		t.Errorf("io.Copy(%v, %v) returned error: %v, want nil.", buf, ar, err)
	}
	if buf.String() != tt.dataStr {
		t.Errorf("got %s, want %s", buf.String(), tt.dataStr)
	}
}

type archiveReaderLZ4Case struct {
	name    string
	setup   func(w io.Writer) *lz4.Writer
	dataStr string
}

func TestArchiveReaderLZ4(t *testing.T) {
	for _, tt := range []archiveReaderLZ4Case{
		{
			name: "non-legacy regular",
			setup: func(w io.Writer) *lz4.Writer {
				return lz4.NewWriter(w)
			},
			dataStr: randomString(1024),
		},
		{
			name: "non-legacy larger data",
			setup: func(w io.Writer) *lz4.Writer {
				return lz4.NewWriter(w)
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "non-legacy short data", // Likley not realistic for most cases in the real world.
			setup: func(w io.Writer) *lz4.Writer {
				return lz4.NewWriter(w)
			},
			dataStr: randomString(100), // Smaller than pre-read size, 1024 bytes.
		},
		{
			name: "legacy regular",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				return lz4w
			},
			dataStr: randomString(1024),
		},
		{
			name: "legacy larger data",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "legacy small data",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				return lz4w
			},
			dataStr: randomString(100), // Smaller than pre-read size, 1024 bytes..
		},
		{
			name: "legacy small data",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				return lz4w
			},
			dataStr: randomString(100), // Smaller than pre-read size, 1024 bytes..
		},
		{
			name: "regular larger data with fast compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Fast))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "legacy larger data with fast compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Fast))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkArchiveReaderLZ4(t, tt)
		})
	}
}

func TestArchiveReaderLZ4SlowCompressed(t *testing.T) {
	for _, tt := range []archiveReaderLZ4Case{
		{
			name: "regular larger data with medium compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Level5))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "regular larger data with slow compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Level9))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "legacy larger data with medium compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Level5))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
		{
			name: "legacy larger data with slow compression",
			setup: func(w io.Writer) *lz4.Writer {
				lz4w := lz4.NewWriter(w)
				lz4w.Apply(lz4.LegacyOption(true))
				lz4w.Apply(lz4.CompressionLevelOption(lz4.Level9))
				return lz4w
			},
			dataStr: randomString(5 * 1024),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkArchiveReaderLZ4(t, tt)
		})
	}
}
