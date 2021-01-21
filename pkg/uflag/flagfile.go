// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uflag supports u-root-custom flags as well as flag files.
package uflag

import (
	"fmt"
	"strconv"
	"strings"
)

// ArgvToFile encodes argv program arguments such that they can be stored in a
// file.
func ArgvToFile(args []string) string {
	// We separate flags in the flags file with new lines, so we
	// have to escape new lines.
	//
	// Go already has a nifty Quote mechanism which will escape
	// more than just new-line, which won't hurt anyway.
	var quoted []string
	for _, arg := range args {
		quoted = append(quoted, strconv.Quote(arg))
	}
	return strings.Join(quoted, "\n")
}

// FileToArgv converts argvs stored in a file back to an array of strings.
func FileToArgv(content string) []string {
	quotedArgs := strings.Split(content, "\n")
	var args []string
	for _, arg := range quotedArgs {
		if len(arg) > 0 {
			s, err := strconv.Unquote(arg)
			if err != nil {
				panic(fmt.Sprintf("flags file encoded wrong, arg %q, error %v", arg, err))
			}
			args = append(args, s)
		}
	}
	return args
}
