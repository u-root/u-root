// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

// Package vartest contains utility functions for testing uefivars and
// subpackages. It is unlikely to be useful outside of testing.
package vartest

import (
	"archive/zip"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/upath"
)

// Extracts testdata zip for use as efivars in tests. Used in uefivars and subpackages.
func SetupVarZip(path string) (efiVarDir string, cleanup func(), err error) {
	efiVarDir, err = os.MkdirTemp("", "gotest-uefivars")
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(efiVarDir)
		}
	}()
	z, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer z.Close()
	for _, zf := range z.File {
		var fname string
		if fname, err = upath.SafeFilepathJoin(efiVarDir, zf.Name); err != nil {
			return
		}
		if zf.FileInfo().IsDir() {
			err = os.MkdirAll(fname, zf.FileInfo().Mode())
			if err != nil {
				return
			}
		} else {
			var fo *os.File
			fo, err = os.Create(fname)
			if err != nil {
				return
			}
			var fi io.ReadCloser
			fi, err = zf.Open()
			if err != nil {
				return
			}
			_, err = io.Copy(fo, fi)
			if err != nil {
				return
			}
			fo.Close()
		}
	}
	cleanup = func() { os.RemoveAll(efiVarDir) }
	return
}
