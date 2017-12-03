package gzip

import (
	"bytes"
	"io"
	"testing"
)

func Test_compress(t *testing.T) {
	type args struct {
		r         io.Reader
		level     int
		blocksize int
		processes int
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.Buffer{}
			if err := compress(tt.args.r, &w, tt.args.level, tt.args.blocksize, tt.args.processes); (err != nil) != tt.wantErr {
				t.Errorf("compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("compress() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
