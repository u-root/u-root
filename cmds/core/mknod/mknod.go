// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Unmount a filesystem at the specified path.
//
// Synopsis:
//     mknod PATH TYPE [MAJOR MINOR]
//
// Description:
//     Creates a special file at PATH of the given TYPE. If TYPE is b, c or u,
//     the MAJOR and MINOR number must be specified. If the TYPE is p, they
//     must not be specified.
package main

import "log"

func main() {
	if err := mknod(); err != nil {
		log.Fatalf("mknod: %v", err)
	}
}
