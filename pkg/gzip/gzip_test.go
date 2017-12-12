package gzip

import (
	"bytes"
	"io"
	"testing"
)

func Test_Compress(t *testing.T) {
	type args struct {
		r         io.Reader
		level     int
		blocksize int
		processes int
	}
	tests := []struct {
		name    string
		args    args
		wantW   []byte
		wantErr bool
	}{
		{
			name: "Basic Compress",
			args: args{
				r:         bytes.NewReader([]byte("Test Test Test")),
				level:     9,
				blocksize: 128,
				processes: 1,
			},
			wantW:   []byte("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\nI-.Q\x80\x13\x00\x00\x00\x00\xff\xff\x01\x00\x00\xff\xffG?\xfc\xcc\x0e\x00\x00\x00"),
			wantErr: false,
		},
		{
			name: "Zeros",
			args: args{
				r:         bytes.NewReader([]byte("000000000000000000000000000000000000000000000000000")),
				level:     9,
				blocksize: 128,
				processes: 1,
			},
			wantW:   []byte("\x1f\x8b\b\x00\x00\tn\x88\x02\xff2 \x1d\x00\x00\x00\x00\xff\xff\x01\x00\x00\xff\xffR6\xe3\xeb3\x00\x00\x00"),
			wantErr: false,
		},
		{
			name: "Nil",
			args: args{
				r:         bytes.NewReader([]byte(nil)),
				level:     1,
				blocksize: 128,
				processes: 1,
			},
			wantW:   []byte("\x1f\x8b\b\x00\x00\tn\x88\x04\xff\x00\x00\x00\xff\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.Buffer{}
			if err := Compress(tt.args.r, &w, tt.args.level, tt.args.blocksize, tt.args.processes); (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotW := w.Bytes()
			if !bytes.Equal(gotW, tt.wantW) {
				t.Errorf("Compress() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
