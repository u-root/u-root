// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// commands.go - defines editor commands
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var errExit = fmt.Errorf("exit")

// A Context is passed to an invoked command
type Context struct {
	cmd       string // full command string
	cmdOffset int    // start of the command after address resolution
	addrs     []int  // resolved addresses
	out       io.Writer
}

// A Command can be run with a Context and returns an error
type Command func(*Context) error

// The cmds map maps single byte commands to their handler functions.
// This is also a good way to check what commands are implemented.
var cmds = map[byte]Command{
	'q': cmdQuit,
	'Q': cmdQuit,
	'd': cmdDelete,
	'l': cmdPrint,
	'p': cmdPrint,
	'n': cmdPrint,
	'h': cmdErr,
	'H': cmdErr,
	'a': cmdInput,
	'i': cmdInput,
	'c': cmdInput,
	'w': cmdWrite,
	'W': cmdWrite,
	'k': cmdMark,
	'e': cmdEdit,
	'E': cmdEdit,
	'r': cmdEdit,
	'f': cmdFile,
	'=': cmdLine,
	'j': cmdJoin,
	'm': cmdMove,
	't': cmdMove,
	'y': cmdCopy,
	'x': cmdPaste,
	'P': cmdPrompt,
	's': cmdSub,
	'u': cmdUndo,
	'D': cmdDump, // var dump the buffer for debug
	'z': cmdScroll,
	'!': cmdCommand,
	'#': func(*Context) (e error) { return },
}

//////////////////////
// Command handlers /
////////////////////

func cmdDelete(ctx *Context) (e error) {
	var r [2]int
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}
	e = buffer.Delete(r)
	return
}

func cmdQuit(ctx *Context) (e error) {
	if ctx.cmd[ctx.cmdOffset] == 'q' && buffer.Dirty() {
		return fmt.Errorf("warning: file modified")
	}
	return errExit
}

func cmdPrint(ctx *Context) (e error) {
	var r [2]int
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}
	for l := r[0]; l <= r[1]; l++ {
		if ctx.cmd[ctx.cmdOffset] == 'n' {
			fmt.Fprintf(ctx.out, "%d\t", l+1)
		}
		line := buffer.GetMust(l, true)
		if ctx.cmd[ctx.cmdOffset] == 'l' {
			line += "$" // TODO: the man pages describes more escaping, but it's not clear what GNU ed actually does.
		}
		fmt.Fprintf(ctx.out, "%s\n", line)
	}
	return
}

func cmdScroll(ctx *Context) (e error) {
	start, e := buffer.AddrValue(ctx.addrs)
	if e != nil {
		return
	}
	// parse win size (if there)
	winStr := ctx.cmd[ctx.cmdOffset+1:]
	if len(winStr) > 0 {
		var win int
		if win, e = strconv.Atoi(winStr); e != nil {
			return fmt.Errorf("invalid window size: %s", winStr)
		}
		state.winSize = win
	}
	end := min(start+state.winSize-1, buffer.Len()-1)
	var ls []string
	if ls, e = buffer.Get([2]int{start, end}); e != nil {
		return
	}
	for _, l := range ls {
		fmt.Fprintf(ctx.out, "%s\n", l)
	}
	return
}

func cmdErr(ctx *Context) (e error) {
	if ctx.cmd[ctx.cmdOffset] == 'h' {
		if state.lastErr != nil {
			fmt.Fprintf(ctx.out, "%s\n", state.lastErr)
			return
		}
	}
	if ctx.cmd[ctx.cmdOffset] == 'H' {
		if state.printErr {
			state.printErr = false
			return
		}
		state.printErr = true
	}
	return
}

func cmdInput(ctx *Context) (e error) {
	scan := bufio.NewScanner(os.Stdin)
	nbuf := []string{}
	if len(ctx.cmd[ctx.cmdOffset+1:]) != 0 && ctx.cmd[ctx.cmdOffset] != 'c' {
		return fmt.Errorf("%c only takes a single line addres", ctx.cmd[ctx.cmdOffset])
	}
	for scan.Scan() {
		line := scan.Text()
		if line == "." {
			break
		}
		nbuf = append(nbuf, line)
	}
	if len(nbuf) == 0 {
		return
	}
	switch ctx.cmd[ctx.cmdOffset] {
	case 'i':
		var line int
		if line, e = buffer.AddrValue(ctx.addrs); e != nil {
			return
		}
		e = buffer.Insert(line, nbuf)
	case 'a':
		var line int
		if line, e = buffer.AddrValue(ctx.addrs); e != nil {
			return
		}
		e = buffer.Insert(line+1, nbuf)
	case 'c':
		var r [2]int
		if r, e = buffer.AddrRange(ctx.addrs); e != nil {
			return
		}
		if e = buffer.Delete(r); e != nil {
			return
		}
		e = buffer.Insert(r[0], nbuf)
	}
	return
}

