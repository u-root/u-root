// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"os"
)

// Usage wraps a passed in func() with a func() that sets
// os.Args[0] to a string and then calls the func().
//
// It is intended to be called with Usage function from a flag package,
// such as flag or spf13/pflag.
// E.g., flag.usage = util.Usage(flag.Usage, "some message")
//
// Usage must not import "flag", since callers might use an alternate flags
// package such as spf13/pflag, and would set Usage for a flag
// package that the caller is not using.
func Usage(wrapUsage func(), message string) func() {
	return func() {
		os.Args[0] = message
		wrapUsage()
	}
}
