// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ulog exposes logging via a Go interface.
//
// ulog has three implementations of the Logger interface: a Go standard
// library "log" package Logger, a kernel syslog (dmesg) Logger, and a test
// Logger that logs via a test's testing.TB.Logf.
package ulog

import (
	"log"
	"os"
	"testing"
)

// Logger is a log receptacle.
//
// It puts your information somewhere for safekeeping.
type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
}

// Log is a Logger that prints to stderr, like the default log package.
var Log = log.New(os.Stderr, "", log.LstdFlags)

// TestLogger is a Logger implementation that logs to testing.TB.Logf.
type TestLogger struct {
	TB testing.TB
}

// Printf formats according to the format specifier and prints to a unit test's log.
func (tl TestLogger) Printf(format string, v ...interface{}) {
	tl.TB.Logf(format, v...)
}

// Print formats according the default formats for v and prints to a unit test's log.
func (tl TestLogger) Print(v ...interface{}) {
	tl.TB.Log(v...)
}
