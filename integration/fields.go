// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

func isWhitespace(b byte) bool {
	return b == '\t' || b == '\n' || b == '\v' ||
		b == '\f' || b == '\r' || b == ' '
}

// fields splits the string s around each instance of one or more consecutive white space
// characters, returning a slice of substrings of s or an
// empty slice if s contains only white space.
//
// fields is similar to strings.Fields() method, two main differences are:
//     fields doesn't split substring of s if substring is inside of double quotes
//     fields works only with ASCII strings.
func fields(s string) []string {
	var ret []string
	var token []byte

	var quotes bool
	for i := range s {
		if s[i] == '"' {
			quotes = !quotes
		}

		if !isWhitespace(s[i]) || quotes {
			token = append(token, s[i])
		} else if len(token) > 0 {
			ret = append(ret, string(token))
			token = token[:0]
		}
	}

	if len(token) > 0 {
		ret = append(ret, string(token))
	}
	return ret
}
