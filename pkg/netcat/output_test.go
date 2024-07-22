// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netcat

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteMultipleCases(t *testing.T) {
	// Setup temporary files for testing
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "output.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpHexFile, err := os.CreateTemp(tmpDir, "output.hex")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpHexFile.Name())

	testcases := []struct {
		name              string
		appendOutput      bool
		data              []byte
		expectFileContent string
		outFilePath       string
		outFileHexPath    string
		wantErr           bool
	}{
		{
			name:         "No output files",
			appendOutput: false,
			data:         []byte("Hello, World\n"),
		},
		{
			name:              "Write without append",
			appendOutput:      false,
			data:              []byte("Hello, World\n"),
			expectFileContent: "Hello, World\n",
			outFilePath:       tmpFile.Name(),
			outFileHexPath:    tmpHexFile.Name(),
		},
		{
			name:              "Write with append",
			appendOutput:      true,
			data:              []byte("appendedoutput\n"),
			expectFileContent: "in_case_of_append\nappendedoutput\n",
			outFilePath:       tmpFile.Name(),
			outFileHexPath:    tmpHexFile.Name(),
		},
		{
			name:        "Write fail",
			data:        []byte("data"),
			outFilePath: "path/to/nonexistent/file",
			wantErr:     true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(tmpFile.Name(), []byte("in_case_of_append\n"), 0o644); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}

			// OutputOptions setup
			opts := OutputOptions{
				AppendOutput:   tt.appendOutput,
				OutFilePath:    tt.outFilePath,
				OutFileHexPath: tt.outFileHexPath,
			}

			n, err := opts.Write(tt.data)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("Write failed: %v", err)
			}

			if n > 0 {
				// Verify file content
				content, err := os.ReadFile(tt.outFilePath)
				if err != nil {
					t.Fatalf("Failed to read file: %v", err)
				}

				if diff := cmp.Diff(tt.expectFileContent, string(content)); diff != "" {
					t.Errorf("File content mismatch (-want +got):\n%s", diff)
				}

				// Verify hex file content
				if tt.outFileHexPath == "" {
					hexContent, err := os.ReadFile(tt.outFileHexPath)
					if err != nil {
						t.Fatalf("Failed to read hex file: %v", err)
					}

					// Since hexContent is already a string in hex format, compare directly
					if diff := cmp.Diff(hex.EncodeToString([]byte(tt.expectFileContent)), string(hexContent)); diff != "" {
						t.Errorf("Hex file content mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
