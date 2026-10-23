// Copyright 2015-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// console implements a basic console. It establishes a pair of files
// to read from, the default being a UART at 0x3f8, but an alternative
// being just stdin and stdout. It will also set up a root file system
// using util.Rootfs, although this can be disabled as well.
// Console uses a Go version of fork_pty to start up a shell, default
// sh. Console runs until the shell exits and then exits itself.
//
// console takes a -serial switch to pick the console to use.
// -serial can be a file name, or an address. If an address,
// if it is a 16-bit number, it is assumed to be accessed
// via read(2)/write(2) to /dev/port. If more than 16 bits,
// it is assumed to use mmap access via /dev/mem.
//
// Raw hardware access is assumed to be using an 8250/16550a
// part. console only uses 2 registers: status and data.
//
// Example: use the uart at 0x102818000
// console -serial 0x102818000
// use the uart, but a different shell:
// console -serial 0x102818000 /usr/bin/bash
// use the uart, but create a root file system,
// useful if we are running as /init
// console -setuproot -serial 0x102818000 /usr/bin/bash
package main
