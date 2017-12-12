package gzip

import (
	"os"
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
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
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
	// TODO: Add test cases.
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
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
				Options: tt.fields.Options,
			}
			if err := f.Cleanup(); (err != nil) != tt.wantErr {
				t.Errorf("File.Cleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_Process(t *testing.T) {
	type fields struct {
		Path    string
		Options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:    tt.fields.Path,
				Options: tt.fields.Options,
			}
			if err := f.Process(); (err != nil) != tt.wantErr {
				t.Errorf("File.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
