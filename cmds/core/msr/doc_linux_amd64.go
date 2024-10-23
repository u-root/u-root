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
//	sudo msr 0 0x3a rd
//	Breaking that down:
//	0 - for cpu 0
//	0x3a - for register 0x3a
//	rd - pop a 32-bit int and a []string and use them to read 1 or more MSRs
//
//	for all:
//	sudo msr "'*" 0x3a rd
//
//	The "'" is needed to quote the * so forth does not think we're multiplying.
//
//	There are two convenience words so you can see how things break down.
//	You do not need them but they can help debug.
//	The above can also be written as
//	sudo msr "'*" cpu 0x3a reg rd
//
//	Here is a breakdown, running msr with each command in turn:
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0
//	0
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 cpu
//	[/dev/cpu/0/msr]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr 0x3a
//	[[/dev/cpu/0/msr] 0x3a]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 msr 0x3a reg
//	[[/dev/cpu/0/msr] 58]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$ ./msr 0 cpu 0x3a reg rd
//	[0]
//	rminnich@xcpu:~/gopath/src/github.com/u-root/u-root/cmds/core/msr$
//
//	the typeof word adds more information.
//	$ sudo ./msr 0 cpu
//	0
//	$ sudo ./msr 0 cpu typeof
//	msr.CPUs
//	$ sudo ./msr 0 cpu 0x3a
//	[0 0x3a]
//	$ sudo ./msr 0 cpu 0x3a reg
//	[0 0x3a]
//	$ sudo ./msr 0 cpu 0x3a reg typeof
//	[0 msr.MSR]
//	$ sudo ./msr 0 cpu 0x3a reg rd
//	[0 0x3a [5]]
//	$
//
//	To read, then write all of them
//	sudo msr 0 0x3a rd 1 or wr
//	to just write them
//
//	Also, note, all the types are checked by assertions. The reg has to be
//	32 bits, the val 64
//
//	more examples:
//	read reg 0x3a and leave the list of MSRs and values on TOS, to be
//	printed at exit.
//
//	sudo msr "'"* 0x3a rd
//	[[/dev/cpu/0/msr /dev/cpu/1/msr /dev/cpu/2/msr /dev/cpu/3/msr] [5 5 5 5]]
//
//	From there, the write is easy:
//	sudo msr "'"* 0x3a rd wr
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
//	msr "'"* cpu
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
