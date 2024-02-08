// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ulogtest implement the Logger interface via a test's testing.TB.Logf.
package ulogtest

import (
	"testing"
)

// Logger is a Logger implementation that logs to testing.TB.Logf.
type Logger struct {
	TB testing.TB
}

// Printf formats according to the format specifier and prints to a unit test's log.
func (tl Logger) Printf(format string, v ...interface{}) {
	tl.TB.Logf(format, v...)
}
