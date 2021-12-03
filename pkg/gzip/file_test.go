// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_file_outputPath(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Stdout",
			fields: fields{Path: "/dev/stdout", Options: &Options{Stdout: true}},
			want:   "/dev/stdout",
		},
		{
			name:   "Test",
			fields: fields{Path: "/dev/null", Options: &Options{Test: true}},
			want:   "/dev/null",
		},
		{
			name:   "Compress",
			fields: fields{Path: "/tmp/test", Options: &Options{Suffix: ".gz"}},
			want:   "/tmp/test.gz",
		},
		{
			name:   "Decompress",
			fields: fields{Path: "/tmp/test.gz", Options: &Options{Decompress: true, Suffix: ".gz"}},
			want:   "/tmp/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
				Options: tt.fields.Options,
			}
			if got := f.outputPath(); got != tt.want {
				t.Errorf("file.outputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_CheckPath(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "skip decompressing already does not have suffix",
			fields: fields{
				Path: "file",
				Options: &Options{
					Decompress: true,
					Suffix:     ".gz",
				},
			},
			wantErr: true,
		},
		{
			name: "skip compressing already has suffix",
			fields: fields{
				Path: "file.gz",
				Options: &Options{
					Decompress: false,
					Suffix:     ".gz",
				},
			},
			wantErr: true,
		},
	}

	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := os.Create(filepath.Join(tempDir, tt.fields.Path))
			if err != nil {
				t.Fatalf("File.CheckPath() error can't create temp file: %v", err)
			}
			defer path.Close()

			f := &File{
				Path:    path.Name(),
				Options: tt.fields.Options,
			}
			if err := f.CheckPath(); (err != nil) != tt.wantErr {
				t.Errorf("File.CheckPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_CheckOutputPath(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "don't check output path if Stdout is true",
			fields: fields{
				Options: &Options{
					Stdout: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
				Options: tt.fields.Options,
			}
			if err := f.CheckOutputPath(); (err != nil) != tt.wantErr {
				t.Errorf("File.CheckOutputPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_CheckOutputStdout(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Stdout compress to device",
			fields: fields{
				Path:    "/dev/null",
				Options: &Options{Stdout: true, Decompress: false, Force: false},
			},
			wantErr: true,
		},
		{
			name: "Stdout compress to device force",
			fields: fields{
				Path:    "/dev/null",
				Options: &Options{Stdout: true, Decompress: false, Force: true},
			},
			wantErr: false,
		},
		{
			name: "Stdout compress redirect to file",
			fields: fields{
				Path:    "/tmp/test",
				Options: &Options{Stdout: true, Decompress: false, Force: false},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
				Options: tt.fields.Options,
			}
			oldStdout := os.Stdout
			var stdout *os.File
			if f.Path[0:4] == "/dev" {
				stdout, _ = os.Open(f.Path)
			} else {
				stdout, _ = os.Create(f.Path)
				defer os.Remove(f.Path)
			}
			defer stdout.Close()

			os.Stdout = stdout
			if err := f.CheckOutputStdout(); (err != nil) != tt.wantErr {
				t.Errorf("File.checkOutStdout() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.Stdout = oldStdout
		})
	}
}

func TestFile_Cleanup(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		exists  bool
		wantErr bool
	}{
		{
			name: "file should be deleted",
			fields: fields{
				Options: &Options{},
			},
			exists:  false,
			wantErr: false,
		},
		{
			name: "file should not be deleted Keep is true",
			fields: fields{
				Options: &Options{Keep: true},
			},
			exists:  true,
			wantErr: false,
		},
		{
			name: "file should not be deleted Stdout is true",
			fields: fields{
				Options: &Options{Stdout: true},
			},
			exists:  true,
			wantErr: false,
		},
		{
			name: "file should not be deleted Test is true",
			fields: fields{
				Options: &Options{Stdout: true},
			},
			exists:  true,
			wantErr: false,
		},
	}

	tempDir := t.TempDir()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := os.CreateTemp(tempDir, "cleanup-test")
			if err != nil {
				t.Errorf("File.Cleanup() error can't create temp file: %v", err)
			}
			defer path.Close()
			f := &File{
				Path:    path.Name(),
				Options: tt.fields.Options,
			}
			if err := f.Cleanup(); (err != nil) != tt.wantErr {
				t.Errorf("File.Cleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, err = os.Stat(f.Path)
			if tt.exists && err != nil {
				t.Errorf("File.Cleanup() file should be deleted")
			}
			if !tt.exists && err == nil {
				t.Errorf("File.Cleanup() file should stay")
			}
		})
	}
}

func TestFile_Process(t *testing.T) {
	tempDir := t.TempDir()
	path, err := os.CreateTemp(tempDir, "process-test")
	if err != nil {
		t.Fatalf("File.Process() error can't create temp file: %v", err)
	}
	defer path.Close()

	f := File{
		Path: path.Name(),
		Options: &Options{
			Decompress: false,
			Blocksize:  128,
			Level:      -1,
			Processes:  1,
			Suffix:     ".gz",
		},
	}

	if err := f.Process(); err != nil {
		t.Errorf("File.Process() compression error = %v", err)
	}

	f.Path = f.Path + f.Options.Suffix
	f.Options.Decompress = true
	if err := f.Process(); err != nil {
		t.Errorf("File.Process() decompression error = %v", err)
	}
}
