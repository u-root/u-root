// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

// Booter is an interface that defines custom boot types. Implementations can be
// like network boot, local boot, etc. Boolean debugEnabled can be used to turning on/off
// the Booter debugging log.
type Booter interface {
	Boot(debugEnabled bool) error
	TypeName() string
}
