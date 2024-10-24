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
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/uroot/util"
)

// flags
var (
	fsuppress bool
	fprompt   string
	usage     = "Usage: ed [-s] [-p <prompt>] [file]\n"
)

func init() {
	flag.BoolVar(&fsuppress, "s", false, "suppress counts")
	flag.StringVar(&fprompt, "p", "*", "specify a command prompt")
	flag.Usage = util.Usage(flag.Usage, usage)
}

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

// Parse input and execute command
func execute(cmd string, output io.Writer) (e error) {
	ctx := &Context{
		cmd: cmd,
		out: output,
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

func runEd(in io.Reader, out io.Writer, suppress bool, prompt, file string) error {
	var e error
	if len(prompt) > 0 {
		state.prompt = true
	}
	buffer = NewFileBuffer(nil)
	if file != "" { // we were given a file name
		state.fileName = file
		// try to read in the file
		if _, e = os.Stat(state.fileName); os.IsNotExist(e) && !suppress {
			fmt.Fprintf(os.Stderr, "%s: No such file or directory", state.fileName)
			// this is not fatal, we just start with an empty buffer
		} else {
			if buffer, e = FileToBuffer(state.fileName); e != nil {
				return e
			}
			if !suppress {
				fmt.Println(buffer.Size())
			}
		}
	}
	state.winSize = 22 // we don't actually support getting the real window size
	inScan := bufio.NewScanner(in)
	if state.prompt {
		fmt.Fprintf(out, "%s", prompt)
	}
	for inScan.Scan() {
		cmd := inScan.Text()
		e = execute(cmd, out)
		if e != nil {
			state.lastErr = e
			if !suppress && state.printErr {
				fmt.Fprintf(out, "%s\n", e)
			} else {
				fmt.Fprintf(out, "?\n")
			}
			if errors.Is(e, errExit) {
				return nil
			}
		}
		if state.prompt {
			fmt.Printf("%s", prompt)
		}
	}
	if inScan.Err() != nil {
		return fmt.Errorf("error reading stdin: %w", inScan.Err())
	}
	return nil
}

// Entry point
func main() {
	flag.Parse()
	file := ""
	switch len(flag.Args()) {
	case 0:
	case 1:
		file = flag.Args()[0]
	default:
		flag.Usage()
		os.Exit(1)
	}
	if err := runEd(os.Stdin, os.Stdout, fsuppress, fprompt, file); err != nil {
		log.Fatal(err)
	}
}
