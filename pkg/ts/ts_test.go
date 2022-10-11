// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ts

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func TestPrependTimestamp(t *testing.T) {
	format := func(time.Time) string { return "#" }
	tests := []struct {
		name, input, want string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "single blank line",
			input: "\n",
			want:  "#\n",
		},
		{
			name:  "blank lines",
			input: "\n\n\n\n",
			want:  "#\n#\n#\n#\n",
		},
		{
			name:  "text",
			input: "hello\nworld\n\n\n",
			want:  "#hello\n#world\n#\n#\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := &PrependTimestamp{
				R:      bytes.NewBufferString(tt.input),
				Format: format,
			}
			got := &bytes.Buffer{}
			if _, err := io.Copy(got, pt); err != nil {
				t.Errorf("io.Copy returned an error: %v", err)
			}
			if !bytes.Equal(got.Bytes(), []byte(tt.want)) {
				t.Errorf("PrependTimestamp = %q; want %q", got.String(), tt.want)
			}
		})
	}
}

// TestPrependTimestampBuffering ensures two important properties with regards
// to buffering which would be easy to miss with just a readline implementation:
//  1. Data is printed before the whole line is available.
//  2. Timestamp is generated at the beginning of the line, not at the end
//     of the previous line.
func TestPrependTimestampBuffering(t *testing.T) {
	// Mock out the format function.
	i := 0
	format := func(time.Time) string {
		return fmt.Sprint(i)
	}

	// Control exactly how many bytes are returned with each call to data.Read.
	data := &bytes.Buffer{}
	pt := &PrependTimestamp{
		R:      data,
		Format: format,
	}

	// These tests must run sequentially, so no tt.Run().
	for _, tt := range []struct {
		i        int
		text     string
		buffer   []byte
		wantText string
		wantErr  error
	}{
		{
			// Data is printed before the whole line is available.
			i:        1,
			text:     "Waiting...",
			buffer:   make([]byte, 100),
			wantErr:  nil,
			wantText: "1Waiting...",
		},
		{
			text:     "DONE\n",
			buffer:   make([]byte, 100),
			wantErr:  nil,
			wantText: "DONE\n",
		},
		{
			// Timestamp is generated at the beginning of the line.
			i:        2,
			text:     "Hello",
			buffer:   make([]byte, 2),
			wantErr:  nil,
			wantText: "2H",
		},
		{
			text:     "",
			buffer:   make([]byte, 2),
			wantErr:  nil,
			wantText: "el",
		},
		{
			text:     "",
			buffer:   make([]byte, 2),
			wantErr:  nil,
			wantText: "lo",
		},
		{
			text:     "",
			buffer:   make([]byte, 2),
			wantErr:  io.EOF,
			wantText: "",
		},
	} {
		data.Write([]byte(tt.text))
		i = tt.i
		n, err := pt.Read(tt.buffer)
		if err != tt.wantErr {
			t.Errorf("PrependTimestamp.Read(%q) err = %v; want %v",
				tt.text, err, tt.wantErr)
		}
		if !bytes.Equal(tt.buffer[:n], []byte(tt.wantText)) {
			t.Errorf("PrependTimestamp.Read(%q) buffer = %q; want %q",
				tt.text, tt.buffer[:n], tt.wantText)
		}
	}
}

// BenchmarkPrependTime measures the throughput of PrependTime where N is
// measured in bytes.
func BenchmarkPrependTime(b *testing.B) {
	line := "hello world\n"
	data := strings.Repeat(line, (b.N+len(line))/len(line))[:b.N]
	pt := New(bytes.NewBufferString(data))
	b.ResetTimer()
	if _, err := io.Copy(io.Discard, pt); err != nil {
		b.Fatal(err)
	}
}
