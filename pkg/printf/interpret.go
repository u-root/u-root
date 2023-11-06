package printf

import (
	"bytes"
	"fmt"
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
			// these codes are ignored because if we dont find a valid argument, we will just set it to 0
			// perhaps better defaults or errors could be possibly set instead, but keeping it like this for now
			if format.width == -1 {
				format.width, _ = readDecimal(bytes.NewBuffer(nextArg()))
			}
			if format.precision == -1 {
				format.precision, _ = readDecimal(bytes.NewBuffer(nextArg()))
			}
			// see if the format code is implemented :)
			formatter, ok := codeMap[format.specifier]
			if !ok || formatter == nil {
				return fmt.Errorf("%w: %s", ErrUnimplemented, string(format.specifier))
			}
			// run the formatter
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
