package printf

import (
	"bufio"
	"bytes"
	"io"
)

// interpreting format code can be a two step process
// first the input type must be normalized, then printed
//
// for instance, `printf '%d' "0x1234"` will return 4660
// in this case, 0x1234 must first be "normalized" to the "%d"'s datatype
// then, the number should be formatted according to the rules of the format code

func newFormatCode[T any](formatter func(io.Writer, T) error, normalizers ...func([]byte) (T, bool)) *formatcode[T] {
	return &formatcode[T]{
		formatter:   formatter,
		normalizers: normalizers,
	}
}

type formatcode[T any] struct {
	normalizers []func([]byte) (T, bool)
	formatter   func(io.Writer, T) error
}

func (c *formatcode[T]) format(w *bytes.Buffer, arg []byte) error {
	var val T
	for _, v := range c.normalizers {
		tmp, ok := v(arg)
		if ok {
			val = tmp
			break
		}
	}
	return c.formatter(w, val)
}

func nIdentity(xs []byte) ([]byte, bool) {
	return xs, true
}

var formatCodeQ = newFormatCode[[]byte](func(w io.Writer, b []byte) error {
	wr := bufio.NewWriterSize(w, len(b))
	defer wr.Flush()
	rd := bytes.NewBuffer(b)
	for {
		char, n, err := rd.ReadRune()
		if n > 0 {
			if isMetaCharacter(char) {
				wr.WriteByte('\\')
			}
			wr.WriteRune(char)
		}
		if err != nil {
			return nil
		}
	}
}, nIdentity)

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
