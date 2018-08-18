// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The mount package implements functions for mounting and unmounting
// file systems and defines the mount interface.
package mount

type Mounter interface {
	Mount() error
	Unmount(flags int) error
}
