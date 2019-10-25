// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// FindFileSystem returns nil if a file system is available for use.
//
// It rereads /proc/filesystems each time as the supported file systems can change
// as modules are added and removed.
func FindFileSystem(fs string) error {
	b, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return err
	}
	for _, l := range strings.Split(string(b), "\n") {
		f := strings.Fields(l)
		if (len(f) > 1 && f[0] == "nodev" && f[1] == fs) || (len(f) > 0 && f[0] != "nodev" && f[0] == fs) {
			return nil
		}
	}
	return fmt.Errorf("%s not found", fs)
}
