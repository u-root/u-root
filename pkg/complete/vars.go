// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import "errors"

// Debug is a package level variable which can be
// set to, e.g., log.Printf if you want lots of debug.
var (
	Debug       = func(s string, v ...interface{}) {}
	ErrEOL      = errors.New("end of line")
	ErrEmptyEnv = errors.New("empty environment variable")
)