var rxWrite = regexp.MustCompile(`^(q)?(?: )?(!)?(.*)`)

func cmdWrite(ctx *Context) (e error) {
	file := state.fileName
	quit := false
	run := false
	var r [2]int
	if ctx.cmdOffset == 0 {
		r[0] = 0
		r[1] = buffer.Len() - 1
	} else {
		if r, e = buffer.AddrRange(ctx.addrs); e != nil {
			return
		}
	}
	m := rxWrite.FindAllStringSubmatch(ctx.cmd[ctx.cmdOffset+1:], -1)
	if m[0][1] == "q" {
		quit = true
	}
	if m[0][2] == "!" {
		run = true
	}
	if len(m[0][3]) > 0 {
		file = m[0][3]
	}
	var lstr []string
	lstr, e = buffer.Get(r)
	if e != nil {
		return
	}
	if run {
		s := System{
			Cmd:    m[0][3],
			Stdin:  bytes.NewBuffer(nil),
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		go func() {
			for _, str := range lstr {
				if _, e = fmt.Fprintf(s.Stdin.(*bytes.Buffer), "%s\n", str); e != nil {
					return
				}
			}
		}()
		return s.Run()
	}

	var f *os.File
	oFlag := os.O_TRUNC
	if ctx.cmd[ctx.cmdOffset] == 'W' {
		oFlag = os.O_APPEND
	}
	if f, e = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|oFlag, 0o666); e != nil {
		return e
	}
	defer f.Close()

	for _, s := range lstr {
		_, e = fmt.Fprintf(f, "%s\n", s)
		if e != nil {
			return
		}
	}
	if quit {
		if e = cmdQuit(ctx); e != nil {
			return
		}
	}
	buffer.Clean()
	return
}

func cmdMark(ctx *Context) (e error) {
	if len(ctx.cmd)-1 <= ctx.cmdOffset {
		e = fmt.Errorf("no mark character supplied")
		return
	}
	c := ctx.cmd[ctx.cmdOffset+1]
	var l int
	if l, e = buffer.AddrValue(ctx.addrs); e != nil {
		return
	}
	e = buffer.SetMark(c, l)
	return
}

func cmdEdit(ctx *Context) (e error) {
	var addr int
	// we do this manually because we allow addr 0
	if len(ctx.addrs) == 0 {
		return ErrINV
	}
	addr = ctx.addrs[len(ctx.addrs)-1]
	if addr != 0 && buffer.OOB(addr) {
		return ErrOOB
	}
	// cmd or filename?
	cmd := ctx.cmd[ctx.cmdOffset]
	force := false
	if cmd == 'E' || cmd == 'r' {
		force = true
	} // else == 'e'
	if buffer.Dirty() && !force {
		return fmt.Errorf("warning: file modified")
	}
	filename := ctx.cmd[ctx.cmdOffset+1:]
	filename = filename[wsOffset(filename):]
	var fh io.Reader
	if len(filename) == 0 {
		filename = state.fileName
	}
	if filename[0] == '!' { // command, not filename
		s := System{
			Cmd:    filename[1:],
			Stdout: bytes.NewBuffer(nil),
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}
		if e = s.Run(); e != nil {
			return
		}
		fh = s.Stdout.(io.Reader)
	} else { // filename
		if _, e = os.Stat(filename); os.IsNotExist(e) && !fsuppress {
			return fmt.Errorf("%s: No such file or directory", filename)
			// this is not fatal, we just start with an empty buffer
		}
		if fh, e = os.Open(filename); e != nil {
			e = fmt.Errorf("could not read file: %w", e)
			return
		}
		state.fileName = filename
	}

	if cmd != 'r' { // other commands replace
		buffer = NewFileBuffer(nil)
		if e = buffer.Read(0, fh); e != nil {
			return
		}
	} else {
		e = buffer.Read(addr, fh)
	}
	if !fsuppress {
		fmt.Fprintf(ctx.out, "%d\n", buffer.Size())
	}
	return
}

