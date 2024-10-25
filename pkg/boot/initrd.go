// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/uio/uio"
)

// CatInitrdsWithFileCache lazily reads up multiple initrds into single tmpfs file
// and return a os.File disguising as a io.ReaderAt.
// It starts processing after first ReadAt call is made.
func CatInitrdsWithFileCache(initrds ...io.Reader) io.ReaderAt {
	var names []string
	for _, initrd := range initrds {
		names = append(names, stringer(initrd))
	}
	return uio.NewLazyOpenerAt(strings.Join(names, ","), func() (io.ReaderAt, error) {
		f, err := os.CreateTemp("", "combined-initrd")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		for i, ireader := range initrds {
			size, err := io.Copy(f, ireader)
			if err != nil {
				return nil, err
			}
			// Don't pad the ending or an already aligned file.
			if i != len(initrds)-1 && size%512 != 0 {
				padding := make([]byte, 512-(size%512))
				nr, err := f.Write(padding)
				if err != nil {
					return nil, err
				}
				if nr != len(padding) {
					return nil, fmt.Errorf("write padding: %w", err)
				}
			}
		}
		if err := f.Sync(); err != nil {
			return nil, err
		}
		// Return a read-only copy.
		readOnlyF, err := os.Open(f.Name())
		if err != nil {
			return nil, err
		}
		return readOnlyF, nil
	})
}

// CatInitrds concatenates initrds on first ReadAt call from a list of
// io.ReaderAts, pads them to a 512 byte boundary.
func CatInitrds(initrds ...io.ReaderAt) io.ReaderAt {
	var names []string
	for _, initrd := range initrds {
		names = append(names, stringer(initrd))
	}

	return uio.NewLazyOpenerAt(strings.Join(names, ","), func() (io.ReaderAt, error) {
		buf := new(bytes.Buffer)
		for i, ireader := range initrds {
			size, err := buf.ReadFrom(uio.Reader(ireader))
			if err != nil {
				return nil, err
			}
			// Don't pad the ending or an already aligned file.
			if i != len(initrds)-1 && size%512 != 0 {
				padding := make([]byte, 512-(size%512))
				buf.Write(padding)
			}
		}
		// Buffer doesn't implement ReadAt, so wrap in NewReader
		return bytes.NewReader(buf.Bytes()), nil
	})
}

// CreateInitrd creates an initrd with the collection of files passed in.
func CreateInitrd(files ...string) (io.ReaderAt, error) {
	b := &bytes.Buffer{}
	archiver, err := cpio.Format("newc")
	if err != nil {
		return nil, err
	}
	w := archiver.Writer(b)
	cr := cpio.NewRecorder()
	// to deconflict names, we may want to prepend the names with
	// kexec_extra/ or something.
	for _, n := range files {
		rec, err := cr.GetRecord(n)
		if err != nil {
			return nil, fmt.Errorf("getting record of %q failed: %w", n, err)
		}
		if err := w.WriteRecord(rec); err != nil {
			return nil, fmt.Errorf("writing record %q failed: %w", n, err)
		}
	}
	if err := cpio.WriteTrailer(w); err != nil {
		return nil, fmt.Errorf("error writing trailer record: %w", err)
	}
	return bytes.NewReader(b.Bytes()), nil
}
