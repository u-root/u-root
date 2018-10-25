// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strace

// gvisor uses cgo and that won't do.
type iovec struct {
	P Addr   /* Starting address */
	S uint32 /* Number of bytes to transfer */
}

type Addr uintptr
type Arg uintptr
