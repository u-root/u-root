// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name      string
		plain     string
		level     int
		blocksize int
		processes int
	}{
		{
			name:      "Basic Compress",
			plain:     "Test Test Test",
			level:     9,
			blocksize: 128,
			processes: 1,
		},
		{
			name:      "Zeplainos",
			plain:     "000000000000000000000000000000000000000000000000000",
			level:     9,
			blocksize: 128,
			processes: 1,
		},
		{
			name:      "Empty stplaining",
			plain:     "",
			level:     1,
			blocksize: 128,
			processes: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ciphertext bytes.Buffer
			if err := Compress(strings.NewReader(tt.plain), &ciphertext, tt.level, tt.blocksize, tt.processes); err != nil {
				t.Fatalf("Compress() error = %v, want nil", err)
			}

			var plaintext strings.Builder
			if err := Decompress(&ciphertext, &plaintext, tt.blocksize, tt.processes); err != nil {
				t.Fatalf("Decompress() = %v, want nil", err)
			}

			got := plaintext.String()
			if !cmp.Equal(got, tt.plain) {
				t.Errorf("Compress() = %q, want %q -- diff\n%s", got, tt.plain, cmp.Diff(got, tt.plain))
			}
		})
	}
}
