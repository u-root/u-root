// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package termios provides two primary types: TTYIO and Termios.
// The purpose of this package is to allow programs to manipulate
// ttys, recording the original configuration if needed, and restoring it
// as needed.
//
// TTYIO is specific to a kernel, and provides basic operations for, e.g, setting
// raw mode. Because of the huge variance between kernels, the type is defined
// for each kernel type. For example, In Plan 9, the fd for data and
// fd for control are different; and ttys only exist in the window system.
//
// A Termios records information about the state of a terminal.
//
// See cmds/stty for some usage examples, but a typical sequence
// for stdin looks like this:
//
//		t, err := termios.GTTY(0)
//		if err != nil {
//			log.Fatalf("termios.GTTY: %v", err)
//		}
//	     t can then be manipulated as needed.
//
//	     Raw is provided, and returns a restorer:
//		restorer, err := termios.Raw(0)
//		if err != nil {
//				log.Fatalf("raw: %v", err)
//		}
//
//		err := termios.SetTermios(0, restorer)
//
// In keeping with all the variations in kernels and usages,
// there are several ways to accomplish this operation: see
// termios.go and the stty command for more.
//
// This package was written a long time ago, and does not match some
// aspects of today's idiomatic Go. But it has wide enough use that we
// do not plan to change it.
package termios