func cmdFile(ctx *Context) (e error) {
	newFile := ctx.cmd[ctx.cmdOffset:]
	newFile = newFile[wsOffset(newFile):]
	if len(newFile) > 0 {
		state.fileName = newFile
		return
	}
	fmt.Fprintf(ctx.out, "%s\n", state.fileName)
	return
}

func cmdLine(ctx *Context) (e error) {
	addr, e := buffer.AddrValue(ctx.addrs)
	if e == nil {
		fmt.Fprintf(ctx.out, "%d\n", addr+1)
	}
	return
}

func cmdJoin(ctx *Context) (e error) {
	var r [2]int
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}
	// Technically only a range works, but a line isn't an error
	if r[0] == r[1] {
		return
	}

	joined := ""
	for l := r[0]; l <= r[1]; l++ {
		joined += buffer.GetMust(l, false)
	}
	if e = buffer.Delete(r); e != nil {
		return
	}
	e = buffer.Insert(r[0], []string{joined})
	return
}

func cmdMove(ctx *Context) (e error) {
	var r [2]int
	var dest int
	var lines []string
	cmd := ctx.cmd[ctx.cmdOffset]
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}
	// must parse the destination
	destStr := ctx.cmd[ctx.cmdOffset+1:]
	var nctx Context
	if nctx.addrs, nctx.cmdOffset, e = buffer.ResolveAddrs(destStr); e != nil {
		return
	}
	// this is a bit hacky, but we're supposed to allow 0
	offset := 1
	last := len(nctx.addrs) - 1
	if nctx.addrs[last] == -1 {
		nctx.addrs[last] = 0
		offset = 0
	}
	if dest, e = buffer.AddrValue(nctx.addrs); e != nil {
		return
	}

	if lines, e = buffer.Get(r); e != nil {
		return
	}
	delt := r[1] - r[0] + 1
	if dest < r[0] {
		r[0] += delt
		r[1] += delt
	} else if dest > r[1] {
		// NOP
	} else {
		return fmt.Errorf("cannot move lines to within their own range")
	}

	// Should we throw an error if there's trailing stuff?
	if e = buffer.Insert(dest+offset, lines); e != nil {
		return
	}
	if cmd == 'm' {
		e = buffer.Delete(r)
	} // else 't'
	return
}

func cmdCopy(ctx *Context) (e error) {
	var r [2]int
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}
	return buffer.Copy(r)
}

func cmdPaste(ctx *Context) (e error) {
	var addr int
	// this is a bit hacky, but we're supposed to allow 0
	offset := 1
	last := len(ctx.addrs) - 1
	if ctx.addrs[last] == -1 {
		ctx.addrs[last] = 0
		offset = 0
	}
	if addr, e = buffer.AddrValue(ctx.addrs); e != nil {
		return
	}
	return buffer.Paste(addr + offset)
}

func cmdPrompt(ctx *Context) (e error) {
	if state.prompt {
		state.prompt = false
	} else if len(fprompt) > 0 {
		state.prompt = true
	}
	return
}

var (
	rxSanitize        = regexp.MustCompile(`\\.`)
	rxBackrefSanitize = regexp.MustCompile(`\\\\`)
	rxBackref         = regexp.MustCompile(`\\([0-9]+)|&`)
	rxSubArgs         = regexp.MustCompile(`g|l|n|p|\d+`)
)

