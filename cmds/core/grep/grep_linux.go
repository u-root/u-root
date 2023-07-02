// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import flag "github.com/spf13/pflag"

func init() {
	flag.BoolVarP(&mainParams.quiet, "quiet", "q", false, "Don't print matches; exit on first match")
}
