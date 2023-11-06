package printf

import (
	"bytes"
	"strconv"
)

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
