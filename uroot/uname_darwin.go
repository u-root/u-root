// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package uroot contains various functions that might be needed more than
// one place.
package uroot

type Utsname struct {
	Sysname    string
	Nodename   string
	Release    string
	Version    string
	Machine    string
	Domainname string
}

// uname does a uname and returns a uroot.Utsname
func Uname() (*Utsname, error) {
	return &Utsname{Sysname: "unknown", Nodename: "unknown", Release: "nastylion", Version: "10.x", Machine: "amd64", Domainname: "unknown"}, nil
}
