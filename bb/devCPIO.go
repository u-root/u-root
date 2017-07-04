// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"syscall"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
)

const (
	d = syscall.S_IFDIR
	c = syscall.S_IFCHR
	b = syscall.S_IFBLK
	f = syscall.S_IFREG

	// This is the literal timezone file for GMT-0. Given that we have no idea
	// where we will be running, GMT seems a reasonable guess. If it matters,
	// setup code should download and change this to something else.
	gmt0       = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"
	nameserver = "nameserver 8.8.8.8\n"
)

// devCPIOrecords are cpio records as defined in the uroot cpio package.
// Most of the bits can be left unspecified: these all have one link,
// they are mostly root:root, for example.
var devCPIO = []cpio.Record{
	{Info: cpio.Info{Name: "tcz", Mode: d | 0755}},
	{Info: cpio.Info{Name: "etc", Mode: d | 0755}},
	{Info: cpio.Info{Name: "dev", Mode: d | 0755}},
	{Info: cpio.Info{Name: "usr", Mode: d | 0755}},
	{Info: cpio.Info{Name: "usr/lib", Mode: d | 0755}},
	{Info: cpio.Info{Name: "lib64", Mode: d | 0755}},
	{Info: cpio.Info{Name: "bin", Mode: d | 0755}},
	{Info: cpio.Info{Name: "dev/console", Mode: c | 0600, Rmajor: 5, Rminor: 1}},
	{Info: cpio.Info{Name: "dev/tty", Mode: c | 0666, Rmajor: 5, Rminor: 0}},
	{Info: cpio.Info{Name: "dev/null", Mode: c | 0666, Rmajor: 1, Rminor: 3}},
	{Info: cpio.Info{Name: "dev/urandom", Mode: c | 0666, Rmajor: 1, Rminor: 9}},
	{Info: cpio.Info{Name: "etc/resolv.conf", Mode: f | 0644, FileSize: uint64(len(nameserver))}, ReadCloser: cpio.NewBytesReadCloser([]byte(nameserver))},
	{Info: cpio.Info{Name: "etc/localtime", Mode: f | 0644, FileSize: uint64(len(gmt0))}, ReadCloser: cpio.NewBytesReadCloser([]byte(gmt0))},
}

// not yet implemented, let's wait and see if we still need them:
// brw-rw----   1 root     wheel         7,6 May 23  2015 dev/loop6
// brw-rw----   1 root     wheel         7,0 May 23  2015 dev/loop0
// brw-rw----   1 root     wheel         7,4 May 23  2015 dev/loop4
// brw-rw----   1 root     wheel         7,3 May 23  2015 dev/loop3
// brw-rw----   1 root     wheel         7,1 May 23  2015 dev/loop1
// brw-rw----   1 root     wheel         7,5 May 23  2015 dev/loop5
// crw-------   1 root     wheel      10,237 May 23  2015 dev/loop-control
// brw-rw----   1 root     wheel         7,7 May 23  2015 dev/loop7
// brw-rw----   1 root     wheel         7,2 May 23  2015 dev/loop2
