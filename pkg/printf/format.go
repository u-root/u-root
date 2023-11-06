package printf

import (
	"bytes"
	"fmt"
	"io"
)

// format specification is
// %[flags][width][.precision][length]specifier,
type format struct {
	leftJustify bool
	precedeSign bool
	blankSpace  bool
	withPound   bool
	padZero     bool

	width     int
	precision int
	length    int
	specifier rune
}

func readFormat(fr *bytes.Buffer) (o *format, err error) {
	o = &format{
		precision: 1, // the default precision is always 1
	}
	stage := 0
	for {
		n, _, err := fr.ReadRune()
		if err != nil {
			break
		}
		if stage <= 0 {
			// read flags
			switch n {
			case '-':
				o.leftJustify = true
				continue
			case '+':
				o.precedeSign = true
				continue
			case ' ':
				o.blankSpace = true
				continue
			case '#':
				o.withPound = true
				continue
			case '0':
				o.padZero = true
				continue
			default:
				stage = 1
			}
		}

		if stage <= 1 {
			stage = 2
			if isDecimal(n) {
				fr.UnreadRune()
				// dont need to check ok, since know its a decimal
				width, _ := readDecimal(fr)
				o.width = width
				stage = 2
				continue
			}
			// didnt read a decimal, so check if n == '*'
			if n == '*' {
				o.width = -1
				stage = 2
				continue
			}
		}
		// read precision
		if stage <= 2 {
			stage = 3
			switch n {
			case '.':
				ma, _, _ := fr.ReadRune()
				if ma == '*' {
					o.precision = -1
				} else {
					fr.UnreadRune()
					o.precision, _ = readDecimal(fr)
				}
				continue
			default:
				// if it's not a period, then we can keep going
			}
		}

		if stage <= 3 {
			stage = 4
			switch n {
			case 'h':
				o.length = 1
				continue
			case 'l':
				o.length = 2
				continue
			case 'L':
				o.length = 3
				continue
			default:
			}
		}
		// try to read a specifier
		switch n {
		case '%', 'b', 'q', 'd', 'i', 'o', 'u', 'x', 'X', 'f', 'e', 'E', 'g', 'G', 'c', 's':
			o.specifier = n
			return o, nil
		default:
			return o, fmt.Errorf("%w: %s", ErrInvalidDirective, string(n))
		}
	}
	return o, nil
}

// interpreting format code can be a two step process
// first the input type must be normalized, then printed
//
// for instance, `printf '%d' "0x1234"` will return 4660
// in this case, 0x1234 must first be "normalized" to the "%d"'s datatype
// then, the number should be formatted according to the rules of the format code

func newFormatCode[T any](formatter func(*format, io.Writer, T) error, normalizers ...func(*format, []byte) (T, bool)) *formatcode[T] {
	return &formatcode[T]{
		formatter:   formatter,
		normalizers: normalizers,
	}
}

type formatcode[T any] struct {
	normalizers []func(*format, []byte) (T, bool)
	formatter   func(*format, io.Writer, T) error
}

func (c *formatcode[T]) format(f *format, w *bytes.Buffer, arg []byte) error {
	var val T
	for _, v := range c.normalizers {
		tmp, ok := v(f, arg)
		if ok {
			val = tmp
			break
		}
	}
	return c.formatter(f, w, val)
}

func nIdentity(f *format, xs []byte) ([]byte, bool) {
	return xs, true
}

// this is very very arbitrary....
func isMetaCharacter(r rune) bool {
	switch r {
	case '|':
	case '&':
	case ';':
	case '<':
	case '>':
	case '(':
	case ')':
	case '$':
	case '`':
	case '\\':
	case '\'':
	case '"':
	case ' ':
	case '\t':
	case '\n':
	case '*':
	case '?':
	case '[':
	case ']':
	case '#':
	case '~':
	case '=':
	default:
		switch {
		default:
			return false
		}
	}
	return true
}
