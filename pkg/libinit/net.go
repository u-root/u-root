// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

var osNetInit = func() {}

// NetInit is u-root network initialization.
func NetInit() {
	osNetInit()
}
