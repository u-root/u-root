// Copyright 2018-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// msr -- read and write MSRs with regular command or Forth
//
// Synopsis:
//
//	msr [OPTIONS] r glob MSR
//	msr [OPTIONS] w glob MSR value
//	msr [OPTIONS] forth-word [forth-word ...]
//
// Description:
//
//	This program reads and writes sets of MSRs, while allowing
//	them to by changed on a core by core or collective basis.
//
//	To read the msrs for 0 (sorry, the command is msr and the forth command
//	is msr, making this a bit confusing):
//	sudo msr 0 msr 0x3a reg rd
//	Breaking that down:
//	0 - for cpu 0
//	msr - for take the glob, in this case 0, and push all matching filenames on the stack
//	0x3a - for register 0x3a
//	reg - convert to 32-bit integer and push
//	rd - pop a 32-bit int and a []string and use them to read 1 or more MSRs
//
//	for all:
//	sudo msr "'*" msr 0x3a reg rd
//
//	The "'" is needed to quote the * so forth does not think we're multiplying.
//
//	Here is a breakdown, running msr with each command in turn:
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0
//	0
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr
//	[/dev/cpu/0/msr]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr 0x3a
//	[[/dev/cpu/0/msr] 0x3a]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr 0x3a reg
//	[[/dev/cpu/0/msr] 58]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr 0x3a reg rd
//	[0]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$
//
//	To read, then write all of them
//	(the dup is so we have the msr list at TOS -- it's just a convenience)
//	sudo msr 0 msr dup 0x3a reg rd 0x3a reg swap 1 u64 or wr
//	to just write them
//
//	Also, note, all the types are checked by assertions. The reg has to be
//	32 bits, the val 64
//
//	more examples:
//	read reg 0x3a and leave the least of MSRs and values on TOS, to be
//	printed at exit.
//
//	sudo msr "'"* msr dup 0x3a reg rd
//	[[/dev/cpu/0/msr /dev/cpu/1/msr /dev/cpu/2/msr /dev/cpu/3/msr] [5 5 5 5]]
//
//	From there, the write is easy:
//	sudo msr "'"* msr dup 0x3a reg rd 0x3a reg swap wr
//
//	For convenience, we maintain the old read and write commands:
//	msr r <glob> <register>
//	msr w <glob> <register> <value>
//
//	Yep, it's a bit inconvenient; the idea is that in the simple case,
//	you will use the r and w commands. For programmatic cases, you can
//	work to build up a working set of arguments.
//
//	For example, I started with
//	msr "'"* msr
//	and, once I saw the MSR selection was OK, built the command up from
//	there. At each step I could see the stack and whether I was going the
//	right direction.
//
// The old commands remain:
// rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ sudo ./msr r 0 0x3a
// [5]
// rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ sudo ./msr w 0 0x3a 5
// [5]
//
// For a view of what Forth is doing, run with -d.
package main
