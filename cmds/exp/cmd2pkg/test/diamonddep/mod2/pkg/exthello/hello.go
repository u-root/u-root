// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package exthello has one external dependency.
package exthello

import (
	hello1 "github.com/u-root/gobusybox/test/diamonddep/mod1/pkg/hello"
	hello3 "github.com/u-root/gobusybox/test/diamonddep/mod3/pkg/hello"
)

func Hello() string {
	return "test/diamonddep/mod2/exthello: " + hello1.Hello() + " and " + hello3.Hello()
}
