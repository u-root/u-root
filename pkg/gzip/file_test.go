package gzip

import "testing"

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

func TestFile_checkOutPath(t *testing.T) {
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
			if err := f.checkOutPath(); (err != nil) != tt.wantErr {
				t.Errorf("File.checkOutPath() error = %v, wantErr %v", err, tt.wantErr)
			}
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
