// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"io"
)

type PortReader interface {
	In(uint16, UintN) error
}

type PortWriter interface {
	Out(uint16, UintN) error
}

type PortReadWriter interface {
	PortReader
	PortWriter
	io.Closer
}
