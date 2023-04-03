// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"bytes"
	"io"
	"testing"
)

func TestDecompress(t *testing.T) {
	type args struct {
		r         io.Reader
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
			name: "Basic Decompress",
			args: args{
				r:         bytes.NewReader([]byte("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\nI-.Q\x80\x13\x00\x00\x00\x00\xff\xff\x01\x00\x00\xff\xffG?\xfc\xcc\x0e\x00\x00\x00")),
				blocksize: 128,
				processes: 1,
			},
			wantW:   []byte("Test Test Test"),
			wantErr: false,
		},
		{
			name: "Zeros",
			args: args{
				r:         bytes.NewReader([]byte("\x1f\x8b\b\x00\x00\tn\x88\x02\xff2 \x1d\x00\x00\x00\x00\xff\xff\x01\x00\x00\xff\xffR6\xe3\xeb3\x00\x00\x00")),
				blocksize: 128,
				processes: 1,
			},
			wantW:   []byte("000000000000000000000000000000000000000000000000000"),
			wantErr: false,
		},
		{
			name: "Nil",
			args: args{
				r:         bytes.NewReader([]byte(nil)),
				blocksize: 128,
				processes: 1,
			},
			wantErr: true,
		},
		{
			name: "Corrupt input",
			args: args{
				r:         bytes.NewReader([]byte("\x1f\x8b\b\x00\x00\t\x88\x04\xff\x00\x00\x00\xff\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")),
				blocksize: 128,
				processes: 1,
			},
			wantErr: true,
		},
		{
			name: "Invalid header trailing garbage",
			args: args{
				r:         bytes.NewReader([]byte("\x1f\x8b\b\x00\x00\tn\x88\x04\xff\x00\x00\x00\xff\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")),
				blocksize: 128,
				processes: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Decompress(tt.args.r, w, tt.args.blocksize, tt.args.processes); (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotW := w.Bytes()
			if !bytes.Equal(gotW, tt.wantW) {
				t.Errorf("Decompress() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
