package main

import (
	"bytes"
	"strconv"
	"strings"
)

func interpretFormat(from string, octalPrefix bool) string {
	fr := strings.NewReader(from)
	o := &bytes.Buffer{}

	for fr.Len() > 0 {
		c, _, err := fr.ReadRune()
		// errors are only EOF
		if err != nil {
			continue
		}
		if c == '\\' {
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
				return o.String()
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
	return o.String()
}

func readOctal(fr *strings.Reader, o *bytes.Buffer) {
	octals := ""
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
	if len(octals) == 0 {
		o.WriteRune('\\')
	} else {
		i, err := strconv.ParseInt(octals, 8, 64)
		if err == nil {
			o.WriteRune(rune(i))
		}
	}
}
func readUnicode(fr *strings.Reader, o *bytes.Buffer, length int) {
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
	if len(hexcode) == 0 {
		o.WriteRune('\\')
	} else {
		i, err := strconv.ParseInt(hexcode, 16, 64)
		if err == nil {
			o.WriteRune(rune(i))
		}
	}
}
