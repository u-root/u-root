// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

func init() {
	SupportedLoaders = append(SupportedLoaders,
		// Note: At the time of writing, the kernel patch implementing
		// kexec_file_load(2) has not yet been merged (and is experimental).
		// Use at your own risk!
		FileLoad,
		// TODO: There no way to pass a modified device tree to
		// kexec_file_load. The only alternative is to use kexec_load (and do
		// all the hard work of constructing the segments ourselves).
	)
}
