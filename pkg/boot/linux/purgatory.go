// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

// Purgatory abstracts a executable kexec purgatory in golang.
type Purgatory struct {
	Name    string
	Hexdump string
	Code    []byte
}
