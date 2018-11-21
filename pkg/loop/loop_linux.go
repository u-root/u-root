// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"golang.org/x/sys/unix"
)

const forceUnmount = unix.MNT_FORCE | unix.MNT_DETACH

// Unmount unmounts and frees a loop. If it is mounted, it will try to unmount it.
// If the unmount fails, we try to free it anyway, after trying a more
// forceful unmount. We don't log errors, but we do return a concatentation
// of whatever errors occur.
func (l *Loop) Unmount(flags int) error {
	if l.Mounted {
		if err := unix.Unmount(l.Dir, flags); err != nil {
			unix.Unmount(l.Dir, flags|forceUnmount)
		}
	}
	return ClearFile(l.Dev)
}
