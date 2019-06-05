// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Archive archives files. The VTOC is at the front; we're not modeling tape drives or
// streams as in tar and cpio. This will greatly speed up listing the archive,
// modifying it, and so on. We think.
// Why a new tool?
package main

import (
	"io"
	"os"
)

func encodeOne(out io.Writer, f *file) error {
	if f.Size == 0 || !f.Mode.IsRegular() {
		return nil
	}

	in, err := os.Open(f.Name)
	if err != nil {
		return err
	}
	defer in.Close()
	amt, err := io.Copy(out, in)
	debug("%s: wrote %d bytes", f.Name, amt)
	return err
}

func encode(out io.Writer, dirs ...string) error {
	var vtoc []*file
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	vtoc, err := buildVTOC(dirs)
	if err != nil {
		return err
	}

	amt, err := writeVTOC(out, vtoc)
	debug("Wrote %d bytes of vtoc", amt)

	for _, v := range vtoc {
		if err = encodeOne(out, v); err != nil {
			break
		}
	}

	return err
}
