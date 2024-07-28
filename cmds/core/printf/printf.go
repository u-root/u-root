// Copyright 2013-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SYNOPSIS
//
//	printf format [arg ...]

// DESCRIPTION
//
//	printf writes formatted data according to format using each arg until
//	drained.
//
//	printf interprets the standard escape sequences \\, \', \", \a, \b, \e,
//	\f, \n, \r, \t, \v, \xH[H], \O[OOO], the sequence \c, which terminates
//	further output if it's found inside format or a %b format string, the
//	format specification %b for an unescaped string and all C printf(3)
//	format specifications ending with csdiouxXaAeEfFgG, including variable
//	width and precision.

// STANDARDS
//
//	The printf utility is compliant with the p1003.1-2013 specification.
//
//	The possibility of specifying 4-digit octals is an extension to that
//	specification.

// Author
//
// xplshn
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s format [arg ...]\n", os.Args[0])
	os.Exit(1)
}

func printf(format string, args []string) (string, error) {
	var output strings.Builder
	argi := 0

	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			if format[i] == '\\' && i+1 < len(format) {
				switch format[i+1] {
				case 'n':
					output.WriteString("\n")
				case 't':
					output.WriteString("\t")
				case '\\':
					output.WriteString("\\")
				case 'a':
					output.WriteString("\a")
				case 'b':
					output.WriteString("\b")
				case 'e':
					output.WriteString("\033")
				case 'f':
					output.WriteString("\f")
				case 'r':
					output.WriteString("\r")
				case 'v':
					output.WriteString("\v")
				case 'c':
					return output.String(), nil
				case 'x':
					if i+3 < len(format) && isHexDigit(format[i+2]) && isHexDigit(format[i+3]) {
						hex := format[i+2 : i+4]
						num, _ := strconv.ParseInt(hex, 16, 32)
						if num < 0 || num > 255 {
							return "", fmt.Errorf("%%x: value out of uint8 range")
						}
						output.WriteByte(byte(num))
						i += 3
						continue
					}
				case '0', '1', '2', '3', '4', '5', '6', '7':
					if i+3 < len(format) && isOctalDigit(format[i+2]) && isOctalDigit(format[i+3]) {
						oct := format[i+1 : i+4]
						num, _ := strconv.ParseInt(oct, 8, 32)
						if num < 0 || num > 255 {
							return "", fmt.Errorf("%%O: value out of uint8 range")
						}
						output.WriteByte(byte(num))
						i += 3
						continue
					}
				default:
					output.WriteByte(format[i])
				}
				i++
				continue
			}
			output.WriteByte(format[i])
			continue
		}

		i++
		if i >= len(format) {
			break
		}

		width := -1
		precision := -1

		// Parse width
		if format[i] == '*' {
			if argi < len(args) {
				width, _ = strconv.Atoi(args[argi])
				argi++
			}
			i++
		} else {
			j := i
			for ; j < len(format) && format[j] >= '0' && format[j] <= '9'; j++ {
			}
			if j > i {
				width, _ = strconv.Atoi(format[i:j])
				i = j
			}
		}

		// Parse precision
		if i < len(format) && format[i] == '.' {
			i++
			if i < len(format) && format[i] == '*' {
				if argi < len(args) {
					precision, _ = strconv.Atoi(args[argi])
					argi++
				}
				i++
			} else {
				j := i
				for ; j < len(format) && format[j] >= '0' && format[j] <= '9'; j++ {
				}
				if j > i {
					precision, _ = strconv.Atoi(format[i:j])
					i = j
				}
			}
		}

		switch format[i] {
		case 's':
			if argi < len(args) {
				s := args[argi]
				if precision >= 0 && precision < len(s) {
					s = s[:precision]
				}
				if width > 0 {
					output.WriteString(fmt.Sprintf(fmt.Sprintf("%%%ds", width), s))
				} else {
					output.WriteString(s)
				}
				argi++
			} else {
				output.WriteString("")
			}
		case 'd', 'i':
			if argi < len(args) {
				num, err := strconv.Atoi(args[argi])
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid number '%s'\n", args[argi])
					output.WriteString("0")
					argi++
					continue
				}
				if width > 0 {
					output.WriteString(fmt.Sprintf(fmt.Sprintf("%%%dd", width), num))
				} else {
					output.WriteString(strconv.Itoa(num))
				}
				argi++
			} else {
				output.WriteString("0")
			}
		case 'f':
			if argi < len(args) {
				num, err := strconv.ParseFloat(args[argi], 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid number '%s'\n", args[argi])
					output.WriteString("0.000000")
					argi++
					continue
				}
				if precision >= 0 {
					output.WriteString(fmt.Sprintf(fmt.Sprintf("%%.%df", precision), num))
				} else {
					output.WriteString(fmt.Sprintf("%.6f", num))
				}
				argi++
			} else {
				output.WriteString("0.000000")
			}
		case 'b':
			if argi < len(args) {
				unescaped, err := unescape(args[argi])
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid escape sequence '%s'\n", args[argi])
					output.WriteString("")
					argi++
					continue
				}
				output.WriteString(unescaped)
				argi++
			} else {
				output.WriteString("")
			}
		case '%':
			output.WriteString("%")
		default:
			return "", fmt.Errorf("%%: invalid format")
		}
	}

	return output.String(), nil
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

func unescape(s string) (string, error) {
	var output strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				output.WriteString("\n")
			case 't':
				output.WriteString("\t")
			case '\\':
				output.WriteString("\\")
			case 'a':
				output.WriteString("\a")
			case 'b':
				output.WriteString("\b")
			case 'e':
				output.WriteString("\033")
			case 'f':
				output.WriteString("\f")
			case 'r':
				output.WriteString("\r")
			case 'v':
				output.WriteString("\v")
			case 'x':
				if i+3 < len(s) && isHexDigit(s[i+2]) && isHexDigit(s[i+3]) {
					hex := s[i+2 : i+4]
					num, _ := strconv.ParseInt(hex, 16, 32)
					if num < 0 || num > 255 {
						return "", fmt.Errorf("%%x: value out of uint8 range")
					}
					output.WriteByte(byte(num))
					i += 3
					continue
				}
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if i+3 < len(s) && isOctalDigit(s[i+2]) && isOctalDigit(s[i+3]) {
					oct := s[i+1 : i+4]
					num, _ := strconv.ParseInt(oct, 8, 32)
					if num < 0 || num > 255 {
						return "", fmt.Errorf("%%O: value out of uint8 range")
					}
					output.WriteByte(byte(num))
					i += 3
					continue
				}
			default:
				output.WriteByte(s[i])
			}
			i++
		} else {
			output.WriteByte(s[i])
		}
	}
	return output.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	format := os.Args[1]
	args := os.Args[2:]

	output, err := printf(format, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Print(output)
}
