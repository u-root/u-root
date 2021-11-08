// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

// Params communicates boot information to
// the purgatory, and possibly the kernel.
type Params struct {
	Entry  uint64
	Params uint64
	_      [5]uint64
}
