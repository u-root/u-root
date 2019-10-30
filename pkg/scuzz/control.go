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

// Response is an interface for scuzz responses.
// It can return an error indicating the drive
// operation did not go well, a raw status block,
// or a string suitable for human consumption.
// sgATA commands are a lot like RPCs, so we
// provide both the kernel error and an error
// in the Response.
//
// tl;dr
// Note the error from the Operate is a transport
// level error and the error in the Response
// describes target errors. For example,
// an error return from Operate might indicate
// permission denied or a similar issue; an error
// in the Response might indicate the drive had
// some issues with the request. Distinguishing
// the two types is complicated by the fact that
// the linux kernel interface returns sets errno
// for all cases and explicitly does not distinguish
// them -- it can't, Linux error return is limited to
// integers in the range 1 to 4095 for implementation
// reasons; it is not possible for the Linux API to be
// informative much beyond "something went wrong."
//
// You can think of the transport and target errors
// this way: in an ssh session, a transport error
// would be a socket write or read failing. A target
// error would be error messages read from the stderr
// channel. The commands we send to the disk over the
// transport (the fd in this case) are executed at the
// drive and status is returned in a status block, much
// as the remote commands return errors over stderr.
// The Linux SG layer combines these errors (it seems)
// such that should the remote device get an error,
// errno is set to EINVAL. It is up to code to figure
// out if the error was from the kernel or the device.
//
// Possibly we should just do what
// Linux does but we want to get some usage first
// before we decide.
//
// TODO
// Because most of Linux SG device error return
// seems to be limited to EINVAL, deeper analysis
// is needed to figure out what really happened.
type Response interface {
	Error() error
	Status() []byte
	String() string
}

// Disk is the interface to a disk, with operations
// to create packets and operate on them.
type Disk interface {
	// UnlockRequest generates an unlock Request.
	UnlockRequest(string, uint, bool) Request
	// IdentifyRequest generates an identify Request.
	IdentifyRequest(timeout uint) Request
	// Operate performs the operation defined in Request on the Disk and returns a Response and an error.
	Operate(Request) (Response, error)
}

// Unlock unlocks a disk using a password and timeout.
func Unlock(d Disk, password string, timeout uint, master bool) (Response, error) {
	p := d.UnlockRequest(password, timeout, master)
	return d.Operate(p)
}

// Identify gets disk identity information.
func Identify(d Disk, timeout uint) (Response, error) {
	p := d.IdentifyRequest(timeout)
	return d.Operate(p)
}
