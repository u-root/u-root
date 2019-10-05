// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"os"
	"path/filepath"
)

var root = os.Getenv("UROOT_ROOT")

// UrootPath returns the absolute path for a uroot file with the UROOT_ROOT
// environment variable taken into account.
//
// It returns a proper value if UROOT_ROOT is not set.  u-root was built to
// assume everything is rooted at /, and in most cases that is still true.  But
// in hosted mode, e.g. on developer mode chromebooks, it's far better if
// u-root can be rooted in /usr/local, so successive kernel/root file system
// upgrades do not wipe it out.
func UrootPath(n ...string) string {
	return filepath.Join("/", root, filepath.Join(n...))
}
