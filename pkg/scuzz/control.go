// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

// Disk is the interface to a disk, with operations
// to create packets and operate on them.
type Disk interface {
	// Unlock unlocks the drive, given a password
	Unlock(string, uint, bool) error
	// Identify returns drive identity information
	Identify(timeout uint) error
}
