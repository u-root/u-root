// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"io"

	"github.com/u-root/u-root/pkg/uio"
)

// CatInitrds concatenates initrds from a list of io.Readers, pads them to a
// 512 byte boundary, and returns a uio.LazyOpenerAt
func CatInitrds(initrds ...io.Reader) *uio.LazyOpenerAt {
	return uio.NewLazyOpenerAt("", func() (io.ReaderAt, error) {
		buf := new(bytes.Buffer)
		for _, ireader := range initrds {
			size, err := buf.ReadFrom(ireader)
			if err != nil {
				return nil, err
			}
			padding := make([]byte, 512-(size%512))
			buf.Write(padding)
		}
		// Buffer doesn't implement ReadAt, so wrap in NewReader
		return bytes.NewReader(buf.Bytes()), nil
	})
}