// FIXME: this is probably more convoluted than it needs to be
func cmdSub(ctx *Context) (e error) {
	cmd := ctx.cmd[ctx.cmdOffset+1:]
	if len(cmd) == 0 {
		if len(state.lastSub) == 0 {
			return fmt.Errorf("invalid substitution")
		}
		cmd = state.lastSub
	}
	state.lastSub = cmd
	del := cmd[0]
	switch del {
	case ' ':
		fallthrough
	case '\n':
		fallthrough
	case 'm':
		fallthrough
	case 'g':
		return fmt.Errorf("invalid pattern delimiter")
	}
	// we replace escapes and their escaped characters with spaces to keep indexing
	sane := rxSanitize.ReplaceAllString(cmd, "  ")

	idx := [2]int{-1, -1}
	idx[0] = strings.Index(sane[1:], string(del)) + 1
	if idx[0] != -1 {
		idx[1] = strings.Index(sane[idx[0]+1:], string(del)) + idx[0] + 1
	}
	if idx[1] == -1 {
		idx[1] = len(cmd) - 1
	}

	mat := cmd[1:idx[0]]
	rep := cmd[idx[0]+1 : idx[1]]
	if rep == "%" {
		rep = state.lastRep
	}
	state.lastRep = rep
	arg := cmd[idx[1]+1:]

	// arg processing
	count := 1
	var printP, printL, printN, global bool

	parsedArgs := rxSubArgs.FindAllStringSubmatch(arg, -1)
	for _, m := range parsedArgs {
		switch m[0] {
		case "g":
			global = true
		case "p":
			printP = true
		case "l":
			printL = true
		case "n":
			printN = true
		default:
			if count, e = strconv.Atoi(m[0]); e != nil || count < 1 {
				return fmt.Errorf("invalid substitution argument")
			}
		}
	}

	repSane := rxBackrefSanitize.ReplaceAllString(rep, "  ")
	refs := rxBackref.FindAllStringSubmatchIndex(repSane, -1)

	var r [2]int
	if r, e = buffer.AddrRangeOrLine(ctx.addrs); e != nil {
		return
	}

	var rx *regexp.Regexp
	if rx, e = regexp.Compile(mat); e != nil {
		return
	}

	last := ""
	lastN := 0
	nMatch := 0
	b, _ := buffer.Get(r)
	// we have to do things a bit manually because we we only have ReplaceAll, and we don't necessarily want that
	for ln, l := range b {
		matches := rx.FindAllStringSubmatchIndex(l, -1)
		if !(len(matches) > 0) {
			continue // skip the rest if we don't have matches
		}
		if !global {
			if len(matches) >= count {
				matches = [][]int{matches[count-1]}
			} else {
				matches = [][]int{}
			}
		}
		// we have matches, deal with them
		fLin := ""
		oLin := 0
		for _, m := range matches {
			nMatch++
			// fRep := rep
			// offset := 0

			// Fill backrefs
			oRep := 0
			fRep := ""
			for _, r := range refs {
				if rep[r[0]:r[1]] == "&" {
					fRep += rep[oRep:r[0]]
					fRep += l[m[0]:m[1]]
					oRep = r[1]
				} else {
					i, _ := strconv.Atoi(rep[r[2]:r[3]])
					if i > len(m)/2-1 { // not enough submatches for backref
						return fmt.Errorf("invalid backref")
					}
					fRep += rep[oRep:r[0]]
					fRep += l[m[2*i]:m[2*i+1]]
					oRep = r[1]
				}
			}
			fRep += rep[oRep:]

			fLin += l[oLin:m[0]]
			fLin += fRep
			oLin = m[1]
		}
		fLin += l[oLin:]
		if e = buffer.Delete([2]int{ln, ln}); e != nil {
			return
		}
		if e = buffer.Insert(ln, []string{fLin}); e != nil {
			return
		}
		last = fLin
		lastN = ln
	}
	if nMatch == 0 {
		e = fmt.Errorf("no match")
	} else {
		if printP {
			fmt.Fprintf(ctx.out, "%s\n", last)
		}
		if printL {
			fmt.Fprintf(ctx.out, "%s$\n", last)
		}
		if printN {
			fmt.Fprintf(ctx.out, "%d\t%s\n", lastN+1, last)
		}
	}
	return
}

func cmdUndo(ctx *Context) (e error) {
	buffer.Rewind()
	return
}

func cmdDump(ctx *Context) (e error) {
	fmt.Fprintf(ctx.out, "%v\n", buffer)
	return
}

var rxCmdSub = regexp.MustCompile(`%`)

func cmdCommand(ctx *Context) (e error) {
	s := System{
		Cmd:    ctx.cmd[ctx.cmdOffset+1:],
		Stdin:  os.Stdin,
		Stdout: ctx.out,
		Stderr: os.Stderr,
	}
	e = s.Run()
	if e != nil {
		return
	}
	fmt.Fprintf(ctx.out, "!")
	return
}
