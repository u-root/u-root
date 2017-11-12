// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

// Uname returns a uroot.Utsname.
//
// Darwin has no uname system call and it's pretty fixed in many ways, so we go
// with this set of usable variables.
func Uname() (*Utsname, error) {
	return &Utsname{
		Sysname:    "unknown",
		Nodename:   "unknown",
		Release:    "nastylion",
		Version:    "10.x",
		Machine:    "amd64",
		Domainname: "unknown",
	}, nil
}
