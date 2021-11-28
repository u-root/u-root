// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// Tests using input from fuzzing runs. Ignores any errors, just checks that the
// inputs do not cause crashes.
//
// To update the input zip after a fuzzing run:
// cd fuzz/corpus
// zip ../../testdata/fuzz_in.zip *
//
// Similarly, the zip can be extracted to use as input corpus. See fuzz.go for
// go-fuzz instructions.
func TestFuzzInputs(t *testing.T) {
	if testing.CoverMode() != "" {
		// NOTE - since this test doesn't validate the outputs, coverage from
		// it is very low value, essentially inflating the coverage numbers.
		t.Skip("this test will inflate coverage")
	}
	// restore log behavior at end
	logOut := log.Writer()
	defer log.SetOutput(logOut)

	// no logging output for this test, to increase speed
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0o200)
	if err != nil {
		panic(err)
	}
	log.SetOutput(null)
	log.SetFlags(0)

	z, err := zip.OpenReader("testdata/fuzz_in.zip")
	if err != nil {
		t.Error(err)
	}
	defer z.Close()
	for i, zf := range z.File {
		name := fmt.Sprintf("%03d_%s", i, zf.Name[:10])
		t.Run(name, func(t *testing.T) {
			f, err := zf.Open()
			if err != nil {
				t.Error(err)
			}
			data, err := io.ReadAll(f)
			if err != nil {
				t.Error(err)
			}
			// ignore any errors - just catch crashes
			list, err := ParseFilePathList(data)
			if err != nil {
				return
			}
			_ = list.String()
			for _, p := range list {
				r, err := p.Resolver()
				if err == nil {
					_ = r.String()
				}
				if strings.Contains(strings.ToLower(p.String()), "acpi") {
					t.Logf("acpi: %s ", p.String())
				}
			}
		})
	}
}
