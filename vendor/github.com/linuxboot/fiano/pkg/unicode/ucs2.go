// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unicode converts between UCS2 cand UTF8.
package unicode

import (
	"log"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// UCS2ToUTF8 converts from UCS2 to UTF8.
func UCS2ToUTF8(input []byte) string {
	e := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	output, _, err := transform.Bytes(e.NewDecoder(), input)
	if err != nil {
		log.Printf("could not decode UCS2: %v", err)
		return string(input)
	}
	// Remove null terminator if one exists.
	if output[len(output)-1] == 0 {
		output = output[:len(output)-1]
	}
	return string(output)
}

// UTF8ToUCS2 converts from UTF8 to UCS2.
func UTF8ToUCS2(input string) []byte {
	e := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	input = input + "\000" // null terminator
	output, _, err := transform.Bytes(e.NewEncoder(), []byte(input))
	if err != nil {
		log.Printf("could not encode UCS2: %v", err)
		return []byte(input)
	}
	return output
}
