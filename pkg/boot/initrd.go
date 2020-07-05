// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"io"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
)

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
