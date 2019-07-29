// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mount implements mounting, moving, and unmounting file systems.
package mount

// Mounter is an object that can be mounted.
type Mounter interface {
	Mount() error
	Unmount(flags int) error
}
