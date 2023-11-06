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
			// we are now parsing a format code.
			format, err := readFormat(fr)
			if err != nil {
				return err
			}
			// now we have the formatCode, so we figure out how many args we need to take.
			if format.width == -1 {
				arg := nextArg()
				format.width, _ = readDecimal(bytes.NewBuffer(arg))
			}
			if format.precision == -1 {
				arg := nextArg()
				format.precision, _ = readDecimal(bytes.NewBuffer(arg))
			}

			formatter, ok := codeMap[format.specifier]
			if !ok {
				return fmt.Errorf("%w: %s", ErrUnimplemented, string(format.specifier))
			}
			err = formatter.format(format, o, nextArg())
			if err != nil {
				// formatter error, return it
				return err
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

func readDecimal(fr *bytes.Buffer) (int, bool) {
	octals := ""
	// read decimals until there are no more decimals to read
	// there is a sanity check at 256. more than that will not be supported for now
	for i := 0; i < 256; i++ {
		dec, _, err := fr.ReadRune()
		if err != nil {
			break
		}
		if dec >= '0' && dec <= '9' {
			octals = octals + string(dec)
			continue
		}
		// not an decimal, so unwind the last read and break
		fr.UnreadRune()
		break
	}
	// if the length of octals is zero, that means this is not actually a number. return nil
	if len(octals) == 0 {
		return 0, false
	}
	// read the base 10 number
	i, err := strconv.ParseInt(octals, 10, 64)
	// this should never actually error, since the input has been pre-sanitized... maybe overflow?
	if err != nil {
		return 0, true
	}
	return int(i), true
}

func isDecimal(r rune) bool {
	return r >= '0' && r <= '9'
}
