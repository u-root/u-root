package gzip

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func Test_appError_Error(t *testing.T) {
	type fields struct {
		msg   string
		level errType
		path  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Path",
			fields: fields{msg: "Test", path: "/tmp"},
			want:   "/tmp Test",
		},
		{
			name:   "No path",
			fields: fields{msg: "Test"},
			want:   "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				msg:   tt.fields.msg,
				level: tt.fields.level,
				path:  tt.fields.path,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("appError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	type args struct {
		err    error
		stdout io.Writer
		stderr io.Writer
	}
	type output struct {
		stdout []byte
		stderr []byte
		exit   int
	}
	tests := []struct {
		name string
		args args
		want output
	}{
		{
			name: "error type",
			args: args{err: errors.New("Test error")},
			want: output{exit: 1, stderr: []byte("error, Test error\n")},
		},
		{
			name: "appError no level",
			args: args{err: &appError{msg: "Test error"}},
			want: output{exit: 0, stdout: []byte("Test error\n")},
		},
		{
			name: "appError info level",
			args: args{err: &appError{msg: "Test error", level: info}},
			want: output{exit: 0, stdout: []byte("Test error\n")},
		},
		{
			name: "appError skipping level",
			args: args{err: &appError{msg: "Test error", level: skipping}},
			want: output{exit: 0, stderr: []byte("skipping, Test error\n")},
		},
		{
			name: "appError warning level",
			args: args{err: &appError{msg: "Test error", level: warning}},
			want: output{exit: 0, stderr: []byte("warning, Test error\n")},
		},
		{
			name: "appError fatal level",
			args: args{err: &appError{msg: "Test error", level: fatal}},
			want: output{exit: 1, stderr: []byte("fatal, Test error\n")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			if got := ErrorHandler(tt.args.err, &stdout, &stderr); got != tt.want.exit {
				t.Errorf("ErrorHandler() = exit: %v, want.exit: %v", got, tt.want.exit)
			}
			if !bytes.Equal(stdout.Bytes(), tt.want.stdout) {
				t.Errorf("ErrorHandler() = stdout: %q, want.stdout: %q", stdout.String(), tt.want.stdout)
			}
			if !bytes.Equal(stderr.Bytes(), tt.want.stderr) {
				t.Errorf("ErrorHandler() = stderr: %q, want.stderr: %q", stderr.String(), tt.want.stderr)
			}
		})
	}
}
