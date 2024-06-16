// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unixflag

import (
	"os"
	"strings"
)

// ArgsToGoArgs converts a Unix-style argument list such that:
// all -- switches before the first non-switch argument are converted to - switches
// all - switches are split into a set of single switches with a -
// first non- args stops the process.
// so, e.g., ps -aux turns into ps -a -u -x
// ls -al --somelongthing becomes ls -a -l -somelongthing
func ArgsToGoArgs(args []string) []string {
	var out []string
	for i, f := range args {
		if strings.HasPrefix(f, "--") {
			out = append(out, f[1:])
			continue
		}
		if strings.HasPrefix(f, "-") {
			fs := strings.Split(f[1:], "")
			for _, ff := range fs {
				out = append(out, "-"+ff)
			}
			continue
		}
		out = append(out, args[i:]...)
		break
	}
	return out
}

// OSArgsToGoArgs converts os.Args to Unix-style args.
// The first argument, i.e. the executable name, is removed.
// ArgsToGoArgs is called with the rest of the args
func OSArgsToGoArgs() []string {
	return ArgsToGoArgs(os.Args[1:])
}
