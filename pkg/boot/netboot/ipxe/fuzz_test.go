// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipxe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/ulog/ulogtest"
)

func FuzzParseIpxeConfig(f *testing.F) {
	seeds, err := filepath.Glob("testdata/fuzz/corpora/*.seed")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}
	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			return
		}

		parser := &parser{
			log: ulogtest.Logger{t},
		}
		parser.parseIpxe(string(data))
	})
}
