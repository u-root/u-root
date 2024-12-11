// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The u-root shell is intended to be very simple, since builtins and extensions
// are written in Go. It should not need YACC. As in the JSON parser, we hope this
// simple state machine will do the job.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type arg struct {
	val string
	mod string
}

// The Command struct is initially filled in by the parser. The shell itself
// adds to it as processing continues, and then uses it to creates os.Commands
type Command struct {
	*exec.Cmd
	// These are filled in by the parser.
	Args  []arg
	fdmap map[int]string
	files map[int]io.Closer
	Link  string
	BG    bool

	// These are set up by the shell as it evaluates the Commands
	// provided by the parser.
	// we separate the command so people don't have to put checks for the length
	// of argv in their builtins. We do that for them.
	cmd  string
	argv []string
}

var (
	cmds  []Command
	punct = "<>|&$ \t\n"
)

func pushback(b *bufio.Reader) {
	// If we can't UnreadByte it will get us an obscure error, but
	// it is hard to tell what else to do.
	_ = b.UnreadByte()
}

func one(b *bufio.Reader) byte {
	c, err := b.ReadByte()
	// On any kind of error, just return 0.
	// This is not a serious kind of shell any more, it
	// is more for diagnostics when things are really broken,
	// so we need not be too picky about errors.
	if err != nil {
		return 0
	}
	return c
}

func next(b *bufio.Reader) byte {
	c := one(b)
	if c == '\\' {
		return one(b)
	}
	return byte(c)
}

// Tokenize stuff coming in from the stream. For everything but an arg, the
// type is just the thing itself, since we can switch on strings.
func tok(b *bufio.Reader) (string, string) {
	var tokType, arg string
	c := next(b)

	switch c {
	case 0, 4:
		return "EOF", ""
	case '>':
		return "FD", "1"
	case '<':
		return "FD", "0"
	// yes, I realize $ handling is still pretty hokey.
	// And, again, this is a diagnostic tool now, not a general
	// purpose shell, so the hell with it.
	case '$':
		arg = ""
		return "ENV", arg
	case '\'':
		for {
			nc := next(b)
			if nc == '\'' {
				return "ARG", arg
			}
			arg = arg + string(nc)
		}
	case ' ', '\t':
		return "white", string(c)
	case '\n', '\r':
		return "EOL", ""
	case '|', '&':
		// peek ahead. We need the literal, so don't use next()
		nc := one(b)
		if nc == c {
			return "LINK", string(c) + string(c)
		}
		if nc != 0 {
			pushback(b)
		}
		if c == '&' {
			tokType = "BG"
			if nc == 0 {
				tokType = "EOL"
			}
			return "BG", tokType
		}
		return "LINK", string(c)
	default:
		for {
			if c == 0 {
				return "ARG", arg
			}
			if strings.Contains(punct, string(c)) {
				pushback(b)
				return "ARG", arg
			}
			arg = arg + string(c)
			c = next(b)
		}

	}
}

// get an ARG. It has to work.
func getArg(b *bufio.Reader, what string) string {
	for {
		nt, s := tok(b)
		if nt == "EOF" || nt == "EOL" {
			// We used to panic here, but what's the sense in that?
			return ""
		}
		if nt == "white" {
			continue
		}
		// It has to work, but if not ... too bad.
		if nt != "ARG" {
			return ""
		}
		return s
	}
}

func parsestring(b *bufio.Reader, c *Command) (*Command, string) {
	t, s := tok(b)
	if s == "\n" || t == "EOF" || t == "EOL" {
		return nil, t
	}
	for {
		switch t {
		// In old rush, env strings were substituted wholesale, and
		// parsed in line. This was very useful, but nobody uses it so...
		// for now, screw it. Do environment later.
		case "ENV":
		case "ARG":
			c.Args = append(c.Args, arg{s, t})
		case "white":
		case "FD":
			x := 0
			_, err := fmt.Sscanf(s, "%v", &x)
			if err != nil {
				panic(fmt.Errorf("bad FD on redirect: %v, %w", s, err))
			}
			// whitespace is allowed
			c.fdmap[x] = getArg(b, t)
		// LINK and BG are similar save that LINK requires another command. If we don't get one, well.
		case "LINK":
			c.Link = s
			return c, t
		case "BG":
			c.BG = true
			return c, t
		case "EOF":
			return c, t
		case "EOL":
			return c, t
		default:
			panic(fmt.Errorf("unknown token type %v", t))
		}
		t, s = tok(b)
	}
}

func parse(b *bufio.Reader) (*Command, string) {
	c := newCommand()
	return parsestring(b, c)
}

func newCommand() *Command {
	return &Command{Cmd: exec.Command(""), fdmap: make(map[int]string), files: make(map[int]io.Closer)}
}

// Just eat it up until you have all the commands you need.
func parsecommands(b *bufio.Reader) ([]*Command, string) {
	cmds := make([]*Command, 0)
	for {
		c, t := parse(b)
		if c == nil {
			return cmds, t
		}
		cmds = append(cmds, c)
		if t == "EOF" || t == "EOL" {
			return cmds, t
		}
	}
}

func getCommand(b *bufio.Reader) (c []*Command, t string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	// TODO: put a recover here that just returns an error.
	c, t = parsecommands(b)
	// the rules.
	// For now, no empty commands.
	// Can't have a redir and a redirect for fd1.
	for i, v := range c {
		if len(v.Args) == 0 {
			return nil, "", errors.New("empty commands not allowed (yet)")
		}
		if v.Link == "|" && v.fdmap[1] != "" {
			return nil, "", errors.New("can't have a pipe and > on one command")
		}
		if v.Link == "|" && i == len(c)-1 {
			return nil, "", errors.New("can't have a pipe to nowhere")
		}
		if i < len(c)-1 && v.Link == "|" && c[i+1].fdmap[0] != "" {
			return nil, "", errors.New("can't have a pipe to command with redirect on stdin")
		}
	}
	return c, t, err
}

/*
func main() {
	b := bufio.NewReader(os.Stdin)
	for {
	    c, t, err := getCommand(b)
		fmt.Printf("%v %v %v\n", c, t, err)
	    if t == "EOF" {
	       break
	       }
	       }
}
*/
