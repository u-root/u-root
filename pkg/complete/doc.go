// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package complete implements a simple completion package
// designed to be used in shells and other programs. It currently
// offers completion functions to implement table-based and
// file search path based completion. It also offers a multi
// completion capability so that you can construct completions
// from other completions.
//
// Goals:
// small code base, so it can easily be embedded in firmware
//
// easily embedded in other programs, like the ip command
//
// friendly to mixed modes, i.e. if we say
// ip l
// and stdin is interactive, it would be nice if ip dropped into
// a command line prompt and let you use completion to get the rest
// of the line, instead of printing out a bnf
//
// The structs should be very light weight and hence cheap to build, use,
// and throw away. They should NOT have lots of state.
//
// Rely on the fact that system calls and kernels are fast and cache file system
// info so you should not. This means that we don't need to put huge effort into building in-memory
// structs representing file system information. Just ask the kernel.
//
// Non-Goals:
// be just like bash or zsh
//
// do extensive caching from the file system or environment. There was a time
// (the 1970s as it happens) when extensive in-shell hash tables made sense.
// Disco balls were also big. We don't need either.
//
// Use:
// see the code, but basically, you can create completer and call a function
// to read one word. The intent is that completers are so cheap that just
// creating them on demand costs nothing. So far this seems to work.
package complete
