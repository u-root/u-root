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
	args  []arg
	fdmap map[int]string
	link  string
	bg    bool

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
	err := b.UnreadByte()
	if err != nil {
		panic(fmt.Sprintf("unreading bufio: %v", err))
	}
}

func one(b *bufio.Reader) byte {
	//fmt.Printf("next\n")
	c, err := b.ReadByte()
	//fmt.Printf("'%v' %v\n", c, err)
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
		return one(b)
	}
	return byte(c)
}

// Tokenize stuff coming in from the stream. For everything but an arg, the
// type is just the thing itself, since we can switch on strings.
func tok(b *bufio.Reader) (string, string) {
	tokType, arg := "white", ""
	c := next(b)

	//fmt.Printf("TOK %v", c)
	switch c {
	case 0:
		return "EOF", ""
	case '>':
		return "FD", "1"
	case '<':
		return "FD", "0"
	// yes, I realize $ handling is still pretty hokey.
	case '$':
		arg = ""
		c = next(b)
		for {
			if strings.Index(punct, string(c)) > -1 {
				pushback(b)
				break
			}
			arg = arg + string(c)
			c = next(b)
		}
		return "ENV", arg
	case '\'':
		for {
			nc := next(b)
			if nc == '\'' {
				return "ARG", arg
			}
			arg = arg + string(nc)
		}
	case ' ':
		return "white", string(c)
	case '\n':
		//fmt.Printf("NEWLINE\n")
		return "EOL", ""
	case '|', '&':
		//fmt.Printf("LINK %v\n", c)
		// peek ahead. We need the literal, so don't use next()
		nc := one(b)
		if nc == c {
			//fmt.Printf("LINK %v\n", string(c)+string(c))
			return "LINK", string(c) + string(c)
		}
		pushback(b)
		if c == '&' {
			//fmt.Printf("BG\n")
			tokType = "BG"
			if nc == 0 {
				tokType = "EOL"
			}
			return "BG", tokType
		}
		//fmt.Printf("LINK %v\n", string(c))
		return "LINK", string(c)
	default:
		for {
			if strings.Index(punct, string(c)) > -1 {
				pushback(b)
				return "ARG", arg
			}
			arg = arg + string(c)
			c = next(b)
		}

	}
	return tokType, arg

}

// get an ARG. It has to work.
func getArg(b *bufio.Reader, what string) string {
	for {
		nt, s := tok(b)
		if nt == "EOF" || nt == "EOL" {
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
func parse(b *bufio.Reader) (*Command, string) {
	t, s := tok(b)
	//fmt.Printf("%v %v\n", t, s)
	// Cover the trivial case that nothing happens.
	if s == "\n" || t == "EOF" || t == "EOL" {
		return nil, t
	}
	//fmt.Printf("%v %v\n", t, s)
	c := newCommand()
	for {
		switch t {
		case "ENV", "ARG":
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
			//fmt.Printf("LINK %v %v\n", c, s)
			return c, t
		case "BG":
			c.bg = true
			return c, t
		case "EOF":
			return c, t
		case "EOL":
			return c, t
		default:
			panic(fmt.Sprintf("unknown token type %v", t))
		}
		t, s = tok(b)
	}
	return c, t
}

func newCommand() *Command {
	return &Command{fdmap: make(map[int]string)}
}

// Just eat it up until you have all the commands you need.
func parsecommands(b *bufio.Reader) ([]*Command, string) {
	cmds := make([]*Command, 0)
	for {
		c, t := parse(b)
		if c == nil {
			return cmds, t
		}
		//fmt.Printf("cmd  %v\n", *c)
		cmds = append(cmds, c)
		if t == "EOF" || t == "EOL" {
			return cmds, t
		}
	}
}

func getCommand(b *bufio.Reader) ([]*Command, string, error) {
	// TODO: put a recover here that just returns an error.
	c, t := parsecommands(b)
	// the rules.
	// For now, no empty commands.
	// Can't have a redir and a redirect for fd1.
	for i := range c {
		if len(c[i].args) == 0 {
			return nil, "", errors.New("empty commands not allowed (yet)\n")
		}
		if c[i].link == "|" && c[i].fdmap[1] != "" {
			return nil, "", errors.New("Can't have a pipe and > on one command\n")
		}
		if c[i].link == "|" && i == len(c)-1 {
			return nil, "", errors.New("Can't have a pipe to nowhere\n")
		}
		if i < len(c)-1 && c[i].link == "|" && c[i+1].fdmap[0] != "" {
			return nil, "", errors.New("Can't have a pipe to command with redirect on stdin\n")
		}
	}
	return c, t, nil
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
