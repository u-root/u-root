// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This `ed` is intended to be a feature-complete mimick of [GNU Ed](https://www.gnu.org/software/ed//).  It is a close enough mimick that the [GNU Ed Man Page](https://www.gnu.org/software/ed/manual/ed_manual.html) should be a reliable source of documentation.  Divergence from the man page is generally considered a bug (unless it's an added feature).
//
// There are a few known differences:
//
// - `ed` uses `go`'s `regexp` package, and as such may have a somewhat different regular expression syntax.  Note, however, that backreferences follow the `ed` syntax of `\<ref>`, not the `go` syntax of `$<ref>`.
// - there has been little/no attempt to make particulars like error messages match `GNU Ed`.
// - rather than being an error, the 'g' option for 's' simply overrides any specified count.
// - does not support "traditional" mode
//
// The following has been implemented:
// - Full line address parsing (including RE and markings)
// - Implmented commands: !, #, =, E, H, P, Q, W, a, c, d, e, f, h, i, j, k, l, m, n, p, q, r, s, t, u, w, x, y, z
//
// The following has *not* yet been implemented, but will be eventually:
// - Unimplemented commands: g, G, v, V
// - does not (yet) support "loose" mode
// - does not (yet) support "restricted" mod
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// flags
var (
	fSuppress = flag.Bool("s", false, "suppress counts")
	fPrompt   = flag.String("p", "*", "specify a command prompt")
)

// current FileBuffer
var buffer *FileBuffer

// current ed state
var state struct {
	fileName string // current filename
	lastErr  error
	printErr bool
	prompt   bool
	winSize  int
	lastRep  string
	lastSub  string
}

// Parse input and run command
func run(cmd string) (e error) {
	ctx := &Context{
		cmd: cmd,
	}
	if ctx.addrs, ctx.cmdOffset, e = buffer.ResolveAddrs(cmd); e != nil {
		return
	}
	if len(cmd) <= ctx.cmdOffset {
		// no command, default to print
		ctx.cmd += "p"
	}
	if exe, ok := cmds[ctx.cmd[ctx.cmdOffset]]; ok {
		buffer.Start()
		if e = exe(ctx); e != nil {
			return
		}
		buffer.End()
	} else {
		return fmt.Errorf("invalid command: %v", cmd[ctx.cmdOffset])
	}
	return
}

// Entry point
func main() {
	var e error
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-s] [-p <prompt>] [file]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "p" {
			state.prompt = true
		}
	})
	args := flag.Args()
	if len(args) > 1 { // we only accept one additional argument
		flag.Usage()
		os.Exit(1)
	}
	buffer = NewFileBuffer(nil)
	if len(args) == 1 { // we were given a file name
		state.fileName = args[0]
		// try to read in the file
		if _, e = os.Stat(state.fileName); os.IsNotExist(e) && !*fSuppress {
			fmt.Fprintf(os.Stderr, "%s: No such file or directory", state.fileName)
			// this is not fatal, we just start with an empty buffer
		} else {
			if buffer, e = FileToBuffer(state.fileName); e != nil {
				fmt.Fprintln(os.Stderr, e)
				os.Exit(1)
			}
			if !*fSuppress {
				fmt.Println(buffer.Size())
			}
		}
	}
	state.winSize = 22 // we don't actually support getting the real window size
	inScan := bufio.NewScanner(os.Stdin)
	if state.prompt {
		fmt.Printf("%s", *fPrompt)
	}
	for inScan.Scan() {
		cmd := inScan.Text()
		e = run(cmd)
		if e != nil {
			state.lastErr = e
			if !*fSuppress && state.printErr {
				fmt.Println(e)
			} else {
				fmt.Println("?")
			}
		}
		if state.prompt {
			fmt.Printf("%s", *fPrompt)
		}
	}
	if inScan.Err() != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v", inScan.Err())
		os.Exit(1)
	}
}
