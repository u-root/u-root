package printf

import (
	"bytes"
	"fmt"
	"strconv"
)

func interpret(w *bytes.Buffer, format []byte, args []string, octalPrefix bool, parseSubstitutions bool) error {
	o := w
	fr := bytes.NewBuffer(format)
	idx := 0
	nextArg := func() []byte {
		if idx >= len(args) {
			return nil
		}
		ans := args[idx]
		idx = idx + 1
		return []byte(ans)
	}

	for fr.Len() > 0 {
		c, _, err := fr.ReadRune()
		// errors are only EOF
		if err != nil {
			continue
		}
		if c == '%' && parseSubstitutions {
			// at this point we are looking for which format code this is
			// read another rune
			n, _, err := fr.ReadRune()
			// error only EOF, so write the original rune and continue
			if err != nil {
				o.WriteRune(c)
				continue
			}
			arg := nextArg()
			switch n {
			case '%':
				o.WriteRune('%')
				continue
			case 'b':
				tmp := &bytes.Buffer{}
				err := interpret(tmp, arg, nil, true, false)
				if err != nil {
					return err
				}
				o.WriteString(tmp.String())
				continue
			case 'q':
				if err := formatCodeQ.format(o, arg); err != nil {
					return err
				}
				continue
			case 'd':
				continue
			case 'i':
				continue
			case 'o':
				continue
			case 'u':
				continue
			case 'x':
				continue
			case 'X':
				continue
			case 'f':
				continue
			case 'e':
				continue
			case 'E':
				continue
			case 'g':
				continue
			case 'G':
				continue
			case 'c':
				continue
			case 's':
				fmt.Fprintf(o, "%s", arg)
				continue
			default:
				return NewErrInvalidDirective(string(n))
			}
		} else if c == '\\' {
			// at this point we are looking for which escape sequence this is
			// read another rune
			n, _, err := fr.ReadRune()
			// error only EOF, so write the original rune and continue
			if err != nil {
				o.WriteRune(c)
				continue
			}
			switch n {
			case '0':
				if octalPrefix {
					readOctal(fr, o)
					continue
				}
			case '"':
				o.WriteRune(n)
				continue
			case '\\':
				o.WriteRune(n)
				continue
			case 'a':
				o.WriteRune('\a')
				continue
			case 'b':
				o.WriteRune('\b')
				continue
			case 'c': //produce no further input
				return nil
			case 'e': //escape
				o.WriteRune(27)
				continue
			case 'f':
				o.WriteRune('\f')
				continue
			case 'n':
				o.WriteRune('\n')
				continue
			case 'r':
				o.WriteRune('\r')
				continue
			case 't':
				o.WriteRune('\t')
				continue
			case 'v':
				o.WriteRune('\v')
				continue
			case 'x':
				readUnicode(fr, o, 2)
				continue
			case 'u':
				readUnicode(fr, o, 4)
				continue
			case 'U':
				readUnicode(fr, o, 8)
				continue
			}
			// if not octal prefix, and its a decimal, then unread the rune and run readOctal
			if !octalPrefix && n >= '0' && n <= '9' {
				fr.UnreadRune()
				readOctal(fr, o)
				continue
			}
			// there's no match, so just write the full sequence
			o.WriteRune('\\')
			o.WriteRune(n)
		} else {
			// otherwise just write the rune
			o.WriteRune(c)
		}
	}
	return nil
}

func readOctal(fr *bytes.Buffer, o *bytes.Buffer) {
	octals := ""
	// read up to three decimals from the stream
	for i := 0; i < 3; i++ {
		dec, _, err := fr.ReadRune()
		if err != nil {
			break
		}
		if dec >= '0' && dec <= '9' {
			octals = octals + string(dec)
			continue
		}
		// not an decimal, so unwind the rune and also break, printing the octal if it exists
		fr.UnreadRune()
		break
	}
	// if the length of octals is zero, that means this is not actually a format code.
	// for instance, if the input is `\0\`, this would ensure that we properly print the ending \
	if len(octals) == 0 {
		o.WriteRune('\\')
	} else {
		i, err := strconv.ParseInt(octals, 8, 32)
		if err == nil {
			o.WriteRune(rune(i))
		}
	}
}

func readUnicode(fr *bytes.Buffer, o *bytes.Buffer, length int) {
	hexcode := ""
	for i := 0; i < length; i++ {
		dec, _, err := fr.ReadRune()
		if err != nil {
			break
		}
		if (dec >= '0' && dec <= '9') || (dec >= 'a' && dec <= 'f') || (dec >= 'A' && dec <= 'F') {
			hexcode = hexcode + string(dec)
			continue
		}
		// not an decimal, so unwind the rune and also break, printing the octal if it exists
		fr.UnreadRune()
		break
	}
	if len(hexcode) != 0 {
		i, err := strconv.ParseInt(hexcode, 16, 32)
		if err == nil {
			o.WriteRune(rune(i))
		}
	}
}
