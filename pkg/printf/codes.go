package printf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type formatter interface {
	format(f *format, w *bytes.Buffer, arg []byte) error
}

var codeMap = map[rune]formatter{
	'%': newFormatCode[[]byte](func(f *format, w io.Writer, b []byte) error {
		w.Write([]byte("%"))
		return nil
	}, nIdentity),
	's': newFormatCode[[]byte](func(f *format, w io.Writer, b []byte) error {
		fmt.Fprintf(w, "%s", b)
		return nil
	}, nIdentity),
	'q': newFormatCode[[]byte](func(f *format, w io.Writer, b []byte) error {
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
	}, nIdentity),
	'd': newFormatCode[int](
		func(f *format, w io.Writer, i int) error {
			// TODO: properly encode the width and precision
			_, err := fmt.Fprintf(w, "%d", i)
			return err
		},
		func(f *format, b []byte) (int, bool) {
			if !bytes.HasPrefix(b, []byte("0x")) {
				return 0, false
			}
			// attempt to read as hex
			i, err := strconv.ParseInt(string(bytes.TrimPrefix(b, []byte("0x"))), 16, 64)
			if err != nil {
				return 0, true
			}
			return int(i), true
		},
		func(f *format, b []byte) (int, bool) {
			return readDecimal(bytes.NewBuffer(b))
		},
	),
	'i': nil,
	'o': nil,
	'u': nil,
	'x': nil,
	'X': nil,
	'f': nil,
	'e': nil,
	'E': nil,
	'g': nil,
	'G': nil,
	'c': nil,
}

func init() {
	var formatCodeB = newFormatCode[[]byte](func(f *format, w io.Writer, b []byte) error {
		tmp := &bytes.Buffer{}
		err := interpret(tmp, b, nil, true, false)
		if err != nil {
			return err
		}
		tmp.WriteTo(w)
		return nil
	}, nIdentity)
	// to avoid an intialization loop, we need to do this.
	codeMap['b'] = formatCodeB
}
