// The u-root shell is intended to be very simple, since builtins and extensions
// are written in Go. It should not need YACC. As in the JSON parser, we hope this
// simple state machine will do the job.
package main

import (
       "bufio"
	"fmt"
	"io"
	"os"
)


type arg struct {
	val string
	modifier string
}

type Command struct {
	args []arg
	cmd string
	fdmap map[int]string
	link string
	bg bool
}

var (
	cmds []Command
)

func pushback(b *bufio.Reader) {
	err := b.UnreadByte()
	if err != nil {
		panic(fmt.Sprintf("unreading bufio: %v", err))
	}
}

func one(b *bufio.Reader) byte {
	fmt.Printf("next\n")
	c, err := b.ReadByte()
	fmt.Printf("'%v' %v\n", c, err)
	if err == io.EOF {
		return 0
	}
	if err != nil {
		panic(fmt.Sprintf("reading bufio: %v", err))
	}
	return c
}

func next(b *bufio.Reader) byte {
	c := one(b)
	if c == '\\' {
		return next(b)
	}
	return byte(c)
}
	
// Tokenize stuff coming in from the stream. For everything but an arg, the 
// type is just the thing itself, since we can switch on strings.
func tok(b *bufio.Reader) (string, string, bool) {
	tokType, arg := "white", ""
	c := next(b)

fmt.Printf("TOK %v", c)
	switch(c) {
		case 0:
			return "EOF", "", true
		case '>':
			return "FD", "1", false
		case '<':
			return "FD", "0", false
		case '\'': 
			for {
				nc := next(b)
				if nc == '\'' {
					return "ARG", arg, false
				}
				arg = arg + string(nc)
			}
		case ' ': 
			return "white", string(c), false
		case '\n': 
fmt.Printf("NEWLINE\n")
			return "EOF", "", true
		case '|', '&':
			fmt.Printf("LINK %v\n", c)
			// peek ahead. We need the literal, so don't use next()
			nc := one(b)
			if nc == c {
				fmt.Printf("LINK %v\n", string(c)+string(c))
				return "LINK", string(c)+string(c), false
			}
			pushback(b)
			if c == '&' {
				fmt.Printf("BG\n")
				return "BG", "&", nc == 0
			}
			fmt.Printf("LINK %v\n", string(c))
			return "LINK", string(c), false
		default:
			for {
				if c == ' ' || c == '\n' {
					pushback(b)
					return "ARG", arg, false
				}
				arg = arg + string(c)
				c = next(b)
			}
		
	}
	return tokType, arg, false
	
}

// get an ARG. It has to work.
func getArg(b *bufio.Reader, what string) string {
			for {
				nt, s, eof := tok(b)
				if eof {
					panic(fmt.Sprintf("%v requires an argument", what))
				}
				if nt == "white" {
					continue
				}
				if nt != "ARG" {
					panic(fmt.Sprintf("%v requires an argument, not %v", what, nt))
				}
				return s
			}
}
func parse(b *bufio.Reader) (*Command, bool) {
	t, s, eof := tok(b)
	// Cover the trivial case that nothing happens.
	if s == "\n" || eof {
		return nil, eof
	}
	fmt.Printf("%v %v %v\n", t, s, eof)
	c := newCommand()
	for {
	switch(t) {
		case "ARG": 
			c.args = append(c.args, arg{s, t})
		case "white":
		case "FD":
			x := 0
			_, err := fmt.Sscanf(s, "%v", &x)
			if err != nil {
				panic(fmt.Sprintf("bad FD on redirect: %v, %v", s, err))
			}
			// whitespace is allowed
			c.fdmap[x] = getArg(b, t)
		// LINK and BG are similar save that LINK requires another command. If we don't get one, well.
		case "LINK":
			c.link = s
			fmt.Printf("LINK %v %v\n", c, s)
			return c, eof
		case "BG":
			c.bg = true
			return c, eof
		case "EOF":
			return c, true
		default:
			panic(fmt.Sprintf("unknown token type %v", t))
	}
	t, s, eof = tok(b)
	}
	return c, eof
}

func newCommand() *Command {
	return &Command{fdmap: make(map[int]string)}
}

// Just eat it up until you have all the commands you need.
func parsecommands(b *bufio.Reader) []*Command {
	cmds := make([]*Command, 0)
	for {
		c, eoc := parse(b)
		if c == nil {
			return cmds
		}
		fmt.Printf("cmd  %v\n", *c)
		cmds = append(cmds, c)
		if eoc {
			break
		}
	}
	fmt.Printf("cmds %v\n", cmds)
	return cmds
}

func main() {
	b := bufio.NewReader(os.Stdin)
	c := parsecommands(b)
	fmt.Printf("%v\n", c)
}
