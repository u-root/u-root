// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"testing"
)

// TestEOLReader tests the EOLReader's ability to append the specified EOL sequence to each line of input.
func TestEOLReader(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		eol           []byte
		expected      string
		expectedCount int64
	}{
		{
			name:          "NO EOL sequence",
			input:         "Hello, world",
			eol:           []byte("\n"),
			expected:      "Hello, world",
			expectedCount: 12,
		},
		{
			name:          "Single line LF",
			input:         "Hello, world\n",
			eol:           []byte("\n"),
			expected:      "Hello, world\n",
			expectedCount: 13,
		},
		{
			name:          "Single line CRLF",
			input:         "Hello, world\n",
			eol:           []byte("\r\n"),
			expected:      "Hello, world\r\n",
			expectedCount: 14,
		},
		{
			name:          "Multiple lines CRLF",
			input:         "Hello\nWorld",
			eol:           []byte("\r\n"),
			expected:      "Hello\r\nWorld",
			expectedCount: 12,
		},
		{
			name:          "Empty input",
			input:         "",
			eol:           []byte("\n"),
			expected:      "",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			eolReader := NewEOLReader(reader, tc.eol)

			var result bytes.Buffer

			n, err := io.Copy(&result, eolReader)
			if err != nil {
				t.Fatalf("Failed to copy data from EOLReader: %v", err)
			}

			if n != tc.expectedCount {
				t.Errorf("Expected count does not match actual count. Expected: %d, Got: %d", tc.expectedCount, n)
			}

			if !reflect.DeepEqual(result.String(), tc.expected) {
				t.Errorf("Expected result does not match actual result. Expected: %q, Got: %q", tc.expected, result.String())
			}
		})
	}
}

// mockWriter is a simple io.Writer implementation that writes to a bytes.Buffer.
// It's used to test the ConcurrentWriter.
type mockWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buf.Write(p)
}

func TestConcurrentWriter(t *testing.T) {
	mock := &mockWriter{}
	cw := NewConcurrentWriter(mock)

	var wg sync.WaitGroup
	writeCount := 100
	data := []byte("data")

	wg.Add(writeCount)
	for range writeCount {
		go func() {
			defer wg.Done()
			if _, err := cw.Write(data); err != nil {
				t.Errorf("Failed to write data: %v", err)
			}
		}()
	}

	wg.Wait()

	if got, want := len(mock.buf.Bytes()), len(data)*writeCount; got != want {
		t.Errorf("Expected buffer length %d, got %d", want, got)
	}
}
