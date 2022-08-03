// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logutil

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileWriter(t *testing.T) {
	for _, tt := range []struct {
		name          string
		dirPath       string
		filename      string
		maxSize       int
		startContent  []byte
		appendContent []byte
		wantContent   []byte
		wantError     bool
	}{
		{
			name:          "append to file",
			dirPath:       "",
			filename:      "file.log",
			maxSize:       1024,
			startContent:  []byte("foo"),
			appendContent: []byte("bar"),
			wantContent:   []byte("foobar"),
			wantError:     false,
		},
		{
			name:          "append to file too large",
			dirPath:       "",
			filename:      "file.log",
			maxSize:       -1,
			startContent:  []byte("foo"),
			appendContent: []byte("bar"),
			wantContent:   []byte("foo"),
			wantError:     false,
		},
		{
			name:          "append overflow",
			dirPath:       "",
			filename:      "file.log",
			maxSize:       5,
			startContent:  []byte("foo"),
			appendContent: []byte("bar"),
			wantContent:   []byte("fooba"),
			wantError:     false,
		},
		{
			name:          "dir missing",
			dirPath:       "dir",
			filename:      "file.log",
			maxSize:       1024,
			startContent:  []byte(""),
			appendContent: []byte("bar"),
			wantContent:   []byte(""),
			wantError:     true,
		},
	} {
		dir, err := os.MkdirTemp("", "testdir")
		if err != nil {
			t.Errorf("TestNewFileWriter(%s): MkdirTemp errored: %v", tt.name, err)
		}
		defer os.RemoveAll(dir)
		if tt.dirPath != "" {
			dir = filepath.Join(dir, tt.dirPath)
		}
		if len(tt.startContent) > 0 {
			f, err := os.Create(filepath.Join(dir, tt.filename))
			if err != nil {
				t.Errorf("TestNewFileWriter(%s): Creating start file errored: %v", tt.name, err)
			}
			n, err := f.Write(tt.startContent)
			if err != nil {
				t.Errorf("TestNewFileWriter(%s): Start file errored: %v", tt.name, err)
			}
			if n != len(tt.startContent) {
				t.Errorf("TestNewFileWriter(%s): Start file Write() got %v, expected %v", tt.name, n, len(tt.startContent))
			}
			f.Close()
		}
		w, err := NewWriterToFile(tt.maxSize, dir, tt.filename)
		if (err != nil) != tt.wantError {
			t.Errorf("TestNewFileWriter(%s): NewWriterToFile errored: %v, expected error: %v", tt.name, err, tt.wantError)
		}
		if tt.wantError {
			continue
		}
		n, err := w.Write(tt.appendContent)
		if err != nil {
			t.Errorf("TestNewFileWriter(%s): Write errored: %v", tt.name, err)
		}
		if n != len(tt.appendContent) {
			t.Errorf("TestNewFileWriter(%s): Write() got %v, want %v", tt.name, n, len(tt.appendContent))
		}

		dat, err := os.ReadFile(filepath.Join(dir, tt.filename))
		if err != nil {
			t.Errorf("TestNewFileWriter(%s): ReadFile errored with: %v", tt.name, err)
		}
		if !bytes.Equal(dat, tt.wantContent) {
			t.Errorf("TestNewFileWriter(%s): got %v, expected %v", tt.name, dat, tt.wantContent)
		}
	}
}
