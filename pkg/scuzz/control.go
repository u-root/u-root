// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

// Request is an interface for scuzz requests.
// Different kernels have different ways of defining
// a request. The most common is ioctl, and will require
// a Cmd and a Packet, passed as uintptr.
// Some systems, e.g. Plan9, will require a string.
type Request interface {
	Cmd() uintptr
	Packet() uintptr
	String() string
}

// Disk is the interface to a disk, with operations
// to create packets and operate on them.
type Disk interface {
	// UnlockRequest generates an unlock Request
	UnlockRequest(string, uint, bool) Request
	// Operate performs the operation defined in Request on the Disk.
	Operate(Request) error
}

// Unlock unlocks a disk using a password and timeout.
func Unlock(d Disk, password string, timeout uint, master bool) error {
	p := d.UnlockRequest(password, timeout, master)
	return d.Operate(p)
}
