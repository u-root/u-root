// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

func init() {
	SupportedLoaders = append(SupportedLoaders,
		ZImageLoad,
		// TODO: uimage
		// TODO: FIT
	)
}
