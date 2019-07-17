// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"io"
)

type multiCloser struct {
	io.Writer
	writers []io.Writer
}

func (mc *multiCloser) Close() error {
	var allErr error
	for _, w := range mc.writers {
		if c, ok := w.(io.Closer); ok {
			if err := c.Close(); err != nil {
				allErr = err
			}
		}
	}
	return allErr
}

func ClosingMultiWriter(w ...io.Writer) io.WriteCloser {
	return &multiCloser{
		Writer:  io.MultiWriter(w...),
		writers: w,
	}
}
